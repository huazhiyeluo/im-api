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
	"sync"
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
		ownGroups, err := model.GetGroupByOwnerUid(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
			return
		}

		ownGroupIds := []uint64{}
		for _, group := range ownGroups {
			ownGroupIds = append(ownGroupIds, group.GroupId)
		}
		groupApplys, err := model.GetGroupApplyList(uid, ownGroupIds)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
			return
		}
		for _, v := range groupApplys {
			tempApplys = append(tempApplys, v)
			allUids = append(allUids, v.FromId)
			allGroupIds = append(allGroupIds, v.ToId)
		}
		allGroups, err := model.FindGroupByGroupIds(allGroupIds)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
			return
		}
		for _, v := range allGroups {
			tempAllGroups[v.GroupId] = v
		}
	}

	if utils.IsContainUint32(ttype, []uint32{0, 1}) {
		friendApplys, err := model.GetFriendApplyList(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
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
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "操作错误"})
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
			fromContactFriendData := &model.ContactFriend{
				FromId:        apply.FromId,
				ToId:          apply.ToId,
				FriendGroupId: 0,
				Level:         1,
				Remark:        apply.Remark,
				JoinTime:      nowtime,
			}
			fromContactFriend, err := model.CreateContactFriend(fromContactFriendData)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", fromContactFriend))
				c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "操作错误"})
				return
			}
			toContactFriendData := &model.ContactFriend{
				FromId:        apply.ToId,
				ToId:          apply.FromId,
				FriendGroupId: 0,
				Level:         1,
				Remark:        "",
				JoinTime:      nowtime,
			}
			toContactFriend, err := model.CreateContactFriend(toContactFriendData)
			if err != nil {
				log.Logger.Info(fmt.Sprintf("%v", toContactFriend))
				c.JSON(http.StatusOK, gin.H{"code": 5, "msg": "操作错误"})
				return
			}
			fromUser, _ := model.FindUserByUid(apply.FromId)
			toUser, _ := model.FindUserByUid(apply.ToId)

			tempApply := schema.GetResApplyUser(apply, fromUser, toUser)

			var wg sync.WaitGroup
			wg.Add(2)

			//1、告诉请求的人消息
			go func() {
				defer wg.Done()
				fromMap := make(map[string]interface{})
				fromMap["apply"] = tempApply
				fromMap["user"] = schema.GetResUser(toUser)
				fromMap["contactFriend"] = schema.GetResContactFriend(fromContactFriend)
				fromMapStr, _ := json.Marshal(fromMap)
				go server.UserFriendNoticeMsg(apply.ToId, apply.FromId, string(fromMapStr), server.MSG_MEDIA_FRIEND_AGREE)
			}()
			//2、告诉收的人消息
			go func() {
				defer wg.Done()
				toMap := make(map[string]interface{})
				toMap["apply"] = tempApply
				toMap["user"] = schema.GetResUser(fromUser)
				toMap["contactFriend"] = schema.GetResContactFriend(toContactFriend)
				toMapStr, _ := json.Marshal(toMap)
				go server.UserFriendNoticeMsg(apply.FromId, apply.ToId, string(toMapStr), server.MSG_MEDIA_FRIEND_AGREE)
			}()
			wg.Wait()
			go server.CreateMsg(&server.Message{FromId: apply.FromId, ToId: apply.ToId, MsgType: server.MSG_TYPE_SINGLE, MsgMedia: server.MSG_MEDIA_TEXT, Content: &server.MessageContent{Data: apply.Reason}})
		}
		if apply.Type == 2 {
			fromContactGroupData := &model.ContactGroup{
				FromId:   apply.FromId,
				ToId:     apply.ToId,
				Level:    1,
				Remark:   apply.Remark,
				Nickname: "",
				JoinTime: nowtime,
			}
			fromContactGroup, err := model.CreateContactGroup(fromContactGroupData)
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

			var wg sync.WaitGroup
			wg.Add(3)
			//1、告诉请求的人消息
			go func() {
				defer wg.Done()
				fromMap := make(map[string]interface{})
				fromMap["apply"] = tempApply
				fromMap["group"] = schema.GetResGroup(group)
				fromMap["contactGroup"] = schema.GetResContactGroup(fromContactGroup)
				fromMapStr, _ := json.Marshal(fromMap)
				go server.UserFriendNoticeMsg(group.OwnerUid, apply.FromId, string(fromMapStr), server.MSG_MEDIA_GROUP_AGREE)
			}()

			//2、告诉管理员消息
			go func() {
				defer wg.Done()
				toMap := make(map[string]interface{})
				toMap["apply"] = tempApply
				toMapStr, _ := json.Marshal(toMap)
				go server.UserFriendNoticeMsg(apply.FromId, group.OwnerUid, string(toMapStr), server.MSG_MEDIA_GROUP_AGREE)
			}()

			//3、告诉群的人消息
			go func() {
				defer wg.Done()
				toGroupMap := make(map[string]interface{})
				toGroupMap["user"] = schema.GetResUser(fromUser)
				toGroupMap["contactGroup"] = schema.GetResContactGroup(fromContactGroup)
				toGroupMapStr, _ := json.Marshal(toGroupMap)
				go server.UserGroupNoticeMsg(apply.ToId, string(toGroupMapStr), server.MSG_MEDIA_GROUP_AGREE)
			}()
			wg.Wait()
			go server.CreateMsg(&server.Message{FromId: apply.FromId, ToId: apply.ToId, MsgType: server.MSG_TYPE_ROOM, MsgMedia: server.MSG_MEDIA_TEXT, Content: &server.MessageContent{Data: apply.Info}})
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
