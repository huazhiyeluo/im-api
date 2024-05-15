package schema

import (
	"imapi/internal/model"
)

// 入参
type EditUser struct {
	Uid      uint64 `json:"uid"`      // UID
	Username string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 头像
	Info     string `json:"info"`     // 简介
}

// 用户表
type User struct {
	Uid        uint64 `json:"uid"`        // UID
	Username   string `json:"username"`   // 昵称
	Email      string `json:"email"`      // 邮箱
	Phone      string `json:"phone"`      // 手机号
	Avatar     string `json:"avatar"`     // 头像
	Sex        uint32 `json:"sex"`        // 性别： 0 未知 1男 2女
	Birthday   int64  `json:"birthday"`   // 生日
	Info       string `json:"info"`       // 简介
	Exp        uint32 `json:"exp"`        // 用户经验
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

/********************************对接口********************************/

func GetUser(m *model.User) *User {
	return &User{
		Uid:        m.Uid,
		Username:   m.Username,
		Email:      m.Email,
		Phone:      m.Phone,
		Avatar:     m.Avatar,
		Sex:        m.Sex,
		Birthday:   m.Birthday,
		Info:       m.Info,
		Exp:        m.Exp,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}
