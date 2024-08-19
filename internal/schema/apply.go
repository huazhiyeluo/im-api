package schema

import "qqapi/internal/model"

type ResApply struct {
	Id            uint32 `json:"id"`
	FromId        uint64 `json:"fromId"`
	ToId          uint64 `json:"toId"`
	Type          uint32 `json:"type"`
	FromName      string `json:"fromName"`
	FromIcon      string `json:"fromIcon"`
	ToName        string `json:"toName"`
	ToIcon        string `json:"toIcon"`
	Reason        string `json:"reason"`
	Remark        string `json:"remark"`
	Info          string `json:"info"`
	FriendGroupId uint32 `json:"friendGroupId"`
	Status        uint32 `json:"status"`
	OperateTime   int64  `json:"operateTime"`
}

func GetResApplyUser(apply *model.Apply, fromUser *model.User, toUser *model.User) *ResApply {
	tempApply := &ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Remark: apply.Remark, Info: apply.Info, FriendGroupId: apply.FriendGroupId, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Nickname
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toUser.Nickname
	tempApply.ToIcon = toUser.Avatar
	return tempApply
}

func GetResApplyGroup(apply *model.Apply, fromUser *model.User, toGroup *model.Group) *ResApply {
	tempApply := &ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Remark: apply.Remark, Info: apply.Info, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Nickname
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toGroup.Name
	tempApply.ToIcon = toGroup.Icon
	return tempApply
}
