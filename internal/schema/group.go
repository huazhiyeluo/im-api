package schema

import "qqapi/internal/model"

type EditGroup struct {
	GroupId  uint64 `json:"groupId"`  // 组ID
	OwnerUid uint64 `json:"ownerUid"` // 拥有者
	Type     uint32 `json:"type"`     // 类型
	Name     string `json:"name"`     // 名称
	Icon     string `json:"icon"`     // 图标
	Info     string `json:"info"`     // 描述
}

// 群表
type ResGroup struct {
	GroupId    uint64 `json:"groupId"`    // 组ID
	OwnerUid   uint64 `json:"ownerUid"`   // 拥有者
	Type       uint32 `json:"type"`       // 类型
	Name       string `json:"name"`       // 名称
	Icon       string `json:"icon"`       // 图标
	Info       string `json:"info"`       // 描述
	Num        uint32 `json:"num"`        // 群人数
	Exp        uint32 `json:"exp"`        // 群经验
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

func GetResGroup(m *model.Group) *ResGroup {
	return &ResGroup{
		GroupId:    m.GroupId,
		OwnerUid:   m.OwnerUid,
		Type:       m.Type,
		Name:       m.Name,
		Icon:       m.Icon,
		Info:       m.Info,
		Num:        m.Num,
		Exp:        m.Exp,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}
