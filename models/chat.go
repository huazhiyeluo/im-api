package models

import (
	"context"
	"demoapi/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/set"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Id            string
	Conn          *websocket.Conn
	Message       chan Message
	Uid           uint64
	IsOnline      bool          //是否在线 0，否 1是
	HeartbeatTime int64         //心跳时间
	LoginTime     int64         //登录时间
	GroupSets     set.Interface //好友 / 群
}

type Message struct {
	SenderId   uint64
	ReceiverId uint64
	Type       uint32
	Media      uint32
	Content    string
	CreateTime int64
}

var clientMap map[uint64]*Client = make(map[uint64]*Client) // 用户ID到WebSocket节点的映射
var rwLocker sync.RWMutex

func Chat(c *gin.Context) {

	query := c.Request.URL.Query()
	uid := uint64(utils.StringToUint32(query.Get("uid")))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Chat", err)
	}

	nowtime := time.Now().Unix()
	client := &Client{
		Id:            conn.RemoteAddr().String(),
		Conn:          conn,
		Message:       make(chan Message),
		Uid:           uid,
		IsOnline:      true,
		HeartbeatTime: nowtime,
		LoginTime:     nowtime,
		GroupSets:     set.New(set.ThreadSafe),
	}

	registerClient(client)

	go client.readData()
	go client.writeData()
}

func registerClient(client *Client) {
	rwLocker.Lock()
	clientMap[client.Uid] = client
	rwLocker.Unlock()
}

func unregisterClient(uid uint64) {
	rwLocker.Lock()
	delete(clientMap, uid)
	rwLocker.Unlock()
}

func dispatch(data []byte) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Println("dispatch", err)
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.ReceiverId, *msg)
	case 2:
		sendGroupMsg(*msg)
	}
}

func sendMsg(receiverId uint64, msg Message) {
	rwLocker.Lock()
	client, ok := clientMap[receiverId]
	rwLocker.Unlock()

	msg.CreateTime = time.Now().Unix()
	senderId := msg.SenderId
	ttype := msg.Type

	if ttype == 2 {
		senderId = 0
	}

	if ok {
		log.Println("sendMsg", receiverId, msg)
		// 将消息加入节点的消息队列
		client.Message <- msg
	}

	ctx := context.Background()
	var rkey string
	if senderId > msg.ReceiverId {
		rkey = fmt.Sprintf("msg_%d_%d", msg.ReceiverId, senderId)
	} else {
		rkey = fmt.Sprintf("msg_%d_%d", senderId, msg.ReceiverId)
	}

	log.Println("sendMsg", rkey)

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
		log.Println(err)
	}
	log.Println("count", count)
}

func sendGroupMsg(msg Message) {
	contacts, _ := GetGroupContactList(msg.ReceiverId, 2)
	for _, v := range contacts {
		if v.Uid != msg.SenderId {
			log.Println("sendGroupMsg", v.Uid)
			sendMsg(v.Uid, msg)
		}

	}
}

func (client *Client) Heartbeat(currentTime int64) {
	client.HeartbeatTime = currentTime
}

func CleanConnection() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	nowtime := time.Now().Unix()
	for _, v := range clientMap {
		if v.HeartbeatTime+25 <= nowtime {
			v.Conn.Close()
			v.IsOnline = false
			unregisterClient(v.Uid)
		}
	}
}

func CheckOnline(uids []uint64) map[uint64]bool {
	res := make(map[uint64]bool)
	for _, uid := range uids {
		if _, ok := clientMap[uid]; ok {
			res[uid] = clientMap[uid].IsOnline
		} else {
			res[uid] = false
		}
	}
	return res
}

//----------------------------------------------------------------	--	----------------------------------------------------------------

var udpData chan Message = make(chan Message)

func (client *Client) readData() {
	for {
		log.Println("readData0", client.Id)
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("readData1", err, p)
			return
		}

		msg := &Message{}
		err = json.Unmarshal(p, msg)
		if err != nil {
			log.Println("readData2", err)
			return
		}

		if msg.Type == 3 {
			log.Println("readData 心跳", string(p))
			nowtime := time.Now().Unix()
			client.Heartbeat(nowtime)
		} else {
			log.Println("readData内容", string(p))
			udpData <- *msg
		}

	}
}

func (client *Client) writeData() {
	for {
		select {
		case msg := <-client.Message:
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error converting struct to JSON:", err)
				continue
			}
			log.Println("writeData内容", string(jsonData))
			if err := client.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Error writing to WebSocket:", err, client.Conn)
				return
			}
		}
	}
}

func init() {
	go readUdp()
	go writeUdp()
}

func writeUdp() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 5000,
	})
	if err != nil {
		log.Println("Error connecting to UDP:", err)
	}
	defer conn.Close()
	for {
		select {
		case message := <-udpData:
			jsonData, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error converting struct to JSON:", err)
				return
			}
			log.Println("writeUdp内容", string(jsonData))
			_, err = conn.Write(jsonData)
			if err != nil {
				fmt.Println("Error write from UDP:", err)
			}
		}
	}
}

func readUdp() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 5000,
	})
	if err != nil {
		log.Println("Error connecting to UDP:", err)
	}
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		data := buffer[:n]
		log.Println("readUdp内容", string(data))
		dispatch(data)
	}
}
