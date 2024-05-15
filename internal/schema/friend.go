package schema

import (
	"imapi/internal/model"
	"imapi/internal/server"
)

type ResContactGroup struct {
	FriendGroupId uint32 `json:"friendGroupId"` // UID
	OwnerUid      uint64 `json:"ownerUid"`      // 用户名
	Name          string `json:"name"`          // 联系人组名
}

type ResFriend struct {
	Uid           uint64 `json:"uid"`           // UID
	Username      string `json:"username"`      // 昵称
	Email         string `json:"email"`         // 邮箱
	Phone         string `json:"phone"`         // 手机号
	Avatar        string `json:"avatar"`        // 头像
	Sex           uint32 `json:"sex"`           // 性别： 0 未知 1男 2女
	Birthday      int64  `json:"birthday"`      // 生日
	Info          string `json:"info"`          // 简介
	Exp           uint32 `json:"exp"`           // 用户经验
	CreateTime    int64  `json:"createTime"`    // 创建时间
	FriendGroupId uint32 `json:"friendGroupId"` // 用户组ID 0 默认分组
	Level         uint32 `json:"level"`         // 用户亲密度
	Remark        string `json:"remark"`        // 备注
	IsTop         uint32 `json:"isTop"`         // 是否置顶 0否1是
	IsHidden      uint32 `json:"isHidden"`      // 是否隐藏 0否1是
	IsQuiet       uint32 `json:"isQuiet"`       // 是否免打扰 0否1是
	JoinTime      int64  `json:"joinTime"`      // 加好友时间
	IsOnline      uint32 `json:"isOnline"`      // 是否在线
}

/********************************对接口********************************/

// 好友和好友备注组合
func GetResFriend(user *model.User, contact *model.ContactUser) *ResFriend {
	onlines := server.CheckUserOnlineStatus([]uint64{user.Uid})
	tempFriend := &ResFriend{
		Uid:           user.Uid,
		Username:      user.Username,
		Email:         user.Email,
		Phone:         user.Phone,
		Avatar:        user.Avatar,
		Sex:           user.Sex,
		Birthday:      user.Birthday,
		Info:          user.Info,
		Exp:           user.Exp,
		CreateTime:    user.CreateTime,
		FriendGroupId: contact.FriendGroupId,
		Level:         contact.Level,
		Remark:        contact.Remark,
		IsTop:         contact.IsTop,
		IsHidden:      contact.IsHidden,
		IsQuiet:       contact.IsQuiet,
		JoinTime:      contact.JoinTime,
		IsOnline:      onlines[user.Uid],
	}
	return tempFriend
}
