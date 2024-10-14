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
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func GetApplyList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "UID不存在"})
		return
	}
	uid := uint64(utils.ToNumber(data["uid"]))
	var ttype uint32 = 0
	if _, ok := data["type"]; ok {
		ttype = uint32(utils.ToNumber(data["type"]))
	}

	var tempApplys []*model.Apply
	allUids := []uint64{}
	allGroupIds := []uint64{}
	tempAllUsers := make(map[uint64]*model.User)
	tempAllGroups := make(map[uint64]*model.Group)
	if utils.IsContainUint32(ttype, []uint32{0, 2}) {
		managerContactGroups, err := model.GetContactGroupManagerList(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}

		managerGroupIds := []uint64{}
		for _, contactGroup := range managerContactGroups {
			managerGroupIds = append(managerGroupIds, contactGroup.ToId)
		}
		groupApplys, err := model.GetGroupApplyList(uid, managerGroupIds)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
		for _, v := range groupApplys {
			tempApplys = append(tempApplys, v)
			allUids = append(allUids, v.FromId)
			allGroupIds = append(allGroupIds, v.ToId)
		}
		allGroups, err := model.FindGroupByGroupIds(allGroupIds)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
		for _, v := range allGroups {
			tempAllGroups[v.GroupId] = v
		}
	}

	if utils.IsContainUint32(ttype, []uint32{0, 1}) {
		friendApplys, err := model.GetFriendApplyList(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
			return
		}
		for _, v := range friendApplys {
			tempApplys = append(tempApplys, v)
			allUids = append(allUids, v.FromId)
			allUids = append(allUids, v.ToId)
		}
	}

	allUsers, err := model.FindUserByUids(allUids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	for _, v := range allUsers {
		tempAllUsers[v.Uid] = v
	}

	var applys []*schema.ResApply
	for _, v := range tempApplys {
		if v.Type == 1 {
			temp := schema.GetResApplyUser(v, tempAllUsers[v.FromId], tempAllUsers[v.ToId])
			applys = append(applys, temp)
		}
		if v.Type == 2 {
			temp := schema.GetResApplyGroup(v, tempAllUsers[v.FromId], tempAllGroups[v.ToId])
			applys = append(applys, temp)
		}
	}
	sort.SliceStable(applys, func(i, j int) bool {
		if applys[i].Id > applys[j].Id {
			return true
		} else if applys[i].Id < applys[j].Id {
			return false
		}
		return true
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": applys,
	})
}

func OperateApply(c *gin.Context) {
	data := make(map[string]interface{})
	c.Bind(&data)
	id := uint32(utils.ToNumber(data["id"]))
	status := uint32(utils.ToNumber(data["status"]))
	apply, err := model.FindApplyById(id)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v", apply))
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if apply.Id == 0 || apply.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "申请状态错误"})
		return
	}

	nowtime := time.Now().Unix()

	apply.OperateTime = nowtime
	//同意
	if status == 1 {
		apply.Status = 1
		apply, err = model.UpdateApply(apply)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", apply))
			c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
			return
		}
		if apply.Type == 1 {
			//1、主数据
			var updatesFromContactFriend []*model.Fields
			updatesFromContactFriend = append(updatesFromContactFriend, &model.Fields{Field: "friend_group_id", Otype: 2, Value: apply.FriendGroupId})
			updatesFromContactFriend = append(updatesFromContactFriend, &model.Fields{Field: "level", Otype: 2, Value: 1})
			updatesFromContactFriend = append(updatesFromContactFriend, &model.Fields{Field: "remark", Otype: 2, Value: apply.Remark})
			updatesFromContactFriend = append(updatesFromContactFriend, &model.Fields{Field: "join_time", Otype: 2, Value: nowtime})
			fromContactFriend, err := model.ActContactFriend(apply.FromId, apply.ToId, updatesFromContactFriend)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", fromContactFriend))
				c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
				return
			}

			//1、从数据
			defaultFriendGroup, err := model.GetFriendGroupByIsDefault(apply.ToId)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", defaultFriendGroup))
				c.JSON(http.StatusOK, gin.H{"code": 5, "msg": "操作错误"})
				return
			}
			var updatesToContactFriend []*model.Fields
			updatesToContactFriend = append(updatesToContactFriend, &model.Fields{Field: "friend_group_id", Otype: 2, Value: defaultFriendGroup.FriendGroupId})
			updatesToContactFriend = append(updatesToContactFriend, &model.Fields{Field: "level", Otype: 2, Value: 1})
			updatesToContactFriend = append(updatesToContactFriend, &model.Fields{Field: "remark", Otype: 2, Value: ""})
			updatesToContactFriend = append(updatesToContactFriend, &model.Fields{Field: "join_time", Otype: 2, Value: nowtime})
			toContactFriend, err := model.ActContactFriend(apply.ToId, apply.FromId, updatesToContactFriend)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", toContactFriend))
				c.JSON(http.StatusOK, gin.H{"code": 5, "msg": "操作错误"})
				return
			}

			fromUser, _ := model.FindUserByUid(apply.FromId)
			toUser, _ := model.FindUserByUid(apply.ToId)
			tempApply := schema.GetResApplyUser(apply, fromUser, toUser)

			//1、告诉请求的人消息
			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMap["user"] = schema.GetResUser(toUser)
			fromMap["contactFriend"] = schema.GetResContactFriend(fromContactFriend)
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(apply.ToId, apply.FromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_AGREE)

			//2、告诉收的人消息
			toMap := make(map[string]interface{})
			toMap["apply"] = tempApply
			toMap["user"] = schema.GetResUser(fromUser)
			toMap["contactFriend"] = schema.GetResContactFriend(toContactFriend)
			toMapStr, _ := json.Marshal(toMap)
			go server.UserFriendNoticeMsg(apply.FromId, apply.ToId, string(toMapStr), server.MSG_MEDIA_FRIEND_AGREE)
		}
		if apply.Type == 2 {
			var updatesFromContactGroup []*model.Fields
			updatesFromContactGroup = append(updatesFromContactGroup, &model.Fields{Field: "level", Otype: 2, Value: 1})
			updatesFromContactGroup = append(updatesFromContactGroup, &model.Fields{Field: "remark", Otype: 2, Value: apply.Remark})
			updatesFromContactGroup = append(updatesFromContactGroup, &model.Fields{Field: "nickname", Otype: 2, Value: ""})
			updatesFromContactGroup = append(updatesFromContactGroup, &model.Fields{Field: "join_time", Otype: 2, Value: nowtime})
			fromContactGroup, err := model.ActContactGroup(apply.FromId, apply.ToId, updatesFromContactGroup)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", fromContactGroup))
				c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
				return
			}

			group, err := model.FindGroupByGroupId(apply.ToId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
				return
			}

			group.Num = group.Num + 1
			group, err = model.UpdateGroup(group)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
				return
			}
			fromUser, _ := model.FindUserByUid(apply.FromId)
			tempApply := schema.GetResApplyGroup(apply, fromUser, group)

			//1、告诉请求的人消息
			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMap["group"] = schema.GetResGroup(group)
			fromMap["contactGroup"] = schema.GetResContactGroup(fromContactGroup)
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_AGREE)

			//2、告诉管理员消息
			contactGroups, err := model.GetGroupUserManager(apply.ToId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
				return
			}
			for _, v := range contactGroups {
				toMap := make(map[string]interface{})
				toMap["apply"] = tempApply
				toMapStr, _ := json.Marshal(toMap)
				go server.UserFriendNoticeMsg(apply.FromId, v.FromId, string(toMapStr), server.MSG_MEDIA_GROUP_AGREE)
			}

			//3、告诉群的人消息
			toGroupMap := make(map[string]interface{})
			toGroupMap["user"] = schema.GetResUser(fromUser)
			toGroupMap["group"] = schema.GetResGroup(group)
			toGroupMap["contactGroup"] = schema.GetResContactGroup(fromContactGroup)
			toGroupMapStr, _ := json.Marshal(toGroupMap)
			go server.UserGroupNoticeMsg(apply.ToId, string(toGroupMapStr), server.MSG_MEDIA_GROUP_AGREE)
		}

	}

	if status == 2 {
		apply.Status = 2
		apply, err = model.UpdateApply(apply)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", apply))
			c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
			return
		}
		if apply.Type == 1 {

			fromUser, _ := model.FindUserByUid(apply.FromId)
			toUser, _ := model.FindUserByUid(apply.ToId)
			tempApply := schema.GetResApplyUser(apply, fromUser, toUser)

			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(apply.ToId, apply.FromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_REFUSE)

			toMap := make(map[string]interface{})
			toMap["apply"] = tempApply
			toMapStr, _ := json.Marshal(toMap)
			go server.UserFriendNoticeMsg(apply.FromId, apply.ToId, string(toMapStr), server.MSG_MEDIA_FRIEND_REFUSE)
		}

		if apply.Type == 2 {
			group, err := model.FindGroupByGroupId(apply.ToId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
				return
			}
			fromUser, _ := model.FindUserByUid(apply.FromId)

			tempApply := schema.GetResApplyGroup(apply, fromUser, group)

			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_REFUSE)

			contactGroups, err := model.GetGroupUserManager(apply.ToId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
				return
			}
			for _, v := range contactGroups {
				toMap := make(map[string]interface{})
				toMap["apply"] = tempApply
				toMapStr, _ := json.Marshal(toMap)
				go server.UserFriendNoticeMsg(apply.FromId, v.FromId, string(toMapStr), server.MSG_MEDIA_GROUP_REFUSE)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}
