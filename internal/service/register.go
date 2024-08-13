package service

import (
	"net/http"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	data := schema.Register{}
	c.Bind(&data)
	if data.Password != data.Repassword {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码和确认密码不一致"})
		return
	}
	nowtime := time.Now().Unix()
	insertData := &model.User{
		Username:   data.Username,
		Password:   utils.GenMd5(data.Password),
		Email:      "",
		Phone:      "",
		CreateTime: nowtime,
	}
	user, err := model.CreateUser(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	token := setToken(user.Uid)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"token": token,
		"data":  schema.GetResUser(user),
	})
}
