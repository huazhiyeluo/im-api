package schema

import "qqapi/internal/model"

type ResContactGroup struct {
	FromId     uint64 `json:"fromId"`     // UID
	ToId       uint64 `json:"toId"`       // UID
	GroupPower uint32 `json:"groupPower"` // 群权限（0 普通 1管理员 2创始人）
	Level      uint32 `json:"level"`      // 我在本群等级
	Remark     string `json:"remark"`     // 群聊备注
	Nickname   string `json:"nickname"`   // 我在本群昵称
	IsTop      uint32 `json:"isTop"`      // 是否置顶 0否1是
	IsHidden   uint32 `json:"isHidden"`   // 是否隐藏 0否1是
	IsQuiet    uint32 `json:"isQuiet"`    // 是否免打扰 0否1是
	JoinTime   int64  `json:"joinTime"`   // 加群时间
}

// 群联系人
func GetResContactGroup(contact *model.ContactGroup) *ResContactGroup {
	tempGroup := &ResContactGroup{
		FromId:     contact.FromId,
		ToId:       contact.ToId,
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
