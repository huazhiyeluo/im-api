package service

import (
	"imapi/internal/model"
	"imapi/internal/schema"
	"imapi/internal/server"
)

// ---------------------------------------------------------------- inner ----------------------------------------------------------------

func getResApplyUser(apply *model.Apply, fromUser *model.User, toUser *model.User) *schema.ResApply {
	tempApply := &schema.ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Username
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toUser.Username
	tempApply.ToIcon = toUser.Avatar
	return tempApply
}

func getResApplyGroup(apply *model.Apply, fromUser *model.User, toGroup *model.Group) *schema.ResApply {
	tempApply := &schema.ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Username
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toGroup.Name
	tempApply.ToIcon = toGroup.Icon
	return tempApply
}

func getResUser(user *model.User) *schema.ResFriend {
	onlines := server.CheckUserOnlineStatus([]uint64{user.Uid})
	tempFriend := &schema.ResFriend{Uid: user.Uid, Username: user.Username, Avatar: user.Avatar, IsOnline: onlines[user.Uid]}
	return tempFriend
}

func getResGroup(group *model.Group) *schema.ResGroup {
	tempGroup := &schema.ResGroup{GroupId: group.GroupId, OwnerUid: group.OwnerUid, Name: group.Name, Icon: group.Icon, Info: group.Info, Num: group.Num}
	return tempGroup
}
