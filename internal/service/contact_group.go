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

// 1-1、群组列表
func GetContactGroupList(c *gin.Context) {

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

	var dataUsers []*schema.ResContactGroup
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
		temp := schema.GetResContactGroup(tempGroups[v.ToId], v)
		dataUsers = append(dataUsers, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataUsers,
	})
}

// 1-2、群组-one
func GetContactGroupOne(c *gin.Context) {

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

	dataGroup := schema.GetResContactGroup(group, contactGroup)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataGroup,
	})
}

// 1-3、群成员列表
func GetContactGroupUser(c *gin.Context) {
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
	var dataGroupUsers []*schema.ResContactGroupUser
	for _, v := range contactUsers {
		dataGroupUsers = append(dataGroupUsers, schema.GetResContactGroupUser(tempToUsers[v.FromId], v))
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataGroupUsers,
	})
}

// 2-1、加入群组
func JoinContactGroup(c *gin.Context) {
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

// 2-2、群删除成员|群解散
func QuitContactGroup(c *gin.Context) {
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
		toMap["group"] = schema.GetResContactGroup(group, &model.ContactGroup{})
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
		toMap["user"] = schema.GetResContactGroupUser(fromUser, &model.ContactGroup{})
		toMap["group"] = schema.GetResContactGroup(group, &model.ContactGroup{})
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

// 3-1、修改群联系人信息
func ActContactGroup(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if contactGroup.FromId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "没有加入群"})
		return
	}
	nowtime := time.Now().Unix()

	var updatesContactGroup []*model.Fields
	updatesContactGroup = append(updatesContactGroup, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesContactGroup = append(updatesContactGroup, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	contactGroup, err = model.ActContactGroup(fromId, toId, updatesContactGroup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contactGroup,
	})
}
