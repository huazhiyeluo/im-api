package server

import (
	"fmt"
	"imapi/internal/model"
	"imapi/internal/utils"
	"imapi/third_party/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// type 消息类型
	MSG_TYPE_HEART     = 0 // 心跳消息
	MSG_TYPE_SINGLE    = 1 // 单聊消息
	MSG_TYPE_ROOM      = 2 // 群聊消息
	MSG_TYPE_BROADCAST = 3 // 广播消息
	MSG_TYPE_NOTICE    = 4 // 通知消息
	MSG_TYPE_ACK       = 5 // 应答消息

	// media（type 1|2） 消息展示样式
	MSG_MEDIA_TEXT  = 1 // 文本
	MSG_MEDIA_IMAGE = 2 // 图片
	MSG_MEDIA_AUDIO = 3 // 音频
	MSG_MEDIA_VIDEO = 4 // 视频
	MSG_MEDIA_FILE  = 5 // 文件
	MSG_MEDIA_EMOJI = 6 // 表情

	// media（type 4） 消息展示样式
	MSG_MEDIA_OFFLINE_PACK = 10 // 挤下线
	MSG_MEDIA_ONLINE       = 11 // 上线
	MSG_MEDIA_OFFLINE      = 12 // 下线

	MSG_MEDIA_FRIEND_ADD    = 21 // 添加好友
	MSG_MEDIA_FRIEND_AGREE  = 22 // 成功添加好友
	MSG_MEDIA_FRIEND_REFUSE = 23 // 拒绝添加好友
	MSG_MEDIA_FRIEND_DELETE = 24 // 删除好友

	MSG_MEDIA_GROUP_JOIN   = 31 // 添加群
	MSG_MEDIA_GROUP_AGREE  = 32 // 成功添加群
	MSG_MEDIA_GROUP_REFUSE = 33 // 拒绝添加群
	MSG_MEDIA_GROUP_DELETE = 24 // 删除好友
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Id            string          // 唯一标识符
	Uid           uint64          // UID
	Conn          *websocket.Conn // websocket 连接
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
		log.Logger.Info("WebSocket Upgrade Error:", log.Any("err", err))
		return nil, err
	}
	return conn, nil
}

// chat连接
func Chat(c *gin.Context) {

	query := c.Request.URL.Query()
	uid := uint64(utils.StringToUint32(query.Get("uid")))

	if oldclient, ok := manager.Clients.Load(uid); ok {
		log.Logger.Info(fmt.Sprintf("下线: %v", uid))
		msg := &Message{FromId: uid, ToId: uid, MsgType: 4, MsgMedia: MSG_MEDIA_OFFLINE_PACK, Content: "下线"}
		oldclient.(*Client).Message <- msg
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Info("Chat", log.Any("err", err))
	}
	nowtime := time.Now().Unix()
	client := &Client{
		Id:            utils.GenGUID(),
		Uid:           uid,
		Conn:          conn,
		Message:       make(chan *Message),
		HeartbeatTime: nowtime,
		LoginTime:     nowtime,
	}

	manager.Register <- client

	go client.ReadData()
	go client.WriteData()
}

// 读取消息
func (client *Client) ReadData() {
	for {
		msg := &Message{}
		err := client.Conn.ReadJSON(msg)
		if err != nil {
			manager.UnRegister <- client
			log.Logger.Info(fmt.Sprintf("ReadData1 : %v , %v", err, msg))
			return
		}
		if msg.MsgType > 0 {
			log.Logger.Info(fmt.Sprintf("ReadData0 : %v , %v", client.Id, client.Uid))
			log.Logger.Info(fmt.Sprintf("ReadData内容: %v ", msg))
		}
		manager.Message <- msg
	}
}

// 写入消息
func (client *Client) WriteData() {
	for {
		select {
		case msg := <-client.Message:
			log.Logger.Info(fmt.Sprintf("WriteData内容 %v ", msg))
			if err := client.Conn.WriteJSON(msg); err != nil {
				log.Logger.Info(fmt.Sprintf("Error writing to WebSocket: %v , %v", err, client.Uid))
				return
			}
		}
	}
}

func Dispatch(msg *Message) {
	switch msg.MsgType {
	case MSG_TYPE_HEART:
		SendHeartMsg(msg)
	case MSG_TYPE_SINGLE:
		SendMsg(msg, msg.ToId)
	case MSG_TYPE_ROOM:
		SendGroupMsg(msg)
	case MSG_TYPE_BROADCAST:
		SendBroadcastMsg(msg)
	case MSG_TYPE_NOTICE:
		SendNoticeMsg(msg)
	}
}

// 0 心跳消息
func SendHeartMsg(msg *Message) {
	if v, ok := manager.Clients.Load(msg.FromId); ok {
		client := v.(*Client)
		client.HeartbeatTime = msg.CreateTime
	}

}

// 1 个人消息
func SendMsg(msg *Message, toId uint64) {
	log.Logger.Info(fmt.Sprintf("sendMsg: %v ", msg))
	// 将消息加入节点的消息队列
	if v, ok := manager.Clients.Load(toId); ok {
		client := v.(*Client)
		client.Message <- msg
	}
}

// 2 群消息
func SendGroupMsg(msg *Message) {
	log.Logger.Info(fmt.Sprintf("sendGroupMsg %v ", msg))
	contacts, _ := model.GetGroupContactList(msg.ToId, 2)
	for _, v := range contacts {
		if v.FromId != msg.FromId {
			log.Logger.Info(fmt.Sprintf("sendGroupMsg %v ", v.FromId))
			SendMsg(msg, v.FromId)
		}
	}
}

// 3 广播消息
func SendBroadcastMsg(msg *Message) {
	manager.Clients.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		msg.ToId = client.Uid
		client.Message <- msg
		return true
	})
}

// 4 通知消息
func SendNoticeMsg(msg *Message) {
	log.Logger.Info(fmt.Sprintf("SendNoticeMsg: %v ", msg))
	// 将消息加入节点的消息队列
	if v, ok := manager.Clients.Load(msg.ToId); ok {
		client := v.(*Client)
		client.Message <- msg
	}
}
