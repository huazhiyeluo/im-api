package schema

type EditGroup struct {
	GroupId  uint64 // 组ID
	OwnerUid uint64 // 拥有者
	Type     uint32 // 类型
	Name     string // 名称
	Icon     string // 图标
	Info     string // 描述
}
