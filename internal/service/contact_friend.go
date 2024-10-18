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

// 1-1-1、好友分组
func GetContactFriendGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["ownerUid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	ownerUid := uint64(utils.ToNumber(data["ownerUid"]))
	contactFriends, err := model.GetFriendGroup(ownerUid)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	var dataContactFriendGroups []*schema.ResContactFriendGroup
	for _, v := range contactFriends {
		dataContactFriendGroups = append(dataContactFriendGroups, schema.GetResContactFriendGroup(v))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactFriendGroups,
	})
}

// 1-1-2、添加|编辑好友分组
func EditContactFriendGroup(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	friendGroupId := uint32(utils.ToNumber(data["friendGroupId"]))
	ownerUid := uint64(utils.ToNumber(data["ownerUid"]))
	name := utils.ToString(data["name"])

	if ownerUid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "分组名称不能为空"})
		return
	}

	friendGroup := &model.FriendGroup{}
	var err error

	if friendGroupId == 0 {
		friendGroup, err = model.GetFriendGroupByName(ownerUid, name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
		if friendGroup.FriendGroupId != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经有分组"})
			return
		}
		insertFriendGroupData := &model.FriendGroup{
			OwnerUid: ownerUid,
			Name:     name,
		}
		friendGroup, err = model.CreateFriendGroup(insertFriendGroupData)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
	} else {
		friendGroup, err = model.GetFriendGroupByName(ownerUid, name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
		if friendGroup.FriendGroupId != 0 && friendGroup.FriendGroupId != friendGroupId {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经有分组"})
			return
		}
		updateFriendGroupData := &model.FriendGroup{
			FriendGroupId: friendGroupId,
			Name:          name,
		}
		friendGroup, err = model.UpdateFriendGroup(updateFriendGroupData)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": schema.GetResContactFriendGroup(friendGroup),
	})
}

// 1-1-3、删除好友分组
func DelContactFriendGroup(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["friendGroupId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "分组ID不存在"})
		return
	}
	friendGroupId := uint32(utils.ToNumber(data["friendGroupId"]))
	friendGroup, err := model.GetFriendGroupByFriendGroupId(friendGroupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if friendGroup.IsDefault == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "默认分组不允许删除"})
		return
	}

	_, err = model.DeleteFriendGroup(friendGroupId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", friendGroup))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	friendGroup, err = model.GetFriendGroupByIsDefault(friendGroup.OwnerUid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	updateFriendGroup := &model.ContactFriend{FromId: friendGroup.OwnerUid, FriendGroupId: friendGroup.FriendGroupId}
	model.UpdateContactFriendByFriendGroupId(friendGroupId, updateFriendGroup)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": schema.GetResContactFriendGroup(friendGroup),
	})
}

// 1-1-4、分排序
func SortContactFriendGroup(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	for _, v := range data["data"].([]interface{}) {
		temp := v.(map[string]interface{})
		updateFriendGroupData := &model.FriendGroup{
			FriendGroupId: uint32(utils.ToNumber(temp["friendGroupId"])),
			Sort:          uint32(utils.ToNumber(temp["friendGroupId"])),
		}
		friendGroup, err := model.UpdateFriendGroup(updateFriendGroupData)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", friendGroup))
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 1-2、好友列表
func GetContactFriendList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	contacts, err := model.GetContactFriendList(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	toIds := []uint64{fromId}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	toUsers, err := model.FindUserByUids(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	var dataUsers []*schema.ResUser
	for _, v := range toUsers {
		temp := schema.GetResUser(v)
		dataUsers = append(dataUsers, temp)
	}

	var dataContactFriends []*schema.ResContactFriend
	for _, v := range contacts {
		temp := schema.GetResContactFriend(v)
		dataContactFriends = append(dataContactFriends, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"contactFriends": dataContactFriends,
			"users":          dataUsers,
		},
	})
}

// 1-3、好友-one
func GetContactFriendOne(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "UID不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	contactFriend, err := model.GetContactFriendOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}
	dataContactFriend := schema.GetResContactFriend(contactFriend)

	toUser, err := model.FindUserByUid(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
		return
	}
	dataUser := schema.GetResUser(toUser)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"contactFriend": dataContactFriend,
			"user":          dataUser,
		},
	})
}

// 2-1、添加好友
func AddContactFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	reason := utils.ToString(data["reason"])
	remark := utils.ToString(data["remark"])
	friendGroupId := uint32(utils.ToNumber(data["friendGroupId"]))

	if fromId == 0 || toId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "参数错误"})
	}

	if fromId == toId {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "不允许添加自己"})
		return
	}

	fromUser, err := model.FindUserByUid(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
		return
	}
	if fromUser.Uid == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "用户不存在"})
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

	cantactFriend, err := model.GetContactFriendOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if cantactFriend.FromId != 0 && cantactFriend.JoinTime > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经是好友"})
		return
	}

	apply, err := model.FindApplyByTwoId(fromId, toId, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		_sendApplyNotic(apply, fromUser, toUser)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经申请过了，请等待"})
		return
	}

	apply, err = model.FindApplyByTwoId(toId, fromId, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id != 0 {
		_sendApplyNotic(apply, toUser, fromUser)
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "对方申请你为好友，请通过"})
		return
	}

	insertApplyData := &model.Apply{
		FromId:        fromId,
		ToId:          toId,
		Type:          1,
		Reason:        reason,
		Remark:        remark,
		FriendGroupId: friendGroupId,
		OperateTime:   time.Now().Unix(),
	}
	apply, err = model.CreateApply(insertApplyData)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", apply))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	_sendApplyNotic(apply, fromUser, toUser)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}

func _sendApplyNotic(apply *model.Apply, fromUser *model.User, toUser *model.User) {
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
}

// 2-2、邀请好友
func InviteContactFriend(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	tempToIds := data["toIds"].([]interface{})
	toIds := []uint64{}
	for _, toId := range tempToIds {
		toIds = append(toIds, uint64(utils.ToNumber(toId)))
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))
	groupId := uint64(utils.ToNumber(data["groupId"]))

	group, err := model.FindGroupByGroupId(groupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}

	contactFriends, err := model.GetContactFriendByToIds(fromId, toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}
	for _, v := range contactFriends {
		toMap := make(map[string]interface{})
		toMap["group"] = schema.GetResGroup(group)
		toMapStr, _ := json.Marshal(toMap)
		go server.UserFriendMsg(fromId, v.ToId, string(toMapStr), server.MSG_MEDIA_INVITE)
	}

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
	toMap["contactFriend"] = schema.GetResContactFriend(&model.ContactFriend{FromId: toId, ToId: fromId})
	toMap["user"] = schema.GetResUser(fromUser)
	toMapStr, _ := json.Marshal(toMap)
	go server.UserFriendNoticeMsg(fromId, toId, string(toMapStr), server.MSG_MEDIA_FRIEND_DELETE)

	fromMap := make(map[string]interface{})
	fromMap["contactFriend"] = schema.GetResContactFriend(&model.ContactFriend{FromId: fromId, ToId: toId})
	fromMap["user"] = schema.GetResUser(toUser)
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
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))
	nowtime := time.Now().Unix()

	var updatesContactFriend []*model.Fields
	updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	contactFriend, err := model.ActContactFriend(fromId, toId, updatesContactFriend)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	dataContactFriend := schema.GetResContactFriend(contactFriend)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactFriend,
	})
}
