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

// 1、创建群
func EditGroup(c *gin.Context) {
	data := schema.EditGroup{}
	c.Bind(&data)

	fromUser, err := model.FindUserByUid(data.OwnerUid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	group, err := model.FindGroupByName(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经存在"})
		return
	}

	nowtime := time.Now().Unix()

	insertData := &model.Group{
		OwnerUid:   data.OwnerUid,
		Type:       data.Type,
		Name:       data.Name,
		Icon:       data.Icon,
		Info:       data.Info,
		Num:        1,
		Exp:        0,
		CreateTime: nowtime,
		UpdateTime: nowtime,
	}
	group, err = model.CreateGroup(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertContactGroupData := &model.ContactGroup{
		FromId:   data.OwnerUid,
		ToId:     group.GroupId,
		Level:    1,
		Remark:   "",
		Nickname: "",
		JoinTime: nowtime,
	}
	contactGroup, err := model.CreateContactGroup(insertContactGroupData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//1、告诉请求的人消息
	fromMap := make(map[string]interface{})
	fromMap["user"] = schema.GetResGroupUser(fromUser, &model.ContactGroup{})
	fromMap["group"] = schema.GetResGroup(group, contactGroup)
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(group.OwnerUid, group.OwnerUid, string(fromMapStr), server.MSG_MEDIA_GROUP_CREATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}

// 2-1、群组列表
func GetGroupList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)
	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))
	contactGroups, err := model.GetContactGroupList(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	var dataUsers []*schema.ResGroup
	if len(contactGroups) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": dataUsers,
		})
		return
	}
	toIds := []uint64{}
	for _, v := range contactGroups {
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

	for _, v := range contactGroups {
		temp := schema.GetResGroup(tempGroups[v.ToId], v)
		dataUsers = append(dataUsers, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 2-2、群组-one
func GetGroupOne(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)
	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "群不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	contactGroup, err := model.GetContactGroupOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	dataGroup := schema.GetResGroup(group, contactGroup)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataGroup,
	})
}

// 3、群成员列表
func GetGroupUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	groupId := uint64(utils.ToNumber(data["groupId"]))
	contactUsers, err := model.GetGroupUser(groupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	fromIds := []uint64{}
	for _, v := range contactUsers {
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
	var dataGroupUsers []*schema.ResGroupUser
	for _, v := range contactUsers {
		dataGroupUsers = append(dataGroupUsers, schema.GetResGroupUser(tempToUsers[v.FromId], v))
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataGroupUsers,
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

	cantactGroup, err := model.GetContactGroupOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if cantactGroup.FromId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经在群里"})
		return
	}

	apply, err := model.FindApplyByTwoId(fromId, toId, 2)
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
	tempApply := schema.GetResApplyGroup(apply, fromUser, group)

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
		toMap["group"] = schema.GetResGroup(group, &model.ContactGroup{})
		toMapStr, _ := json.Marshal(toMap)
		go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_GROUP_DISBAND)

		group, err := model.DeleteGroup(toId)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", group))
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}

		contact, err := model.DeleteContactGroupAll(toId)
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
		toMap["user"] = schema.GetResGroupUser(fromUser, &model.ContactGroup{})
		toMap["group"] = schema.GetResGroup(group, &model.ContactGroup{})
		toMapStr, _ := json.Marshal(toMap)
		go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_GROUP_DELETE)

		contact, err := model.DeleteContactGroup(fromId, toId)
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
