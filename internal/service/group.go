package service

import (
	"encoding/json"
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditGroup(c *gin.Context) {
	data := schema.EditGroup{}
	c.Bind(&data)

	fromUser, err := model.FindUserByUid(data.OwnerUid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	group, err := model.FindGroupByName(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经存在"})
		return
	}

	insertData := &model.Group{
		OwnerUid: data.OwnerUid,
		Type:     data.Type,
		Name:     data.Name,
		Icon:     data.Icon,
		Info:     data.Info,
		Num:      1,
	}
	group, err = model.CreateGroup(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertContactData := &model.Contact{
		FromId: data.OwnerUid,
		ToId:   group.GroupId,
		Type:   2,
		Remark: "",
	}
	_, err = model.CreateContact(insertContactData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//1、告诉请求的人消息
	fromMap := make(map[string]interface{})
	fromMap["user"] = getResUser(fromUser)
	fromMap["group"] = getResGroup(group)
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(group.OwnerUid, group.OwnerUid, string(fromMapStr), server.MSG_MEDIA_GROUP_CREATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}
