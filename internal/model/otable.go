package model

// 用户表
type User struct {
	Uid      uint64 `gorm:"column:uid;primary_key;AUTO_INCREMENT"` // UID
	Username string `gorm:"column:username;NOT NULL"`              // 用户名
	Password string `gorm:"column:password;NOT NULL"`              // 密码
	Info     string `gorm:"column:info;NOT NULL"`                  // 简介
	Email    string `gorm:"column:email;NOT NULL"`                 // 邮箱
	Phone    string `gorm:"column:phone;NOT NULL"`                 // 手机号
	Avatar   string `gorm:"column:avatar;NOT NULL"`                // 头像
	IsLogout uint32 `gorm:"column:is_logout;default:0;NOT NULL"`   // 是否退出登录0否 1是
}

func (m *User) TableName() string {
	return "qqim.user"
}

// 群管理表
type Group struct {
	GroupId  uint64 `gorm:"column:group_id;primary_key;AUTO_INCREMENT"` // ID
	OwnerUid uint64 `gorm:"column:owner_uid;default:0"`                 // 拥有者
	Type     uint32 `gorm:"column:type;default:0"`                      // 群类型
	Name     string `gorm:"column:name;NOT NULL"`                       // 名称
	Icon     string `gorm:"column:icon;NOT NULL"`                       // 图标
	Info     string `gorm:"column:info;NOT NULL"`                       // 描述
	Num      uint32 `gorm:"column:num;default:0"`                       // 群人数
}

func (m *Group) TableName() string {
	return "qqim.group"
}

// 联系人表
type Contact struct {
	Id             uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT"` // ID
	FromId         uint64 `gorm:"column:from_id;default:0"`             // ID [主]
	ToId           uint64 `gorm:"column:to_id;default:0"`               // ID [从]
	Type           uint32 `gorm:"column:type;default:0"`                // 联系人类型 1用户 2群
	ContactGroupId uint32 `gorm:"column:contact_group_id;default:0"`    // 联系人类型 1用户 2群
	Remark         string `gorm:"column:remark;NOT NULL"`               // 备注
}

func (m *Contact) TableName() string {
	return "qqim.contact"
}

// 联系人分组表
type ContactGroup struct {
	ContactGroupId uint32 `gorm:"column:contact_group_id;primary_key;AUTO_INCREMENT"` // ID
	OwnerUid       uint64 `gorm:"column:owner_uid;default:0"`                         // 拥有者
	Name           string `gorm:"column:name;NOT NULL"`                               // 名称
}

func (m *ContactGroup) TableName() string {
	return "qqim.contact_group"
}

// 消息表
type Message struct {
	Id         string `gorm:"column:id;NOT NULL;primary_key;"` // ID
	FromId     uint64 `gorm:"column:from_id;default:0"`        // ID [主]
	ToId       uint64 `gorm:"column:to_id;default:0"`          // ID [从]
	MsgType    uint32 `gorm:"column:msg_type;default:0"`       // 消息类型 1私信 2群 3广播
	MsgMedia   uint32 `gorm:"column:msg_media;default:0"`      // 图片类型 1文字 2图片 3 音频 4 视频
	Content    string `gorm:"column:content"`                  // 内容
	CreateTime int64  `gorm:"column:create_time;default:0"`    // 创建时间
	Status     uint32 `gorm:"column:status;default:0"`         // 状态
}

func (m *Message) TableName() string {
	return "qqim.message"
}

// 消息表-未读
type MessageUnread struct {
	Uid        uint64 `gorm:"column:uid;default:0"`         // 目标UID
	MsgId      string `gorm:"column:msg_id;NOT NULL"`       // ID [主]
	CreateTime int64  `gorm:"column:create_time;default:0"` // 创建时间
}

func (m *MessageUnread) TableName() string {
	return "qqim.message_unread"
}

// 申请联系人表
type Apply struct {
	Id          uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT"` // ID
	FromId      uint64 `gorm:"column:from_id;default:0;NOT NULL"`    // ID [主]
	ToId        uint64 `gorm:"column:to_id;default:0;NOT NULL"`      // ID  [从]
	Type        uint32 `gorm:"column:type;default:0;NOT NULL"`       // 联系人类型 1用户 2群
	Reason      string `gorm:"column:reason;NOT NULL"`               // 原因
	Status      uint32 `gorm:"column:status;default:0;NOT NULL"`     // 状态 0默认 1同意 2拒绝
	OperateTime int64  `gorm:"column:operate_time;default:0"`        // 创建时间
}

func (m *Apply) TableName() string {
	return "qqim.apply"
}
