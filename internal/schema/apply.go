package schema

import "imapi/internal/model"

type ResApply struct {
	Id          uint32 `json:"id"`
	FromId      uint64 `json:"fromId"`
	ToId        uint64 `json:"toId"`
	Type        uint32 `json:"type"`
	FromName    string `json:"fromName"`
	FromIcon    string `json:"fromIcon"`
	ToName      string `json:"toName"`
	ToIcon      string `json:"toIcon"`
	Reason      string `json:"reason"`
	Status      uint32 `json:"status"`
	OperateTime int64  `json:"operateTime"`
}

func GetResApplyUser(apply *model.Apply, fromUser *model.User, toUser *model.User) *ResApply {
	tempApply := &ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Username
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toUser.Username
	tempApply.ToIcon = toUser.Avatar
	return tempApply
}

func GetResApplyGroup(apply *model.Apply, fromUser *model.User, toGroup *model.Group) *ResApply {
	tempApply := &ResApply{Id: apply.Id, FromId: apply.FromId, ToId: apply.ToId, Type: apply.Type, Reason: apply.Reason, Status: apply.Status, OperateTime: apply.OperateTime}
	tempApply.FromName = fromUser.Username
	tempApply.FromIcon = fromUser.Avatar
	tempApply.ToName = toGroup.Name
	tempApply.ToIcon = toGroup.Icon
	return tempApply
}
