package schema

import "qqapi/internal/model"

type ResContactGroup struct {
	GroupId    uint64 `json:"groupId"`    // 群ID
	OwnerUid   uint64 `json:"ownerUid"`   //创建者
	Name       string `json:"name"`       // 用户名
	Icon       string `json:"icon"`       // 头像
	Info       string `json:"info"`       // 简介
	Num        uint32 `json:"num"`        // 人数
	Exp        uint32 `json:"exp"`        // 经验值
	CreateTime int64  `json:"createTime"` // 注册时间
	GroupPower uint32 `json:"groupPower"` // 群权限（0 普通 1管理员 2创始人）
	Level      uint32 `json:"level"`      // 我在本群等级
	Remark     string `json:"remark"`     // 群聊备注
	Nickname   string `json:"nickname"`   // 我在本群昵称
	IsTop      uint32 `json:"isTop"`      // 是否置顶 0否1是
	IsHidden   uint32 `json:"isHidden"`   // 是否隐藏 0否1是
	IsQuiet    uint32 `json:"isQuiet"`    // 是否免打扰 0否1是
	JoinTime   int64  `json:"joinTime"`   // 加群时间
}

type ResContactGroupUser struct {
	Uid        uint64 `json:"uid"`        // UID
	Username   string `json:"username"`   // 昵称
	Email      string `json:"email"`      // 邮箱
	Phone      string `json:"phone"`      // 手机号
	Avatar     string `json:"avatar"`     // 头像
	Sex        uint32 `json:"sex"`        // 性别： 0 未知 1男 2女
	Birthday   int64  `json:"birthday"`   // 生日
	Info       string `json:"info"`       // 简介
	Exp        uint32 `json:"exp"`        // 用户经验
	CreateTime int64  `json:"createTime"` // 创建时间
	GroupPower uint32 `json:"groupPower"` // 群权限（0 普通 1管理员 2创始人）
	Level      uint32 `json:"level"`      // 用户亲密度
	Remark     string `json:"remark"`     // 备注
	Nickname   string `json:"nickname"`   // 备注
	IsTop      uint32 `json:"isTop"`      // 是否置顶 0否1是
	IsHidden   uint32 `json:"isHidden"`   // 是否隐藏 0否1是
	IsQuiet    uint32 `json:"isQuiet"`    // 是否免打扰 0否1是
	JoinTime   int64  `json:"joinTime"`   // 加群时间
}

// 群和群备注组合
func GetResContactGroup(group *model.Group, contact *model.ContactGroup) *ResContactGroup {
	tempGroup := &ResContactGroup{
		GroupId:    group.GroupId,
		OwnerUid:   group.OwnerUid,
		Name:       group.Name,
		Icon:       group.Icon,
		Info:       group.Info,
		Num:        group.Num,
		Exp:        group.Exp,
		CreateTime: group.CreateTime,
		GroupPower: contact.GroupPower,
		Level:      contact.Level,
		Remark:     contact.Remark,
		Nickname:   contact.Nickname,
		IsTop:      contact.IsTop,
		IsHidden:   contact.IsHidden,
		IsQuiet:    contact.IsQuiet,
		JoinTime:   contact.JoinTime,
	}
	return tempGroup
}

// 群成员和群备注组合
func GetResContactGroupUser(user *model.User, contact *model.ContactGroup) *ResContactGroupUser {
	tempGroup := &ResContactGroupUser{
		Uid:        user.Uid,
		Username:   user.Username,
		Email:      user.Email,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Sex:        user.Sex,
		Birthday:   user.Birthday,
		Info:       user.Info,
		Exp:        user.Exp,
		CreateTime: user.CreateTime,
		GroupPower: contact.GroupPower,
		Level:      contact.Level,
		Remark:     contact.Remark,
		Nickname:   contact.Nickname,
		IsTop:      contact.IsTop,
		IsHidden:   contact.IsHidden,
		IsQuiet:    contact.IsQuiet,
		JoinTime:   contact.JoinTime,
	}
	return tempGroup
}
