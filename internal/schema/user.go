package schema

type Login struct {
	Username string // 用户名
	Password string // 密码
}

type Register struct {
	Username   string // 用户名
	Password   string // 密码
	Repassword string // 确认密码
}

type ChatMsg struct {
	FromId  uint64 // 发送者
	ToId    uint64 // 接受者
	MsgType uint32 // 1单聊消息 2群聊消息
	Start   int64  // 开始
	End     int64  // 结束
	IsRev   uint32 // 是否正序
}

type EditUser struct {
	Uid      uint64 // UID
	Username string // 用户名
	Info     string // 简介
	Avatar   string //头像
}
