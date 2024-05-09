package schema

type EditGroup struct {
	GroupId  uint64 `json:"groupId"`  // 组ID
	OwnerUid uint64 `json:"ownerUid"` // 拥有者
	Type     uint32 `json:"type"`     // 类型
	Name     string `json:"name"`     // 名称
	Icon     string `json:"icon"`     // 图标
	Info     string `json:"info"`     // 描述
}
