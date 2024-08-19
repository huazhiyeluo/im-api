package server

import (
	"context"
	"encoding/json"
	"fmt"
	"qqapi/internal/model"
	"qqapi/internal/utils"
	"qqapi/third_party/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	manager = NewManager()
)

var mu sync.Mutex

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
			msg.Id = utils.GenGUID()
			msg.CreateTime = time.Now().Unix()
			go m.StoreData(msg)
			mu.Lock()
			Dispatch(msg)
			mu.Unlock()
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
	//未读消息就不在存储了
	if msg.Status == 1 {
		return
	}

	jsonData, _ := json.Marshal(msg)
	log.Logger.Info(fmt.Sprintf("StoreData %v", string(jsonData)))
	if !utils.IsContainUint32(msg.MsgType, []uint32{1, 2, 3, 4}) {
		return
	}
	if utils.IsContainUint32(msg.MsgType, []uint32{3}) {
		if utils.IsContainUint32(msg.MsgMedia, []uint32{10, 11, 12}) {
			return
		}
	}
	if utils.IsContainUint32(msg.MsgType, []uint32{4}) {
		if utils.IsContainUint32(msg.MsgMedia, []uint32{3, 4, 5}) {
			return
		}
	}

	m.StoreRedisMessage(msg)

	content, _ := json.Marshal(msg.Content)

	go model.CreateMessage(&model.Message{
		Id:         msg.Id,
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
		rkey = model.Rkmsg(0, msg.FromId)
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
	go model.CreateMessageUnread(&model.MessageUnread{
		Uid:        uid,
		MsgId:      msg.Id,
		CreateTime: msg.CreateTime,
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
	log.Logger.Info(fmt.Sprintf("Liao CreateMsg %v", string(jsonData)))
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
func CheckUserOnlineStatus(uids []uint64) map[uint64]uint32 {
	res := make(map[uint64]uint32)
	for _, uid := range uids {
		if _, ok := manager.Clients.Load(uid); ok {
			res[uid] = 1
		} else {
			res[uid] = 0
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
		ids := []string{}
		for _, v := range temps {
			ids = append(ids, v.MsgId)
		}
		msgs, err := model.GetMessageAll(ids)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("PushUnreadMessage %v", err))
		}
		msgMaps := make(map[string]*model.Message)
		for _, v := range msgs {
			msgMaps[v.Id] = v
		}

		for _, v := range temps {
			if _, ok := msgMaps[v.MsgId]; !ok {
				continue
			}
			content := &MessageContent{}
			json.Unmarshal([]byte(msgMaps[v.MsgId].Content), content)
			msg := &Message{
				FromId:     msgMaps[v.MsgId].FromId,
				ToId:       msgMaps[v.MsgId].ToId,
				MsgType:    msgMaps[v.MsgId].MsgType,
				MsgMedia:   msgMaps[v.MsgId].MsgMedia,
				Content:    content,
				CreateTime: msgMaps[v.MsgId].CreateTime,
				Status:     1,
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
			msg.Status = 1
			CreateMsg(msg)
		}
	}
	utils.RDB.Del(ctx, rkey)
	model.DeleteMessageUnread(uid)
}
