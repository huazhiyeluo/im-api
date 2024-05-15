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

// 1-1、好友列表
func GetFriendList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	contacts, err := model.GetContactUserList(fromId)
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
		temp := schema.GetResFriend(tempToUsers[v.ToId], v)
		dataUsers = append(dataUsers, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 1-2、好友-one
func GetFriendOne(c *gin.Context) {

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

	contactUser, err := model.GetContactUserOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	toUser, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	dataUser := schema.GetResFriend(toUser, contactUser)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUser,
	})
}

// 1-3、好友分组
func GetFriendGroup(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["ownUid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	ownUid := uint32(utils.ToNumber(data["ownUid"]))
	contactGroups, err := model.GetFriendGroup(ownUid)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	var dataContactGroups []*schema.ResContactGroup
	for _, v := range contactGroups {
		dataContactGroups = append(dataContactGroups, &schema.ResContactGroup{FriendGroupId: v.FriendGroupId, OwnerUid: v.OwnerUid, Name: v.Name})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactGroups,
	})
}

// 1-4、添加好友
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

	cantactUser, err := model.GetContactUserOne(fromId, toId)
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

// 1-5、删除好友
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

	contact, err := model.DeleteContactUser(fromId, toId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	contact, err = model.DeleteContactUser(toId, fromId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contact))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	toMap := make(map[string]interface{})
	toMap["user"] = schema.GetResFriend(fromUser, &model.ContactUser{})
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(fromId, toId, string(toMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	fromMap := make(map[string]interface{})
	fromMap["user"] = schema.GetResFriend(toUser, &model.ContactUser{})
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(toId, fromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
