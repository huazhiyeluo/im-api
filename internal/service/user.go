package service

import (
	"imapi/internal/model"
	"imapi/internal/schema"
	"net/http"

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
	updateData := &model.User{
		Uid:      data.Uid,
		Username: data.Username,
		Info:     data.Info,
		Avatar:   data.Avatar,
	}
	user, err = model.UpdateUser(updateData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}
