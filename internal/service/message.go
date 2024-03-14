package service

import (
	"context"
	"encoding/json"
	"fmt"
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
