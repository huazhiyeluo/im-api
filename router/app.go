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

	r.POST("/user/editUser", service.EditUser)

	r.POST("/user/getFriendGroup", service.GetFriendGroup)
	r.POST("/user/getFriendList", service.GetFriendList)
	r.POST("/user/getFriendOne", service.GetFriendOne)
	r.POST("/user/addFriend", service.AddFriend)
	r.POST("/user/delFriend", service.DelFriend)

	r.POST("/user/getGroupList", service.GetGroupList)
	r.POST("/user/getGroupOne", service.GetGroupOne)
	r.POST("/user/getGroupUser", service.GetGroupUser)
	r.POST("/user/editGroup", service.EditGroup)
	r.POST("/user/joinGroup", service.JoinGroup)
	r.POST("/user/quitGroup", service.QuitGroup)

	r.POST("/user/getApplyList", service.GetApplyList)
	r.POST("/user/operateApply", service.OperateApply)

	r.POST("/user/chatMsg", service.ChatMsg)
	r.POST("/attach/upload", service.Upload)

	//websocket
	r.GET("/chat", server.Chat)

	return r
}
