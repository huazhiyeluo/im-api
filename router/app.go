package router

import (
	"qqapi/internal/server"
	"qqapi/internal/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "static")

	// login
	r.POST("/user/login", service.Login)

	// register
	r.POST("/user/register", service.Register)
	r.POST("/user/bind", service.Bind)

	// user
	r.POST("/user/actUser", service.ActUser)
	r.POST("/user/searchUser", service.SearchUser)
	r.POST("/user/actDeviceToken", service.ActDeviceToken)

	// group
	r.POST("/user/createGroup", service.CreateGroup)
	r.POST("/user/actGroup", service.ActGroup)
	r.POST("/user/searchGroup", service.SearchGroup)

	// contact_friend
	r.POST("/user/getContactFriendGroup", service.GetContactFriendGroup)
	r.POST("/user/addContactFriendGroup", service.AddContactFriendGroup)
	r.POST("/user/delContactFriendGroup", service.DelContactFriendGroup)

	r.POST("/user/getContactFriendList", service.GetContactFriendList)
	r.POST("/user/getContactFriendOne", service.GetContactFriendOne)
	r.POST("/user/addContactFriend", service.AddContactFriend)
	r.POST("/user/inviteContactFriend", service.InviteContactFriend)
	r.POST("/user/delContactFriend", service.DelContactFriend)
	r.POST("/user/actContactFriend", service.ActContactFriend)

	// contact_group
	r.POST("/user/getContactGroupList", service.GetContactGroupList)
	r.POST("/user/getContactGroupOne", service.GetContactGroupOne)
	r.POST("/user/getContactGroupUser", service.GetContactGroupUser)
	r.POST("/user/joinContactGroup", service.JoinContactGroup)
	r.POST("/user/quitContactGroup", service.QuitContactGroup)
	r.POST("/user/delContactGroup", service.DelContactGroup)
	r.POST("/user/actContactGroup", service.ActContactGroup)

	// apply
	r.POST("/user/getApplyList", service.GetApplyList)
	r.POST("/user/operateApply", service.OperateApply)

	// message
	r.POST("/user/chatMsg", service.ChatMsg)

	// upload
	r.POST("/attach/upload", service.Upload)

	//websocket
	r.GET("/chat", server.Chat)

	return r
}
