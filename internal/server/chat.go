package server

import (
	"context"
	"encoding/json"
	"errors"
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
	MSG_TYPE_NOTICE    = 3 // 通知消息
	MSG_TYPE_ACK       = 4 // 应答消息
	MSG_TYPE_BROADCAST = 5 // 广播消息

	// media（type 1|2） 消息展示样式
	MSG_MEDIA_TEXT       = 1  // 文本
	MSG_MEDIA_IMAGE      = 2  // 图片
	MSG_MEDIA_AUDIO      = 3  // 音频
	MSG_MEDIA_VIDEO      = 4  // 视频
	MSG_MEDIA_FILE       = 5  // 文件
	MSG_MEDIA_EMOJI      = 6  // 表情
	MSG_MEDIA_NOT_ONLINE = 10 // 不在线
	MSG_MEDIA_NO_CONNECT = 11 // 未接通
	MSG_MEDIA_TIMES      = 12 // 通话时长
	MSG_MEDIA_OFF        = 13 // 挂断电话

	// media（type 3） 消息展示样式
	MSG_MEDIA_OFFLINE_PACK = 10 // 挤下线
	MSG_MEDIA_ONLINE       = 11 // 上线
	MSG_MEDIA_OFFLINE      = 12 // 下线
	MSG_MEDIA_USERINFO     = 13 // 用户信息

	MSG_MEDIA_FRIEND_ADD    = 21 // 添加好友
	MSG_MEDIA_FRIEND_AGREE  = 22 // 成功添加好友
	MSG_MEDIA_FRIEND_REFUSE = 23 // 拒绝添加好友
	MSG_MEDIA_FRIEND_DELETE = 24 // 删除好友

	MSG_MEDIA_GROUP_CREATE  = 30 // 创建群
	MSG_MEDIA_GROUP_JOIN    = 31 // 添加群
	MSG_MEDIA_GROUP_AGREE   = 32 // 成功添加群
	MSG_MEDIA_GROUP_REFUSE  = 33 // 拒绝添加群
	MSG_MEDIA_GROUP_DELETE  = 34 // 退出群
	MSG_MEDIA_GROUP_DISBAND = 35 // 解散群

	// media（type 4） 消息展示样式
	MSG_MEDIA_PHONE_CONTACT = 0 // 发起聊天
	MSG_MEDIA_PHONE_QUIT    = 1 // 退出聊天
	MSG_MEDIA_PHONE_OPEN    = 2 // 接通聊天
	MSG_MEDIA_PHONE_ICE     = 3 // ICE候选
	MSG_MEDIA_PHONE_OFFER   = 4 // offer
	MSG_MEDIA_PHONE_ANSWER  = 5 // answer
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

type MessageContent struct {
	Data string `json:"data"` //数据
	Url  string `json:"url"`  //链接地址
	Name string `json:"name"` //文件名
}

type Message struct {
	Id         string          `json:"id"`         // ID
	FromId     uint64          `json:"fromId"`     // ID [主]
	ToId       uint64          `json:"toId"`       // ID [从]
	MsgType    uint32          `json:"msgType"`    // 消息类型 1私信 2群 3广播 4通知
	MsgMedia   uint32          `json:"msgMedia"`   // 图片类型 1文字 2图片 3 音频 4 视频
	Content    *MessageContent `json:"content"`    // 内容
	CreateTime int64           `json:"createTime"` // 创建时间
	Status     uint32          `json:"status"`     // 状态
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

	checkSessionKey(c, uid)

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
			// 记录错误类型
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Logger.Error(fmt.Sprintf("Liao ReadData Unexpected close error: %v", err))
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Logger.Info(fmt.Sprintf("Liao ReadData Normal close: %v", err))
			} else {
				var closeErr *websocket.CloseError
				if errors.As(err, &closeErr) {
					log.Logger.Error(fmt.Sprintf("Liao ReadData Close error: Code: %v, Text: %v", closeErr.Code, closeErr.Text))
				} else {
					log.Logger.Error(fmt.Sprintf("Liao ReadData Read error: %v", err))
				}
			}
			// 从管理器中注销客户端
			manager.UnRegister <- client
			log.Logger.Info(fmt.Sprintf("Liao Client unregistered: %v , %v", client.Id, client.Uid))
			return
		}

		log.Logger.Info(fmt.Sprintf("Liao ReadData0 : %v , %v", client.Id, client.Uid))
		jsonData, _ := json.Marshal(msg)
		log.Logger.Info(fmt.Sprintf("Liao ReadData内容: %v ", string(jsonData)))

		CreateMsg(msg)
	}
}

// 写入消息
func (client *Client) WriteData() {
	for {
		select {
		case msg := <-client.Message:
			jsonData, _ := json.Marshal(msg)
			log.Logger.Info(fmt.Sprintf("Liao WriteData内容 %v ", string(jsonData)))

			if err := client.Conn.WriteJSON(msg); err != nil {
				log.Logger.Info(fmt.Sprintf("Liao Error writing to WebSocket: %v , %v", err, client.Uid))
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
	case MSG_TYPE_NOTICE:
		SendNoticeMsg(msg)
	case MSG_TYPE_ACK:
		SendAckMsg(msg)
	case MSG_TYPE_BROADCAST:
		SendBroadcastMsg(msg)
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
	} else {
		StoreUnreadMessage(toId, msg)
	}
}

// 2 群消息
func SendGroupMsg(msg *Message) {
	log.Logger.Info(fmt.Sprintf("sendGroupMsg %v ", msg))
	contacts, _ := model.GetGroupUser(msg.ToId)
	for _, v := range contacts {
		if v.FromId != msg.FromId {
			log.Logger.Info(fmt.Sprintf("sendGroupMsg %v ", v.FromId))
			SendMsg(msg, v.FromId)
		}
	}
}

// 3 通知消息
func SendNoticeMsg(msg *Message) {
	log.Logger.Info(fmt.Sprintf("SendNoticeMsg: %v ", msg))
	SendMsg(msg, msg.ToId)
}

// 4 应答消息
func SendAckMsg(msg *Message) {
	log.Logger.Info(fmt.Sprintf("SendAckMsg: %v ", msg))
	// 将消息加入节点的消息队列
	if v, ok := manager.Clients.Load(msg.ToId); ok {
		client := v.(*Client)
		client.Message <- msg
	} else {
		if utils.IsContainUint32(msg.MsgMedia, []uint32{0, 2}) {
			go CreateMsg(&Message{FromId: msg.ToId, ToId: msg.FromId, MsgType: MSG_TYPE_ACK, MsgMedia: MSG_MEDIA_PHONE_QUIT, Content: &MessageContent{Data: ""}})
			go CreateMsg(&Message{FromId: msg.ToId, ToId: msg.FromId, MsgType: MSG_TYPE_SINGLE, MsgMedia: MSG_MEDIA_NOT_ONLINE, Content: &MessageContent{Data: "对方不在线"}})
			go CreateMsg(&Message{FromId: msg.FromId, ToId: msg.ToId, MsgType: MSG_TYPE_SINGLE, MsgMedia: MSG_MEDIA_NO_CONNECT, Content: &MessageContent{Data: "呼叫未接通"}})
		}
	}
}

// 5 广播消息
func SendBroadcastMsg(msg *Message) {
	manager.Clients.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		msg.ToId = client.Uid
		client.Message <- msg
		return true
	})
}

func setSessionKey(uid uint64, sessionkey string) string {
	rkey := model.RksessionKey(uid)
	utils.RDB.Set(context.TODO(), rkey, sessionkey, time.Minute*time.Duration(0))
	utils.RDB.ExpireAt(context.TODO(), rkey, time.Now().Add(time.Minute*60*24*2))
	return sessionkey
}

func getSessionKey(uid uint64) string {
	rkey := model.RksessionKey(uid)
	sessionkey, err := utils.RDB.Get(context.TODO(), rkey).Result()
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", err))
	}
	return sessionkey
}

func checkSessionKey(c *gin.Context, uid uint64) {
	sessionKey := ""
	session, _ := c.Request.Cookie("sessionKey")
	if session != nil {
		sessionKey = session.Value
	}

	oldSessionKey := getSessionKey(uid)
	if oldSessionKey != "" {
		if oldSessionKey != sessionKey {
			if oldclient, ok := manager.Clients.Load(uid); ok {
				log.Logger.Info(fmt.Sprintf("下线: %v", uid))
				msg := &Message{FromId: uid, ToId: uid, MsgType: MSG_TYPE_NOTICE, MsgMedia: MSG_MEDIA_OFFLINE_PACK, Content: &MessageContent{Data: "下线"}}
				oldclient.(*Client).Message <- msg
			}
		}
	} else {
		setSessionKey(uid, sessionKey)
	}

}
