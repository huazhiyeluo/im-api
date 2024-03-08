package service

import (
	"demoapi/model"
	"demoapi/schema"
	"demoapi/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetContactList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	uid := uint64(utils.ToNumber(data["uid"]))

	contacts, err := model.GetContactList(uid, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toIds := []uint64{}

	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	targetUsers, err := model.GetUserByUids(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetUsers := make(map[uint64]*model.User)
	for _, v := range targetUsers {
		tempTargetUsers[v.Uid] = v
	}

	// onlines := data.CheckOnline(toIds)

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.ToId,
			Username: tempTargetUsers[v.ToId].Username,
			Avatar:   tempTargetUsers[v.ToId].Avatar,
			// IsOnline: onlines[v.ToId],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

func GetGroupList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	uid := uint64(utils.ToNumber(data["uid"]))

	contacts, err := model.GetContactList(uid, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toIds := []uint64{}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	targetGroups, err := model.GetGroupByGroupIds(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetGroups := make(map[uint64]*model.Group)
	for _, v := range targetGroups {
		tempTargetGroups[v.GroupId] = v
	}

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.ToId,
			Username: tempTargetGroups[v.ToId].Name,
			Avatar:   tempTargetGroups[v.ToId].Icon,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

func GetGroupUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	groupId := uint64(utils.ToNumber(data["group_id"]))
	contacts, err := model.GetGroupContactList(groupId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	fromIds := []uint64{}
	for _, v := range contacts {
		fromIds = append(fromIds, v.FromId)
	}

	targetUsers, err := model.GetUserByUids(fromIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetUsers := make(map[uint64]*model.User)
	for _, v := range targetUsers {
		tempTargetUsers[v.Uid] = v
	}

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.FromId,
			Username: tempTargetUsers[v.FromId].Username,
			Avatar:   tempTargetUsers[v.FromId].Avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

func AddFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	uid := uint64(utils.ToNumber(data["uid"]))
	targetId := uint64(utils.ToNumber(data["targetId"]))

	insertContachData := &model.Contact{
		FromId: uid,
		ToId:   targetId,
		Type:   1,
		Desc:   "",
	}
	contact, err := model.CreateContact(insertContachData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertFriendContactData := &model.Contact{
		FromId: uid,
		ToId:   targetId,
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

func JoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	uid := uint64(utils.ToNumber(data["uid"]))
	targetId := uint64(utils.ToNumber(data["targetId"]))

	insertContachData := &model.Contact{
		FromId: uid,
		ToId:   targetId,
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
