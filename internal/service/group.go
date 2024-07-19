package service

import (
	"encoding/json"
	"net/http"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/server"
	"qqapi/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 1、创建群
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

	nowtime := time.Now().Unix()

	insertData := &model.Group{
		OwnerUid:   data.OwnerUid,
		Type:       data.Type,
		Name:       data.Name,
		Icon:       data.Icon,
		Info:       data.Info,
		Num:        1,
		Exp:        0,
		CreateTime: nowtime,
		UpdateTime: nowtime,
	}
	group, err = model.CreateGroup(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertContactGroupData := &model.ContactGroup{
		FromId:   data.OwnerUid,
		ToId:     group.GroupId,
		Level:    1,
		Remark:   "",
		Nickname: "",
		JoinTime: nowtime,
	}
	contactGroup, err := model.CreateContactGroup(insertContactGroupData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//1、告诉请求的人消息
	fromMap := make(map[string]interface{})
	fromMap["user"] = schema.GetResContactGroupUser(fromUser, &model.ContactGroup{})
	fromMap["group"] = schema.GetResContactGroup(group, contactGroup)
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(group.OwnerUid, group.OwnerUid, string(fromMapStr), server.MSG_MEDIA_GROUP_CREATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}

// 2、编辑群
func ActGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["groupId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "用户组不存在"})
		return
	}
	groupId := uint64(utils.ToNumber(data["groupId"]))

	group, err := model.FindGroupByGroupId(groupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户组不存在"})
		return
	}
	nowtime := time.Now().Unix()

	var updatesGroup []*model.Fields
	updatesGroup = append(updatesGroup, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesGroup = append(updatesGroup, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	group, err = model.ActGroup(groupId, updatesGroup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}
