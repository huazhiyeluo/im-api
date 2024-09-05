package schema

import "github.com/gin-gonic/gin"

type CommonData struct {
	Devname  string `json:"devname"`
	Deviceid string `json:"deviceid"`
}

func GetHeader(c *gin.Context) *CommonData {
	return &CommonData{
		Devname:  c.GetHeader("devname"),
		Deviceid: c.GetHeader("deviceid"),
	}
}
