package schema

import (
	"qqapi/internal/model"
	"qqapi/internal/server"
)

type ResContactFriendGroup struct {
	FriendGroupId uint32 `json:"friendGroupId"` // UID
	OwnerUid      uint64 `json:"ownerUid"`      // 用户名
	Name          string `json:"name"`          // 联系人组名
}

type ResContactFriend struct {
	FromId        uint64 `json:"fromId"`        // UID
	ToId          uint64 `json:"toId"`          // UID
	FriendGroupId uint32 `json:"friendGroupId"` // 用户组ID 0 默认分组
	Level         uint32 `json:"level"`         // 用户亲密度
	Remark        string `json:"remark"`        // 备注
	Desc          string `json:"desc"`          // 备注
	IsTop         uint32 `json:"isTop"`         // 是否置顶 0否1是
	IsHidden      uint32 `json:"isHidden"`      // 是否隐藏 0否1是
	IsQuiet       uint32 `json:"isQuiet"`       // 是否免打扰 0否1是
	JoinTime      int64  `json:"joinTime"`      // 加好友时间
	IsOnline      uint32 `json:"isOnline"`      // 是否在线
}

/********************************对接口********************************/

// 好友联系人
func GetResContactFriend(contact *model.ContactFriend) *ResContactFriend {
	onlines := server.CheckUserOnlineStatus([]uint64{contact.ToId})
	tempFriend := &ResContactFriend{
		FromId:        contact.FromId,
		ToId:          contact.ToId,
		FriendGroupId: contact.FriendGroupId,
		Level:         contact.Level,
		Remark:        contact.Remark,
		Desc:          contact.Desc,
		IsTop:         contact.IsTop,
		IsHidden:      contact.IsHidden,
		IsQuiet:       contact.IsQuiet,
		JoinTime:      contact.JoinTime,
		IsOnline:      onlines[contact.ToId],
	}
	return tempFriend
}
