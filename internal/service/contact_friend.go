package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/server"
	"qqapi/internal/utils"
	"qqapi/third_party/log"
	"time"

	"github.com/gin-gonic/gin"
)

// 1-1、好友分组
func GetContactFriendGroup(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["ownUid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	ownUid := uint32(utils.ToNumber(data["ownUid"]))
	contactFriends, err := model.GetFriendGroup(ownUid)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	var datacontactFriends []*schema.ResContactFriendGroup
	for _, v := range contactFriends {
		datacontactFriends = append(datacontactFriends, &schema.ResContactFriendGroup{FriendGroupId: v.FriendGroupId, OwnerUid: v.OwnerUid, Name: v.Name})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": datacontactFriends,
	})
}

// 1-2、好友列表
func GetContactFriendList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	contacts, err := model.GetContactFriendList(fromId)
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
	var dataUsers []*schema.ResContactFriend
	for _, v := range contacts {
		temp := schema.GetResContactFriend(tempToUsers[v.ToId], v)
		dataUsers = append(dataUsers, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 1-3、好友-one
func GetContactFriendOne(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	ContactFriend, err := model.GetContactFriendOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toUser, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	dataUser := schema.GetResContactFriend(toUser, ContactFriend)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUser,
	})
}

// 2-1、添加好友
func AddContactFriend(c *gin.Context) {
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

	cantactUser, err := model.GetContactFriendOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if cantactUser.FromId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经是好友"})
		return
	}

	apply, err := model.FindApplyByTwoId(fromId, toId, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经申请过了，请等待"})
		return
	}

	apply, err = model.FindApplyByTwoId(toId, fromId, 1)
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
	tempApply := schema.GetResApplyUser(apply, fromUser, toUser)

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

// 2-2、删除好友
func DelContactFriend(c *gin.Context) {
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

	contact, err := model.DeleteContactFriend(fromId, toId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	contact, err = model.DeleteContactFriend(toId, fromId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	toMap := make(map[string]interface{})
	toMap["user"] = schema.GetResContactFriend(fromUser, &model.ContactFriend{})
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(fromId, toId, string(toMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	fromMap := make(map[string]interface{})
	fromMap["user"] = schema.GetResContactFriend(toUser, &model.ContactFriend{})
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(toId, fromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 3-1、修改好友联系人信息
func ActContactFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	contactFriend, err := model.GetContactFriendOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if contactFriend.FromId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "没有加好友"})
		return
	}
	nowtime := time.Now().Unix()

	var updatesContactFriend []*model.Fields
	updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	contactFriend, err = model.ActContactFriend(fromId, toId, updatesContactFriend)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contactFriend,
	})
}
