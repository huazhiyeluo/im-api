package service

import (
	"demoapi/models"
	"demoapi/schema"
	"demoapi/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetContactList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	uid := uint64(utils.ToNumber(data["uid"]))

	contacts, err := models.GetContactList(uid, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	targetIds := []uint64{}

	for _, v := range contacts {
		targetIds = append(targetIds, v.TargetId)
	}
	targetUsers, err := models.GetUserByUids(targetIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetUsers := make(map[uint64]*models.User)
	for _, v := range targetUsers {
		tempTargetUsers[v.Uid] = v
	}

	onlines := models.CheckOnline(targetIds)

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.TargetId,
			Username: tempTargetUsers[v.TargetId].Username,
			Avatar:   tempTargetUsers[v.TargetId].Avatar,
			IsOnline: onlines[v.TargetId],
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

	contacts, err := models.GetContactList(uid, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	targetIds := []uint64{}
	for _, v := range contacts {
		targetIds = append(targetIds, v.TargetId)
	}
	targetGroups, err := models.GetGroupByGroupIds(targetIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetGroups := make(map[uint64]*models.Group)
	for _, v := range targetGroups {
		tempTargetGroups[v.GroupId] = v
	}

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.TargetId,
			Username: tempTargetGroups[v.TargetId].Name,
			Avatar:   tempTargetGroups[v.TargetId].Icon,
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
	contacts, err := models.GetGroupContactList(groupId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	targetIds := []uint64{}

	for _, v := range contacts {
		targetIds = append(targetIds, v.Uid)
	}

	targetUsers, err := models.GetUserByUids(targetIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempTargetUsers := make(map[uint64]*models.User)
	for _, v := range targetUsers {
		tempTargetUsers[v.Uid] = v
	}

	var dataUsers []*schema.ResContact

	for _, v := range contacts {
		dataUsers = append(dataUsers, &schema.ResContact{
			Uid:      v.Uid,
			Username: tempTargetUsers[v.Uid].Username,
			Avatar:   tempTargetUsers[v.Uid].Avatar,
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

	insertContachData := &models.Contact{
		Uid:      uid,
		TargetId: targetId,
		Type:     1,
		Desc:     "",
	}
	contact, err := models.CreateContact(insertContachData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertFriendContactData := &models.Contact{
		Uid:      targetId,
		TargetId: uid,
		Type:     1,
		Desc:     "",
	}
	contact, err = models.CreateContact(insertFriendContactData)
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

	insertContachData := &models.Contact{
		Uid:      uid,
		TargetId: targetId,
		Type:     2,
		Desc:     "",
	}
	contact, err := models.CreateContact(insertContachData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contact,
	})

}
