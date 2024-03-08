package service

import (
	"context"
	"demoapi/model"
	"demoapi/schema"
	"demoapi/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// type 消息类型
	MSG_TYPE_SINGLE = 1 // 单聊消息
	MSG_TYPE_ROOM   = 2 // 群聊消息
	MSG_TYPE_HEART  = 3 // 心跳消息
	MSG_TYPE_NOTICE = 4 // 通知消息
	MSG_TYPE_ACK    = 5 // 应答消息

	// media（type 1|2） 消息展示样式
	MSG_MEDIA_TEXT  = 1 // 文本
	MSG_MEDIA_IMAGE = 2 // 图片
	MSG_MEDIA_AUDIO = 3 // 音频
	MSG_MEDIA_VIDEO = 4 // 视频
	MSG_MEDIA_FILE  = 5 // 文件
	MSG_MEDIA_EMOJI = 6 // 表情

	// media（type 4） 消息展示样式
	MSG_MEDIA_ONLINE  = 1 // 上线
	MSG_MEDIA_OFFLINE = 2 // 下线
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Id            string          //唯一标识符
	Uid           uint64          // UID
	Manager       *Manager        //对应管理者
	Conn          *websocket.Conn //websocket 连接
	Message       chan *Message   // 消息
	HeartbeatTime int64           //心跳时间
	LoginTime     int64           //登录时间
}

type Message struct {
	FromId     uint64 // ID [主]
	ToId       uint64 // ID [从]
	MsgType    uint32 // 消息类型 1私信 2群 3广播
	MsgMedia   uint32 // 图片类型 1文字 2图片 3 音频 4 视频
	Content    string // 内容
	CreateTime int64  // 创建时间
	Status     uint32 // 状态
}

func UpgradeWebSocket(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return nil, err
	}
	return conn, nil
}

// chat连接
func Chat(c *gin.Context, manager *Manager) {

	query := c.Request.URL.Query()
	uid := uint64(utils.StringToUint32(query.Get("uid")))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Chat", err)
	}

	nowtime := time.Now().Unix()
	client := &Client{
		Id:            utils.GenGUID(),
		Uid:           uid,
		Manager:       manager,
		Conn:          conn,
		Message:       make(chan *Message),
		HeartbeatTime: nowtime,
		LoginTime:     nowtime,
	}
	client.Manager.Register <- client

	go client.ReadData()
	go client.WriteData()
}

// 读取消息
func (client *Client) ReadData() {
	for {
		log.Println("ReadData0", client.Id, client.Uid)
		msg := &Message{}
		err := client.Conn.ReadJSON(msg)
		if err != nil {
			client.Manager.UnRegisterClient(client)
			log.Println("ReadData1", err, msg)
			return
		}
		log.Println("ReadData内容", msg)
		client.Manager.Message <- msg
	}
}

// 写入消息
func (client *Client) WriteData() {
	for {
		select {
		case msg := <-client.Message:
			log.Println("WriteData内容", msg)
			if err := client.Conn.WriteJSON(msg); err != nil {
				log.Println("Error writing to WebSocket:", err, client.Conn)
				return
			}
		}
	}
}

func (client *Client) Dispatch(msg *Message) {
	switch msg.MsgType {
	case MSG_TYPE_SINGLE:
		client.SendMsg(msg, msg.ToId)
	case MSG_TYPE_ROOM:
		client.SendGroupMsg(msg)
	case MSG_TYPE_HEART:
		client.SendHeartMsg(msg)
	case MSG_TYPE_NOTICE:
		client.SendNoticeMsg(msg)
	}
}

func (client *Client) SendMsg(msg *Message, toId uint64) {
	log.Println("sendMsg", msg)
	// 将消息加入节点的消息队列
	if c, ok := client.Manager.Clients.Load(toId); ok {
		c.(*Client).Message <- msg
	}
}

func (client *Client) SendGroupMsg(msg *Message) {
	log.Println("sendGroupMsg", msg)
	contacts, _ := model.GetGroupContactList(msg.ToId, 2)
	for _, v := range contacts {
		if v.FromId != msg.FromId {
			log.Println("sendGroupMsg", v.FromId)
			client.SendMsg(msg, v.FromId)
		}
	}
}

func (client *Client) SendHeartMsg(msg *Message) {
	client.HeartbeatTime = msg.CreateTime
}

func (client *Client) SendNoticeMsg(msg *Message) {
	log.Println("SendNoticeMsg", msg)
	// 将消息加入节点的消息队列
	if c, ok := client.Manager.Clients.Load(msg.ToId); ok {
		c.(*Client).Message <- msg
	}
}

// ----------------------------------------------------------------	----------------------------------------------------------------
func ChatMsg(c *gin.Context) {
	data := schema.ChatMsg{}
	c.Bind(&data)
	ctx := context.Background()
	var rkey string
	if data.MsgType == 1 {
		if data.FromId > data.ToId {
			rkey = fmt.Sprintf("msg_%d_%d", data.ToId, data.FromId)
		} else {
			rkey = fmt.Sprintf("msg_%d_%d", data.FromId, data.ToId)
		}
	}
	if data.MsgType == 2 {
		rkey = fmt.Sprintf("msg_%d_%d", 0, data.FromId)
	}
	var chats []string
	var err error
	if data.IsRev == 1 {
		chats, err = utils.RDB.ZRevRange(ctx, rkey, data.Start, data.End).Result()
	} else {
		chats, err = utils.RDB.ZRange(ctx, rkey, data.Start, data.End).Result()
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "数据错误"})
		return
	}
	newChats := utils.ReverseStringArray(chats)
	var tempChats []*model.Message
	for _, v := range newChats {
		msg := &model.Message{}
		json.Unmarshal([]byte(v), msg)
		tempChats = append(tempChats, msg)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": tempChats,
	})
}
