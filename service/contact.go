package service

import (
	"demoapi/model"
	"demoapi/schema"
	"demoapi/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1、联系人列表
func GetContactList(c *gin.Context, manager *Manager) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	contacts, err := model.GetContactList(fromId, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toIds := []uint64{}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	toUsers, err := model.GetUserByUids(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempToUsers := make(map[uint64]*model.User)
	for _, v := range toUsers {
		tempToUsers[v.Uid] = v
	}

	onlines := manager.CheckUserOnlineStatus(toIds)

	var dataUsers []*schema.ResFriend
	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResFriend{
			Uid:      v.ToId,
			Username: tempToUsers[v.ToId].Username,
			Avatar:   tempToUsers[v.ToId].Avatar,
			IsOnline: onlines[v.ToId],
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 2、群组列表
func GetGroupList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	contacts, err := model.GetContactList(fromId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toIds := []uint64{}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	groups, err := model.GetGroupByGroupIds(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempGroups := make(map[uint64]*model.Group)
	for _, v := range groups {
		tempGroups[v.GroupId] = v
	}

	var dataUsers []*schema.ResGroup
	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResGroup{
			GroupId: v.ToId,
			Name:    tempGroups[v.ToId].Name,
			Icon:    tempGroups[v.ToId].Icon,
			Info:    tempGroups[v.ToId].Info,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 3、群成员列表
func GetGroupUser(c *gin.Context, manager *Manager) {
	data := make(map[string]interface{})
	c.Bind(&data)

	groupId := uint64(utils.ToNumber(data["groupId"]))
	contacts, err := model.GetGroupContactList(groupId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	fromIds := []uint64{}
	for _, v := range contacts {
		fromIds = append(fromIds, v.FromId)
	}

	toUsers, err := model.GetUserByUids(fromIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempToUsers := make(map[uint64]*model.User)
	for _, v := range toUsers {
		tempToUsers[v.Uid] = v
	}

	onlines := manager.CheckUserOnlineStatus(fromIds)

	var dataUsers []*schema.ResFriend
	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResFriend{
			Uid:      v.ToId,
			Username: tempToUsers[v.ToId].Username,
			Avatar:   tempToUsers[v.ToId].Avatar,
			IsOnline: onlines[v.ToId],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 4、添加好友
func AddFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))

	insertContachData := &model.Contact{
		FromId: fromId,
		ToId:   toId,
		Type:   1,
		Desc:   "",
	}
	contact, err := model.CreateContact(insertContachData)
	if err != nil {
		log.Println(contact)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertFriendContactData := &model.Contact{
		FromId: toId,
		ToId:   fromId,
		Type:   1,
		Desc:   "",
	}
	contact, err = model.CreateContact(insertFriendContactData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contact,
	})

}

// 5、加入群组
func JoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))

	insertContachData := &model.Contact{
		FromId: fromId,
		ToId:   toId,
		Type:   2,
		Desc:   "",
	}
	contact, err := model.CreateContact(insertContachData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contact,
	})

}
