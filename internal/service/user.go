package service

import (
	"encoding/json"
	"log"
	"net/http"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/server"
	"qqapi/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 0、单个用户
func GetOneUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "条件不存在"})
		return
	}
	uid := uint64(utils.ToNumber(data["uid"]))

	user, err := model.FindUserByUid(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": schema.GetResUser(user),
	})
}

// 1、编辑用户
func ActUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
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
	toMap["user"] = schema.GetResUser(user)
	toMapStr, _ := json.Marshal(toMap)
	go server.UserInfoNoticeMsg(uid, string(toMapStr))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 3、搜索用户
func SearchUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["keyword"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "条件不存在"})
		return
	}
	keyword := utils.ToString(data["keyword"])
	pageSize := uint32(utils.ToNumber(data["pageSize"]))
	if pageSize == 0 {
		pageSize = 10
	}
	pageNum := uint32(utils.ToNumber(data["pageNum"]))
	if pageNum == 0 {
		pageNum = 1
	}
	users, count, err := model.FindUserByKeyword(pageSize, pageNum, keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	var dataUsers []*schema.ResUser
	for _, v := range users {
		temp := schema.GetResUser(v)
		dataUsers = append(dataUsers, temp)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"users": dataUsers,
			"count": count,
		},
	})
}

// 4、更新设备Token
func ActDeviceToken(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
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

	var updatesDeviceToken []*model.Fields
	updatesDeviceToken = append(updatesDeviceToken, &model.Fields{Field: "last_login_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesDeviceToken = append(updatesDeviceToken, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	deviceToken, err := model.ActDeviceToken(uid, updatesDeviceToken)
	if err != nil {
		log.Print(deviceToken)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
