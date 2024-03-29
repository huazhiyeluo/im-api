package schema

type ResFriend struct {
	Uid      uint64 // UID
	Username string // 用户名
	Avatar   string // 头像
	IsOnline bool   // 是否在线
}

type ResGroup struct {
	GroupId  uint64 // 群ID
	OwnerUid uint64 //	创建者
	Name     string // 用户名
	Icon     string // 头像
	Info     string //	简介
	Num      uint32 // 人数
}
