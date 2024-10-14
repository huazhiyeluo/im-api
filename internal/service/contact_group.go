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
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))
	contacts, err := model.GetContactGroupList(fromId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	toIds := []uint64{}
	for _, v := range contacts {
		toIds = append(toIds, v.ToId)
	}
	groups, err := model.FindGroupByGroupIds(toIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	var dataGroups []*schema.ResGroup
	for _, v := range groups {
		temp := schema.GetResGroup(v)
		dataGroups = append(dataGroups, temp)
	}

	var dataContactGroups []*schema.ResContactGroup
	for _, v := range contacts {
		temp := schema.GetResContactGroup(v)
		dataContactGroups = append(dataContactGroups, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"contactGroups": dataContactGroups,
			"groups":        dataGroups,
		},
	})
}

// 1-2、群组-one
func GetContactGroupOne(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)
	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "群不存在"})
		return
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	contactGroup, err := model.GetContactGroupOne(fromId, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	dataContactGroup := schema.GetResContactGroup(contactGroup)

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	dataGroup := schema.GetResGroup(group)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"contactGroup": dataContactGroup,
			"group":        dataGroup,
		},
	})
}

// 1-3、群成员列表
func GetContactGroupUser(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	groupId := uint64(utils.ToNumber(data["groupId"]))
	contactUsers, err := model.GetGroupUser(groupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	fromIds := []uint64{}
	for _, v := range contactUsers {
		fromIds = append(fromIds, v.FromId)
	}

	toUsers, err := model.FindUserByUids(fromIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	var dataUsers []*schema.ResUser
	for _, v := range toUsers {
		temp := schema.GetResUser(v)
		dataUsers = append(dataUsers, temp)
	}

	var dataContactGroups []*schema.ResContactGroup
	for _, v := range contactUsers {
		dataContactGroups = append(dataContactGroups, schema.GetResContactGroup(v))
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"contactGroups": dataContactGroups,
			"users":         dataUsers,
		},
	})
}

// 2-1、加入群组
func JoinContactGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	fromId := uint64(utils.ToNumber(data["fromId"]))
	toId := uint64(utils.ToNumber(data["toId"]))
	reason := utils.ToString(data["reason"])
	remark := utils.ToString(data["remark"])
	info := utils.ToString(data["info"])

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
		Remark:      remark,
		Info:        info,
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

// 2-2、退出群|群解散
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
		toMap["group"] = schema.GetResGroup(group)
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
		toMap["user"] = schema.GetResUser(fromUser)
		toMap["group"] = schema.GetResGroup(group)
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

// 2-3、群删除成员-批量
func DelContactGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	tempFromIds := data["fromIds"].([]interface{})
	fromIds := []uint64{}
	for _, fromId := range tempFromIds {
		fromIds = append(fromIds, uint64(utils.ToNumber(fromId)))
	}
	toId := uint64(utils.ToNumber(data["toId"]))

	group, err := model.FindGroupByGroupId(toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}

	contactGroups, err := model.FindContactGroupByFromIds(fromIds, toId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "操作错误"})
		return
	}
	desFromIds := []uint64{}
	for _, v := range contactGroups {
		desFromIds = append(desFromIds, v.FromId)
	}

	group.Num = group.Num - uint32(len(desFromIds))
	group, err = model.UpdateGroup(group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	users, err := model.FindUserByUids(desFromIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	mapUsers := make(map[uint64]*model.User)
	for _, v := range users {
		mapUsers[v.Uid] = v
	}

	for _, v := range contactGroups {
		toMap := make(map[string]interface{})
		toMap["user"] = schema.GetResUser(mapUsers[v.FromId])
		toMap["group"] = schema.GetResGroup(group)
		toMapStr, _ := json.Marshal(toMap)
		go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_GROUP_DELETE)
	}

	contactGroup, err := model.DeleteContactGroupByFromIds(desFromIds, toId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", contactGroup))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
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
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "群不存在"})
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
	dataContactGroup := schema.GetResContactGroup(contactGroup)

	toMap := make(map[string]interface{})
	toMap["contactGroup"] = dataContactGroup
	toMapStr, _ := json.Marshal(toMap)
	go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_CONTACT_GROUP_UPDATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactGroup,
	})
}

// 4-1、添加群管理员
func AddGroupManger(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "群不存在"})
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
	updatesContactGroup = append(updatesContactGroup, &model.Fields{Field: "group_power", Otype: 2, Value: 1})
	contactGroup, err = model.ActContactGroup(fromId, toId, updatesContactGroup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	dataContactGroup := schema.GetResContactGroup(contactGroup)

	toMap := make(map[string]interface{})
	toMap["contactGroup"] = dataContactGroup
	toMapStr, _ := json.Marshal(toMap)
	go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_CONTACT_GROUP_UPDATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactGroup,
	})
}

// 4-2、删除群管理员
func DelGroupManger(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["fromId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	fromId := uint64(utils.ToNumber(data["fromId"]))

	if _, ok := data["toId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "群不存在"})
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
	updatesContactGroup = append(updatesContactGroup, &model.Fields{Field: "group_power", Otype: 2, Value: 0})
	contactGroup, err = model.ActContactGroup(fromId, toId, updatesContactGroup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	dataContactGroup := schema.GetResContactGroup(contactGroup)

	toMap := make(map[string]interface{})
	toMap["contactGroup"] = dataContactGroup
	toMapStr, _ := json.Marshal(toMap)
	go server.UserGroupNoticeMsg(toId, string(toMapStr), server.MSG_MEDIA_CONTACT_GROUP_UPDATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataContactGroup,
	})
}
