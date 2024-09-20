package service

import (
	"encoding/json"
	"net/http"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/server"
	"qqapi/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 1、创建群
func CreateGroup(c *gin.Context) {
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
		FromId:     data.OwnerUid,
		ToId:       group.GroupId,
		GroupPower: 2,
		Level:      1,
		Remark:     "",
		Nickname:   "",
		JoinTime:   nowtime,
	}
	contactGroup, err := model.CreateContactGroup(insertContactGroupData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	//1、告诉请求的人消息
	fromMap := make(map[string]interface{})
	fromMap["user"] = schema.GetResUser(fromUser)
	fromMap["group"] = schema.GetResGroup(group)
	fromMap["contactGroup"] = schema.GetResContactGroup(contactGroup)
	fromMapStr, _ := json.Marshal(fromMap)
	go server.UserFriendNoticeMsg(group.OwnerUid, group.OwnerUid, string(fromMapStr), server.MSG_MEDIA_GROUP_CREATE)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": schema.GetResGroup(group),
	})
}

// 2、编辑群
func ActGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["groupId"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "用户组不存在"})
		return
	}
	groupId := uint64(utils.ToNumber(data["groupId"]))

	group, err := model.FindGroupByGroupId(groupId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户组不存在"})
		return
	}
	nowtime := time.Now().Unix()

	var updatesGroup []*model.Fields
	updatesGroup = append(updatesGroup, &model.Fields{Field: "update_time", Otype: 2, Value: nowtime})
	for key, val := range data {
		newkey := utils.CamelToSnakeCase(key)
		updatesGroup = append(updatesGroup, &model.Fields{Field: newkey, Otype: 2, Value: val})
	}
	group, err = model.ActGroup(groupId, updatesGroup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	toMap := make(map[string]interface{})
	toMap["group"] = schema.GetResGroup(group)
	toMapStr, _ := json.Marshal(toMap)
	go server.GroupInfoNoticeMsg(groupId, string(toMapStr))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// 3、搜索群
func SearchGroup(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["keyword"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "条件不存在"})
		return
	}
	keyword := utils.ToString(data["keyword"])
	pageSize := uint32(utils.ToNumber(data["pageSize"]))
	if pageSize == 0 {
		pageSize = 10
	}
	pageNum := uint32(utils.ToNumber(data["pageNum"]))
	if pageNum == 0 {
		pageNum = 1
	}
	groups, count, err := model.FindGroupByKeyword(pageSize, pageNum, keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	var dataGroups []*schema.ResGroup
	for _, v := range groups {
		temp := schema.GetResGroup(v)
		dataGroups = append(dataGroups, temp)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"groups": dataGroups,
			"count":  count,
		},
	})
}
