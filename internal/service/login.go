package service

import (
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	data := schema.Login{}
	c.Bind(&data)
	user, err := model.FindUserByName(data.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if user.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}
	if user.Password != utils.GenMd5(data.Password) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码错误"})
		return
	}

	token := setToken(user.Uid)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"token": token,
		"data":  schema.GetUser(user),
	})
}
