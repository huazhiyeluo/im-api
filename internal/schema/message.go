package schema

type ChatMsg struct {
	FromId  uint64 `json:"fromId"`  // 发送者
	ToId    uint64 `json:"toId"`    // 接受者
	MsgType uint32 `json:"msgType"` // 1单聊消息 2群聊消息
	Start   int64  `json:"start"`   // 开始
	End     int64  `json:"end"`     // 结束
	IsRev   uint32 `json:"isRev"`   // 是否正序
}
