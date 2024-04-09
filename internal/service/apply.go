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
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func GetApplyList(c *gin.Context) {

	data := make(map[string]interface{})
	c.Bind(&data)

	if _, ok := data["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": 100, "message": "UID不存在"})
		return
	}

	uid := uint64(utils.ToNumber(data["uid"]))
	ownGroups, err := model.GetGroupByOwnerUid(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	ownGroupIds := []uint64{}
	for _, group := range ownGroups {
		ownGroupIds = append(ownGroupIds, group.GroupId)
	}

	friendApplys, err := model.GetFriendApplyList(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	groupApplys, err := model.GetGroupApplyList(uid, ownGroupIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}

	allUids := []uint64{}
	allGroupIds := []uint64{}
	var tempApplys []*model.Apply
	for _, v := range friendApplys {
		tempApplys = append(tempApplys, v)
		allUids = append(allUids, v.FromId)
		allUids = append(allUids, v.ToId)
	}
	for _, v := range groupApplys {
		tempApplys = append(tempApplys, v)
		allUids = append(allUids, v.FromId)
		allGroupIds = append(allGroupIds, v.ToId)
	}

	allUsers, err := model.FindUserByUids(allUids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempAllUsers := make(map[uint64]*model.User)
	for _, v := range allUsers {
		tempAllUsers[v.Uid] = v
	}
	allGroups, err := model.FindGroupByGroupIds(allGroupIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
		return
	}
	tempAllGroups := make(map[uint64]*model.Group)
	for _, v := range allGroups {
		tempAllGroups[v.GroupId] = v
	}

	var applys []*schema.ResApply
	for _, v := range tempApplys {
		if v.Type == 1 {
			temp := getResApplyUser(v, tempAllUsers[v.FromId], tempAllUsers[v.ToId])
			applys = append(applys, temp)
		}
		if v.Type == 2 {
			temp := getResApplyGroup(v, tempAllUsers[v.FromId], tempAllGroups[v.ToId])
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
	apply.OperateTime = time.Now().Unix()
	//同意
	if status == 1 {
		apply.Status = 1
		apply, err = model.UpdateApply(apply)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", apply))
			c.JSON(http.StatusOK, gin.H{"code": 3, "msg": "操作错误"})
			return
		}

		insertContachData := &model.Contact{
			FromId: apply.FromId,
			ToId:   apply.ToId,
			Type:   apply.Type,
			Remark: "",
		}
		contact, err := model.CreateContact(insertContachData)
		if err != nil {
			log.Logger.Info(fmt.Sprintf("%v", contact))
			c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
			return
		}

		if apply.Type == 1 {
			insertFriendContactData := &model.Contact{
				FromId: apply.ToId,
				ToId:   apply.FromId,
				Type:   apply.Type,
				Remark: "",
			}
			contact, err = model.CreateContact(insertFriendContactData)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", contact))
				c.JSON(http.StatusOK, gin.H{"code": 5, "msg": "操作错误"})
				return
			}
			fromUser, _ := model.FindUserByUid(apply.FromId)
			toUser, _ := model.FindUserByUid(apply.ToId)

			tempApply := getResApplyUser(apply, fromUser, toUser)
			//1、告诉请求的人消息
			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMap["user"] = getResUser(toUser)
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(apply.ToId, apply.FromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_AGREE)

			//2、告诉收的人消息
			toMap := make(map[string]interface{})
			toMap["apply"] = tempApply
			toMap["user"] = getResUser(fromUser)
			toMapStr, _ := json.Marshal(toMap)
			go server.UserFriendNoticeMsg(apply.FromId, apply.ToId, string(toMapStr), server.MSG_MEDIA_FRIEND_AGREE)
		}
		if apply.Type == 2 {
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

			tempApply := getResApplyGroup(apply, fromUser, group)

			//1、告诉请求的人消息
			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMap["user"] = getResUser(fromUser)
			fromMap["group"] = getResGroup(group)
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_AGREE)

			//2、告诉管理员消息
			toMap := make(map[string]interface{})
			toMap["apply"] = tempApply
			toMapStr, _ := json.Marshal(toMap)
			go server.UserFriendNoticeMsg(apply.FromId, group.OwnerUid, string(toMapStr), server.MSG_MEDIA_GROUP_AGREE)

			//3、告诉群的人消息
			toGroupMap := make(map[string]interface{})
			toGroupMap["user"] = getResUser(fromUser)
			toGroupMap["group"] = getResGroup(group)
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
			tempApply := getResApplyUser(apply, fromUser, toUser)

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

			tempApply := getResApplyGroup(apply, fromUser, group)

			fromMap := make(map[string]interface{})
			fromMap["apply"] = tempApply
			fromMapStr, _ := json.Marshal(fromMap)
			go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_REFUSE)

			toMap := make(map[string]interface{})
			toMap["apply"] = tempApply
			toMapStr, _ := json.Marshal(toMap)
			go server.UserFriendNoticeMsg(apply.FromId, group.OwnerUid, string(toMapStr), server.MSG_MEDIA_GROUP_REFUSE)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})

}
