package schema

import (
	"qqapi/internal/model"
)

// 入参
type EditUser struct {
	Uid      uint64 `json:"uid"`      // UID
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Info     string `json:"info"`     // 简介
}

// 用户表
type ResUser struct {
	Uid       uint64 `json:"uid"`       // UID
	Nickname  string `json:"nickname"`  // 昵称
	Avatar    string `json:"avatar"`    // 头像
	Sex       uint32 `json:"sex"`       // 性别： 0 未知 1男 2女
	Birthday  int64  `json:"birthday"`  // 生日
	Info      string `json:"info"`      // 简介
	Exp       uint32 `json:"exp"`       // 用户经验
	RegTime   int64  `json:"regTime"`   // 注册时间
	LoginTime int64  `json:"loginTime"` // 最后登录时间
	Username  string `json:"username"`  // 用户名
	Email     string `json:"email"`     // 邮箱
	Phone     string `json:"phone"`     // 手机号
}

/********************************对接口********************************/

func GetResUser(m *model.User) *ResUser {
	return &ResUser{
		Uid:       m.Uid,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		Sex:       m.Sex,
		Birthday:  m.Birthday,
		Info:      m.Info,
		Exp:       m.Exp,
		RegTime:   m.RegTime,
		LoginTime: m.LoginTime,
	}
}
