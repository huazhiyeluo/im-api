package service

import (
	"encoding/json"
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/server"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/***********************************************/
//编辑用户
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
	toMap["user"] = schema.GetResFriend(user, &model.ContactUser{})
	toMapStr, _ := json.Marshal(toMap)
	go server.UserInfoNoticeMsg(data.Uid, string(toMapStr))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}
