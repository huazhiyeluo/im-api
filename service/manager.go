package service

import (
	"context"
	"demoapi/model"
	"demoapi/utils"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

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
		case message := <-m.Message:
			// 处理接收到的消息
			// 示例中，直接将消息转发给目标客户端
			message.CreateTime = time.Now().Unix()
			m.StoreRedis(message)
			if client, ok := m.Clients.Load(message.FromId); ok {
				client.(*Client).Dispatch(message)
			}
		}
	}
}

// 注册客户端
func (m *Manager) RegisterClient(client *Client) {
	m.Clients.Store(client.Uid, client)
	go m.SendUserStatus(client.Uid, MSG_MEDIA_ONLINE)
	fmt.Println("Client Registered:", client.Uid)
}

// 关闭客户端
func (m *Manager) UnRegisterClient(client *Client) {
	client.Conn.Close()
	client.IsOnline = false
	m.Clients.Delete(client.Uid)
	go m.SendUserStatus(client.Uid, MSG_MEDIA_OFFLINE)
	fmt.Println("Client Unregistered:", client.Uid)
}

// 设置用户状态
func (m *Manager) SendUserStatus(uid uint64, msgMedia uint32) {
	contacts, err := model.GetContactList(uid, 1)
	if err != nil {
		log.Println(err)
	}
	for _, v := range contacts {
		if client, ok := m.Clients.Load(v.ToId); ok {
			msg := &Message{FromId: v.FromId, ToId: v.ToId, MsgType: 4, MsgMedia: msgMedia}
			client.(*Client).Manager.Message <- msg
		}
	}
}

// 清理连接
func (m *Manager) CleanConnection() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection", r)
		}
	}()
	nowtime := time.Now().Unix()
	m.Clients.Range(func(k, v interface{}) bool {
		c := v.(*Client)
		if c.HeartbeatTime+25 <= nowtime {
			fmt.Println("cleanConnection", c.Uid, c.HeartbeatTime)
			m.UnRegisterClient(c)
		}
		return true
	})

}

// 消息存贮
func (m *Manager) StoreRedis(msg *Message) {
	fmt.Println("StoreRedis", msg)

	ctx := context.Background()
	var rkey string

	if !utils.IsContainUint32(msg.MsgType, []uint32{1, 2}) {
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
		fmt.Println("Error converting struct to JSON:", err)
		return
	}
	res, err := utils.RDB.ZRevRange(ctx, rkey, 0, -1).Result()
	if err != nil {
		log.Println(err)
	}
	score := float64(cap(res)) + 1
	count, err := utils.RDB.ZAdd(ctx, rkey, redis.Z{score, jsonData}).Result()
	if err != nil {
		log.Println("StoreRedis", err)
	}
	log.Println("count", count, "StoreRedis", rkey)
}
