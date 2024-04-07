package server

import (
	"context"
	"encoding/json"
	"fmt"
	"imapi/internal/model"
	"imapi/internal/utils"
	"imapi/third_party/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	manager = NewManager()
)

func App() {
	manager.Start()
}

type Manager struct {
	Clients    sync.Map
	Register   chan *Client
	UnRegister chan *Client
	Message    chan *Message
}

func NewManager() *Manager {
	return &Manager{
		Clients:    sync.Map{},
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Message:    make(chan *Message),
	}
}

// 启动服务
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.RegisterClient(client)
		case client := <-m.UnRegister:
			m.UnRegisterClient(client)
		case msg := <-m.Message:
			// 处理接收到的消息
			jsonData, _ := json.Marshal(msg)
			log.Logger.Info(fmt.Sprintf("Start %v", string(jsonData)))
			msg.CreateTime = time.Now().Unix()
			m.StoreData(msg)
			Dispatch(msg)
		}
	}
}

// 注册客户端
func (m *Manager) RegisterClient(client *Client) {
	m.Clients.Store(client.Uid, client)
	go UserStatusNoticeMsg(client.Uid, MSG_MEDIA_ONLINE)
	go PushUnreadMessage(client.Uid)
	fmt.Println("Client Registered:", client.Uid)
}

// 关闭客户端
func (m *Manager) UnRegisterClient(client *Client) {
	if client != nil {
		client.Conn.Close()
		m.Clients.Delete(client.Uid)
	}
	go UserStatusNoticeMsg(client.Uid, MSG_MEDIA_OFFLINE)
}

// 消息存贮
func (m *Manager) StoreData(msg *Message) {
	jsonData, _ := json.Marshal(msg)
	log.Logger.Info(fmt.Sprintf("StoreData %v", string(jsonData)))
	if !utils.IsContainUint32(msg.MsgType, []uint32{1, 2, 3}) {
		return
	}
	if utils.IsContainUint32(msg.MsgType, []uint32{3}) {
		if utils.IsContainUint32(msg.MsgMedia, []uint32{10, 11, 12}) {
			return
		}
	}

	m.StoreRedisMessage(msg)

	content, _ := json.Marshal(msg.Content)

	go model.CreateMessage(&model.Message{
		FromId:     msg.FromId,
		ToId:       msg.ToId,
		MsgType:    msg.MsgType,
		MsgMedia:   msg.MsgMedia,
		Content:    string(content),
		CreateTime: msg.CreateTime,
		Status:     msg.Status,
	})
}

func (m *Manager) StoreRedisMessage(msg *Message) {
	ctx := context.Background()
	var rkey string

	if msg.MsgType == 1 {
		if msg.FromId > msg.ToId {
			rkey = model.Rkmsg(msg.ToId, msg.FromId)
		} else {
			rkey = model.Rkmsg(msg.FromId, msg.ToId)
		}
	}
	if msg.MsgType == 2 {
		rkey = fmt.Sprintf("msg_%d_%d", 0, msg.FromId)
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", err))
		return
	}
	res, err := utils.RDB.ZRevRange(ctx, rkey, 0, -1).Result()
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", err))
	}
	score := float64(cap(res)) + 1
	count, err := utils.RDB.ZAdd(ctx, rkey, redis.Z{Score: score, Member: jsonData}).Result()
	if err != nil {
		log.Logger.Info(fmt.Sprintf("StoreRedis %v", err))
	}
	utils.RDB.ExpireAt(ctx, rkey, time.Now().Add(time.Hour*24*7))
	log.Logger.Info(fmt.Sprintf("count %v StoreRedis %v", count, rkey))

}

// 不在线存贮消息
func StoreUnreadMessage(uid uint64, msg *Message) {
	fmt.Println("StoreUnreadMessage", uid, msg)
	if !utils.IsContainUint32(msg.MsgType, []uint32{1, 2, 3}) {
		return
	}
	if utils.IsContainUint32(msg.MsgType, []uint32{3}) {
		if utils.IsContainUint32(msg.MsgMedia, []uint32{10, 11, 12}) {
			return
		}
	}
	StoreUnreadRedisMessage(uid, msg)

	content, _ := json.Marshal(msg.Content)
	go model.CreateMessageUnread(&model.MessageUnread{
		Uid:        uid,
		FromId:     msg.FromId,
		ToId:       msg.ToId,
		MsgType:    msg.MsgType,
		MsgMedia:   msg.MsgMedia,
		Content:    string(content),
		CreateTime: msg.CreateTime,
		Status:     msg.Status,
	})
}
func StoreUnreadRedisMessage(uid uint64, msg *Message) {
	ctx := context.Background()

	rkey := model.RkUreadMsg(uid)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", err))
		return
	}
	count, err := utils.RDB.LPush(ctx, rkey, jsonData).Result()
	if err != nil {
		log.Logger.Info(fmt.Sprintf("StoreRedis %v", err))
	}
	utils.RDB.ExpireAt(ctx, rkey, time.Now().Add(time.Hour*24*7))
	log.Logger.Info(fmt.Sprintf("count %v StoreUnreadMessage %v", count, rkey))
}

// 设置创建通知
func CreateMsg(msg *Message) {
	jsonData, _ := json.Marshal(msg)
	log.Logger.Info(fmt.Sprintf("CreateMsg %v", string(jsonData)))
	manager.Message <- msg
}

// 清理连接
func CleanConnection() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection", r)
		}
	}()
	nowtime := time.Now().Unix()
	manager.Clients.Range(func(k, v interface{}) bool {
		c := v.(*Client)
		if c.HeartbeatTime+25 <= nowtime {
			fmt.Println("cleanConnection", c.Uid, c.HeartbeatTime)
			manager.UnRegisterClient(c)
		}
		return true
	})
}

// 用户在线状态检查
func CheckUserOnlineStatus(uids []uint64) map[uint64]bool {
	res := make(map[uint64]bool)
	for _, uid := range uids {
		if _, ok := manager.Clients.Load(uid); ok {
			res[uid] = true
		} else {
			res[uid] = false
		}
	}
	return res
}

// 推送用户未读消息
func PushUnreadMessage(uid uint64) {
	ctx := context.Background()
	rkey := model.RkUreadMsg(uid)
	count, err := utils.RDB.LLen(ctx, rkey).Result()
	if err != nil {
		log.Logger.Info(fmt.Sprintf("PushUnreadMessage %v", err))
	}
	countDb, err := model.GetMessageUnreadCount(uid)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("PushUnreadMessage %v", err))
	}
	if count == 0 || count != countDb {
		temps, err := model.GetMessageUnread(uid)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("PushUnreadMessage %v", err))
		}
		for _, v := range temps {
			content := &MessageContent{}
			json.Unmarshal([]byte(v.Content), content)
			msg := &Message{
				FromId:     v.FromId,
				ToId:       v.ToId,
				MsgType:    v.MsgType,
				MsgMedia:   v.MsgMedia,
				Content:    content,
				CreateTime: v.CreateTime,
				Status:     v.Status,
			}
			CreateMsg(msg)
		}
	} else {
		var i int64
		for i = 0; i < count; i++ {
			temp, err := utils.RDB.RPop(ctx, rkey).Result()
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", err))
				continue
			}
			msg := &Message{}
			json.Unmarshal([]byte(temp), msg)
			CreateMsg(msg)
		}
	}
	utils.RDB.Del(ctx, rkey)
	model.DeleteMessageUnread(uid)
}
