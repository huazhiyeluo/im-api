package router

import (
	"demoapi/service"
	"time"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	//启动管理者服务
	manager := service.NewManager()
	InitTimer(manager)
	go manager.Start()

	r := gin.Default()
	r.Static("/static", "static")

	r.POST("/user/login", service.Login)
	r.POST("/user/register", service.Register)

	r.POST("/user/getContactList", func(c *gin.Context) {
		service.GetContactList(c, manager)
	})
	r.POST("/user/getGroupList", service.GetGroupList)
	r.POST("/user/getGroupUser", func(c *gin.Context) {
		service.GetGroupUser(c, manager)
	})

	r.POST("/user/chatMsg", service.ChatMsg)
	r.POST("/attach/upload", service.Upload)

	r.POST("/user/editUser", service.EditUser)
	r.POST("/user/editGroup", service.EditGroup)

	r.POST("/user/addFriend", service.AddFriend)
	r.POST("/user/joinGroup", service.JoinGroup)

	//websocket
	r.GET("/chat", func(c *gin.Context) {
		service.Chat(c, manager)
	})

	return r
}

func InitTimer(m *service.Manager) {
	Timer(time.Second*3, time.Second*3, m.CleanConnection)
}

func Timer(delay time.Duration, tick time.Duration, fun func()) {
	go func() {
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				// 定时器触发的处理逻辑
				fun()
				t.Reset(tick)
			}
		}
	}()
}
