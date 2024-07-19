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

// 1、编辑用户
func EditUser(c *gin.Context) {
	data := schema.EditUser{}
	c.Bind(&data)
	user, err := model.FindUserByUid(data.Uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if user.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}
	nowtime := time.Now().Unix()

	updateData := &model.User{
		Uid:        data.Uid,
		Username:   data.Username,
		Info:       data.Info,
		Avatar:     data.Avatar,
		UpdateTime: nowtime,
	}
	user, err = model.UpdateUser(updateData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	toMap := make(map[string]interface{})
	toMap["user"] = schema.GetResContactFriend(user, &model.ContactFriend{})
	toMapStr, _ := json.Marshal(toMap)
	go server.UserInfoNoticeMsg(data.Uid, string(toMapStr))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}

// 2、编辑用户
func ActUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	uid := uint64(utils.ToNumber(data["uid"]))

	user, err := model.FindUserByUid(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if user.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}
	nowtime := time.Now().Unix()

	var updatesUser []*model.Fields
	updatesUser = append(updatesUser, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesUser = append(updatesUser, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	user, err = model.ActUser(uid, updatesUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	toMap := make(map[string]interface{})
	toMap["user"] = schema.GetResContactFriend(user, &model.ContactFriend{})
	toMapStr, _ := json.Marshal(toMap)
	go server.UserInfoNoticeMsg(uid, string(toMapStr))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}
