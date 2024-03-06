package router

import (
	"demoapi/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "static")
	//websocket
	r.GET("/chat", service.Chat)

	r.POST("/user/login", service.Login)
	r.POST("/user/register", service.Register)

	r.POST("/user/getContactList", service.GetContactList)
	r.POST("/user/getGroupList", service.GetGroupList)
	r.POST("/user/getGroupUser", service.GetGroupUser)

	r.POST("/user/chatMsg", service.ChatMsg)
	r.POST("/attach/upload", service.Upload)

	r.POST("/user/editUser", service.EditUser)
	r.POST("/user/editGroup", service.EditGroup)

	r.POST("/user/addFriend", service.AddFriend)
	r.POST("/user/joinGroup", service.JoinGroup)

	return r
}
