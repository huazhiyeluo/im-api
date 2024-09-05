package service

import (
	"fmt"
	"log"
	"net/http"
	"qqapi/internal/login"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/utils"

	"github.com/gin-gonic/gin"
)

// 目前支持用户名密码注册
func Register(c *gin.Context) {
	data := schema.RegisterData{}
	c.Bind(&data)
	if data.Password != data.Repassword {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码和确认密码不一致"})
		return
	}

	if data.Username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "用户名不得为空"})
		return
	}

	usermapSso, err := model.GetUsermapSsoByUsername(data.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": err.Error()})
		return
	}
	if usermapSso.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "账号已经注册"})
		return
	}

	phone := fmt.Sprintf("p_%s", data.Username)
	email := fmt.Sprintf("e_%s", data.Username)
	usermapSso, err = model.CreateUsermapSso(&model.UsermapSso{Siteuid: utils.GenGUID(), Username: data.Username, Phone: phone, Email: email, Password: utils.GenMd5(data.Password)})
	if err != nil {
		log.Print(usermapSso)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "注册失败"})
		return
	}

	loginData := &schema.LoginData{Platform: "account", Username: data.Username, Password: data.Password, Nickname: data.Nickname, Avatar: data.Avatar}
	cin := schema.GetHeader(c)

	res, err := login.Login(cin, loginData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": res,
	})
}

func Bind(c *gin.Context) {
	data := &schema.BindData{}
	c.Bind(data)
	if data.Password != data.Repassword {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码和确认密码不一致"})
		return
	}
	if data.Username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "用户名不得为空"})
		return
	}

	usermapSso, err := model.GetUsermapSsoByUsername(data.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": err.Error()})
		return
	}
	if usermapSso.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "账号已经注册"})
		return
	}

	phone := fmt.Sprintf("p_%s", data.Username)
	email := fmt.Sprintf("e_%s", data.Username)
	usermapSso, err = model.CreateUsermapSso(&model.UsermapSso{Siteuid: utils.GenGUID(), Username: data.Username, Phone: phone, Email: email, Password: utils.GenMd5(data.Password)})
	if err != nil {
		log.Print(usermapSso)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "注册失败"})
		return
	}

	usermapBind, err := model.CreateUsermapBind(&model.UsermapBind{Uid: data.Uid, Siteuid: usermapSso.Siteuid, Sid: 1})
	if err != nil {
		log.Print(usermapBind)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "注册失败"})
		return
	}

	user, err := model.FindUserByUid(data.Uid)
	if err != nil {
		log.Print(user)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "注册失败"})
		return
	}
	res := schema.GetResUser(user)

	usermapSso, err = model.GetUserMapSsoMix(data.Uid)
	if err != nil {
		log.Print(usermapBind)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "注册失败"})
		return
	}
	if usermapSso.Id != 0 {
		res.Username = usermapSso.Username
		res.Phone = usermapSso.Phone
		res.Email = usermapSso.Email
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": res,
	})
}
