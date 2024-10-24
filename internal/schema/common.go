package schema

import "github.com/gin-gonic/gin"

type CommonData struct {
	Devname  string
	Deviceid string
}

func GetHeader(c *gin.Context) *CommonData {
	return &CommonData{
		Devname:  c.GetHeader("devname"),
		Deviceid: c.GetHeader("deviceid"),
	}
}
