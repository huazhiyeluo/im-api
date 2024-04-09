package service

import (
	"encoding/json"
	"fmt"
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/server"
	"imapi/internal/utils"
	"imapi/third_party/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 1、联系人列表
func GetFriendList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}

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
	toUsers, err := model.FindUserByUids(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempToUsers := make(map[uint64]*model.User)
	for _, v := range toUsers {
		tempToUsers[v.Uid] = v
	}
	var dataUsers []*schema.ResFriend
	for _, v := range contacts {
		dataUsers = append(dataUsers, getResUser(tempToUsers[v.ToId]))
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
	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))
	contacts, err := model.GetContactList(fromId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	var dataUsers []*schema.ResGroup
	if len(contacts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": dataUsers,
		})
	}
	toIds := []uint64{}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	groups, err := model.FindGroupByGroupIds(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempGroups := make(map[uint64]*model.Group)
	for _, v := range groups {
		tempGroups[v.GroupId] = v
	}

	for _, v := range contacts {
		dataUsers = append(dataUsers, getResGroup(tempGroups[v.ToId]))
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

	toUsers, err := model.FindUserByUids(fromIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempToUsers := make(map[uint64]*model.User)
	for _, v := range toUsers {
		tempToUsers[v.Uid] = v
	}
	var dataUsers []*schema.ResFriend
	for _, v := range contacts {
		dataUsers = append(dataUsers, getResUser(tempToUsers[v.FromId]))
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

	if fromId == 0 || toId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "参数错误"})
	}

	if fromId == toId {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "不允许添加自己"})
		return
	}

	fromUser, err := model.FindUserByUid(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if fromUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	toUser, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if toUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	cantact, err := model.GetGroupContactOne(fromId, toId, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if cantact.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经是好友"})
		return
	}

	apply, err := model.FindApplyByTwoId(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经申请过了，请等待"})
		return
	}

	apply, err = model.FindApplyByTwoId(toId, fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "对方申请你为好友，请通过"})
		return
	}

	insertApplyData := &model.Apply{
		FromId:      fromId,
		ToId:        toId,
		Type:        1,
		Reason:      reason,
		OperateTime: time.Now().Unix(),
	}
	apply, err = model.CreateApply(insertApplyData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", apply))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//应用数据处理
	tempApply := getResApplyUser(apply, fromUser, toUser)

	toMap := make(map[string]interface{})
	toMap["apply"] = tempApply
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(apply.FromId, apply.ToId, string(toMapStr), server.MSG_MEDIA_FRIEND_ADD)

	fromMap := make(map[string]interface{})
	fromMap["apply"] = tempApply
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(apply.ToId, apply.FromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_ADD)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}

// 4-4、删除好友
func DelFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))

	if fromId == 0 || toId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	fromUser, err := model.FindUserByUid(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if fromUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	toUser, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if toUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

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

	toMap := make(map[string]interface{})
	toMap["user"] = getResUser(fromUser)
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(fromId, toId, string(toMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	fromMap := make(map[string]interface{})
	fromMap["user"] = getResUser(toUser)
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(toId, fromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_DELETE)

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

	fromUser, err := model.FindUserByUid(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if fromUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "群不存在"})
		return
	}

	cantact, err := model.GetGroupContactOne(fromId, toId, 2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if cantact.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经在群里"})
		return
	}

	apply, err := model.FindApplyByTwoId(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经申请过了，请等待"})
		return
	}

	insertApplyData := &model.Apply{
		FromId:      fromId,
		ToId:        toId,
		Type:        2,
		Reason:      reason,
		OperateTime: time.Now().Unix(),
	}
	apply, err = model.CreateApply(insertApplyData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", apply))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//应用数据处理
	tempApply := getResApplyGroup(apply, fromUser, group)

	toMap := make(map[string]interface{})
	toMap["apply"] = tempApply
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(apply.FromId, group.OwnerUid, string(toMapStr), server.MSG_MEDIA_GROUP_JOIN)

	fromMap := make(map[string]interface{})
	fromMap["apply"] = tempApply
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_JOIN)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}

// 5-4、群删除成员|群解散
func QuitGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))

	fromUser, err := model.FindUserByUid(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if fromUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}

	if group.OwnerUid == fromId {

		toMap := make(map[string]interface{})
		toMap["group"] = getResGroup(group)
		toMapStr, _ := json.Marshal(toMap)
		go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_GROUP_DISBAND)

		group, err := model.DeleteGroup(toId)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", group))
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}

		contact, err := model.DeleteContactGroup(toId)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", contact))
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}

	} else {
		group.Num = group.Num - 1
		group, err = model.UpdateGroup(group)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}

		toMap := make(map[string]interface{})
		toMap["user"] = getResUser(fromUser)
		toMap["group"] = getResGroup(group)
		toMapStr, _ := json.Marshal(toMap)
		go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_GROUP_DELETE)

		contact, err := model.DeleteContact(fromId, toId, 2)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", contact))
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
