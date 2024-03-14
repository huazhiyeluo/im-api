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
			msg.CreateTime = time.Now().Unix()
			m.StoreData(msg)
			//除了广播消息其他消息都需要验证发送人在不在线
			if !utils.IsContainUint32(msg.MsgType, []uint32{MSG_TYPE_BROADCAST}) {
				if _, ok := m.Clients.Load(msg.FromId); ok {
					Dispatch(msg)
				}
			}
		}
	}
}

// 注册客户端
func (m *Manager) RegisterClient(client *Client) {
	m.Clients.Store(client.Uid, client)
	go UserStatusNoticeMsg(client.Uid, MSG_MEDIA_ONLINE)
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
	fmt.Println("StoreData", msg)

	ctx := context.Background()
	var rkey string

	if !utils.IsContainUint32(msg.MsgType, []uint32{1, 2, 3}) {
		return
	}
	if msg.MsgType == 1 {
		if msg.FromId > msg.ToId {
			rkey = fmt.Sprintf("msg_%d_%d", msg.ToId, msg.FromId)
		} else {
			rkey = fmt.Sprintf("msg_%d_%d", msg.FromId, msg.ToId)
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
	log.Logger.Info(fmt.Sprintf("count %v StoreRedis %v", count, rkey))

	go model.CreateMessage(&model.Message{
		FromId:     msg.FromId,
		ToId:       msg.ToId,
		MsgType:    msg.MsgType,
		MsgMedia:   msg.MsgMedia,
		Content:    msg.Content,
		CreateTime: msg.CreateTime,
		Status:     msg.Status,
	})
}

// 设置创建通知
func CreateNoticeMsg(msg *Message) {
	if _, ok := manager.Clients.Load(msg.ToId); ok {
		log.Logger.Info(fmt.Sprintf("CreateNoticeMsg %v", msg))
		manager.Message <- msg
	}
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
