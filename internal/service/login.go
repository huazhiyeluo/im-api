package service

import (
	"net/http"
	"qqapi/internal/login"
	"qqapi/internal/schema"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	data := &schema.LoginData{}
	c.Bind(data)

	cin := schema.GetHeader(c)

	spew.Dump(cin)

	res, err := login.Login(cin, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": res,
	})
}
