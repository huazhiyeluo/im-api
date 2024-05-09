package schema

import "imapi/internal/model"

type Login struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type Register struct {
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Repassword string `json:"repassword"` // 确认密码
}

type ChatMsg struct {
	FromId  uint64 `json:"fromId"`  // 发送者
	ToId    uint64 `json:"toId"`    // 接受者
	MsgType uint32 `json:"msgType"` // 1单聊消息 2群聊消息
	Start   int64  `json:"start"`   // 开始
	End     int64  `json:"end"`     // 结束
	IsRev   uint32 `json:"isRev"`   // 是否正序
}

type EditUser struct {
	Uid      uint64 `json:"uid"`      // UID
	Username string `json:"username"` // 用户名
	Info     string `json:"info"`     // 简介
	Avatar   string `json:"avatar"`   //头像
}

// 用户表
type User struct {
	Uid      uint64 `json:"uid"`      // UID
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Info     string `json:"info"`     // 简介
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
	Avatar   string `json:"avatar"`   // 头像
	IsLogout uint32 `json:"isLogout"` // 是否退出登录0否 1是
}

/********************************对接口********************************/

func GetUser(m *model.User) *User {
	return &User{
		Uid:      m.Uid,
		Username: m.Username,
		Password: m.Password,
		Info:     m.Info,
		Email:    m.Email,
		Phone:    m.Phone,
		Avatar:   m.Avatar,
		IsLogout: m.IsLogout,
	}
}
