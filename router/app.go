package router

import (
	"imapi/internal/server"
	"imapi/internal/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "static")

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
	r.POST("/user/agreeFriend", service.AgreeFriend)
	r.POST("/user/refuseFriend", service.RefuseFriend)
	r.POST("/user/deleteFriend", service.DeleteFriend)

	r.POST("/user/joinGroup", service.JoinGroup)
	r.POST("/user/agreeJoinGroup", service.AgreeJoinGroup)
	r.POST("/user/refuseJoinGroup", service.RefuseJoinGroup)
	r.POST("/user/deleteJoinGroup", service.DeleteJoinGroup)

	//websocket
	r.GET("/chat", server.Chat)

	return r
}
