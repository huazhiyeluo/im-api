package schema

type ResContactGroup struct {
	ContactGroupId uint32 `json:"contactGroupId"` // UID
	OwnerUid       uint64 `json:"ownerUid"`       // 用户名
	Name           string `json:"name"`           // 联系人组名
}

type ResFriend struct {
	Uid            uint64 `json:"uid"`            // UID
	Username       string `json:"username"`       // 用户名
	Avatar         string `json:"avatar"`         // 头像
	IsOnline       uint32 `json:"isOnline"`       // 是否在线
	ContactGroupId uint32 `json:"contactGroupId"` //分组ID
	Remark         string `json:"remark"`         //备注
}

type ResGroup struct {
	GroupId  uint64 `json:"groupId"`  // 群ID
	OwnerUid uint64 `json:"ownerUid"` //	创建者
	Name     string `json:"name"`     // 用户名
	Icon     string `json:"icon"`     // 头像
	Info     string `json:"info"`     // 简介
	Num      uint32 `json:"num"`      // 人数
	Remark   string `json:"remark"`   // 备注
}
