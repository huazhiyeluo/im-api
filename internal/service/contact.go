package service

import (
	"fmt"
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/server"
	"imapi/internal/utils"
	"imapi/third_party/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1、联系人列表
func GetContactList(c *gin.Context) {

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
	onlines := server.CheckUserOnlineStatus(toIds)

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
func GetGroupUser(c *gin.Context) {
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

	onlines := server.CheckUserOnlineStatus(fromIds)

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

// 4-1、添加好友
func AddFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	reason := utils.ToString(data["reason"])

	user, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if user.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	insertApplyContactData := &model.ApplyContact{
		FromId: fromId,
		ToId:   toId,
		Type:   1,
		Reason: reason,
	}
	applyContact, err := model.CreateApplyContact(insertApplyContactData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(applyContact.FromId, applyContact.ToId, server.MSG_MEDIA_FRIEND_ADD)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}

// 4-2、同意添加好友
func AgreeFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	id := uint32(utils.ToNumber(data["id"]))
	applyContact, err := model.FindApplyById(id)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	if applyContact.Id == 0 || applyContact.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "申请状态错误"})
		return
	}
	applyContact.Status = 1
	applyContact, err = model.UpdateApplyContact(applyContact)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}

	insertContachData := &model.Contact{
		FromId: applyContact.FromId,
		ToId:   applyContact.ToId,
		Type:   1,
		Remark: "",
	}
	contact, err := model.CreateContact(insertContachData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
		return
	}

	insertFriendContactData := &model.Contact{
		FromId: applyContact.ToId,
		ToId:   applyContact.FromId,
		Type:   1,
		Remark: "",
	}
	contact, err = model.CreateContact(insertFriendContactData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 5, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(applyContact.ToId, applyContact.FromId, server.MSG_MEDIA_FRIEND_AGREE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 4-3、拒绝添加好友
func RefuseFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	id := uint32(utils.ToNumber(data["id"]))
	applyContact, err := model.FindApplyById(id)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	if applyContact.Id == 0 || applyContact.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "申请状态错误"})
		return
	}
	applyContact.Status = 2
	applyContact, err = model.UpdateApplyContact(applyContact)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(applyContact.ToId, applyContact.FromId, server.MSG_MEDIA_FRIEND_REFUSE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 4-4、删除好友
func DeleteFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	contact, err := model.DeleteContact(fromId, toId, 1)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	contact, err = model.DeleteContact(toId, fromId, 1)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	go server.UserFriendNoticeMsg(fromId, toId, server.MSG_MEDIA_FRIEND_DELETE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 5-1、加入群组
func JoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	reason := utils.ToString(data["reason"])

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "群不存在"})
		return
	}

	insertApplyContactData := &model.ApplyContact{
		FromId: fromId,
		ToId:   toId,
		Type:   2,
		Reason: reason,
	}
	applyContact, err := model.CreateApplyContact(insertApplyContactData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(applyContact.FromId, group.OwnerUid, server.MSG_MEDIA_GROUP_JOIN)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}

// 5-2、同意加入群组
func AgreeJoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	id := uint32(utils.ToNumber(data["id"]))
	applyContact, err := model.FindApplyById(id)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	if applyContact.Id == 0 || applyContact.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "申请状态错误"})
		return
	}

	insertContachData := &model.Contact{
		FromId: applyContact.FromId,
		ToId:   applyContact.ToId,
		Type:   2,
		Remark: "",
	}
	contact, err := model.CreateContact(insertContachData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}

	group, err := model.FindGroupByGroupId(applyContact.ToId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(group.OwnerUid, applyContact.FromId, server.MSG_MEDIA_GROUP_AGREE)
	go server.UserGroupNoticeMsg(applyContact.ToId, server.MSG_MEDIA_GROUP_AGREE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 5-3、拒绝加入群组
func RefuseJoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	id := uint32(utils.ToNumber(data["id"]))
	applyContact, err := model.FindApplyById(id)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	if applyContact.Id == 0 || applyContact.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "申请状态错误"})
		return
	}
	applyContact.Status = 2
	applyContact, err = model.UpdateApplyContact(applyContact)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", applyContact))
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}

	group, err := model.FindGroupByGroupId(applyContact.ToId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(group.OwnerUid, applyContact.FromId, server.MSG_MEDIA_GROUP_REFUSE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 5-4、群删除成员
func DeleteJoinGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	contact, err := model.DeleteContact(fromId, toId, 2)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}

	go server.UserFriendNoticeMsg(group.OwnerUid, fromId, server.MSG_MEDIA_GROUP_DELETE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
