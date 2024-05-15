package model

// 1、用户表
type User struct {
	Uid        uint64 `gorm:"column:uid;primary_key;AUTO_INCREMENT"` // UID
	Username   string `gorm:"column:username;NOT NULL"`              // 昵称
	Password   string `gorm:"column:password;NOT NULL"`              // 密码
	Email      string `gorm:"column:email;NOT NULL"`                 // 邮箱
	Phone      string `gorm:"column:phone;NOT NULL"`                 // 手机号
	Avatar     string `gorm:"column:avatar;NOT NULL"`                // 头像
	Sex        uint32 `gorm:"column:sex;NOT NULL"`                   // 性别： 0 未知 1男 2女
	Birthday   int64  `gorm:"column:birthday;NOT NULL"`              // 生日
	Info       string `gorm:"column:info;NOT NULL"`                  // 简介
	Exp        uint32 `gorm:"column:exp;default:0;NOT NULL"`         // 用户经验
	CreateTime int64  `gorm:"column:create_time;NOT NULL"`           // 创建时间
	UpdateTime int64  `gorm:"column:update_time;NOT NULL"`           // 更新时间
}

func (m *User) TableName() string {
	return "qqim.user"
}

// 2、好友组表
type FriendGroup struct {
	FriendGroupId uint32 `gorm:"column:friend_group_id;primary_key;AUTO_INCREMENT"` // 用户组ID
	OwnerUid      uint64 `gorm:"column:owner_uid;default:0;NOT NULL"`               // 拥有者
	Name          string `gorm:"column:name;NOT NULL"`                              // 用户组名
	CreateTime    int64  `gorm:"column:create_time;default:0;NOT NULL"`             // 创建时间
}

func (m *FriendGroup) TableName() string {
	return "qqim.friend_group"
}

// 3-1、群管理表
type Group struct {
	GroupId    uint64 `gorm:"column:group_id;primary_key;AUTO_INCREMENT"` // ID
	OwnerUid   uint64 `gorm:"column:owner_uid;default:0;NOT NULL"`        // 创建人
	Type       uint32 `gorm:"column:type;default:0;NOT NULL"`             // 群类型
	Name       string `gorm:"column:name;NOT NULL"`                       // 名称
	Icon       string `gorm:"column:icon;NOT NULL"`                       // 图标
	Info       string `gorm:"column:info;NOT NULL"`                       // 描述
	Num        uint32 `gorm:"column:num;default:0;NOT NULL"`              // 群人数
	Exp        uint32 `gorm:"column:exp;default:0;NOT NULL"`              // 群经验
	CreateTime int64  `gorm:"column:create_time;default:0;NOT NULL"`      // 创建时间
	UpdateTime int64  `gorm:"column:update_time;default:0;NOT NULL"`      // 更新时间
}

func (m *Group) TableName() string {
	return "qqim.group"
}

// 3-2 群公告
type GroupTips struct {
	Id         uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT"`  // ID
	GroupId    uint64 `gorm:"column:group_id;default:0;NOT NULL"`    // 群ID
	Content    string `gorm:"column:content;NOT NULL"`               // 群公告
	CreateTime int64  `gorm:"column:create_time;default:0;NOT NULL"` // 创建时间
}

func (m *GroupTips) TableName() string {
	return "qqim.group_tips"
}

// 4、好友联系人表
type ContactUser struct {
	FromId        uint64 `gorm:"column:from_id;default:0;NOT NULL"`         // ID [主]
	ToId          uint64 `gorm:"column:to_id;default:0;NOT NULL"`           // ID  [从]
	FriendGroupId uint32 `gorm:"column:friend_group_id;default:0;NOT NULL"` // 用户组ID 0 默认分组
	Level         uint32 `gorm:"column:level;default:0;NOT NULL"`           // 用户亲密度
	Remark        string `gorm:"column:remark;NOT NULL"`                    // 备注
	IsTop         uint32 `gorm:"column:is_top;default:0;NOT NULL"`          // 是否置顶 0否1是
	IsHidden      uint32 `gorm:"column:is_hidden;default:0"`                // 是否隐藏 0否1是
	IsQuiet       uint32 `gorm:"column:is_quiet;default:0"`                 // 是否免打扰 0否1是
	JoinTime      int64  `gorm:"column:join_time;default:0;NOT NULL"`       // 加好友时间
}

func (m *ContactUser) TableName() string {
	return "qqim.contact_user"
}

// 5、组联系人表
type ContactGroup struct {
	FromId     uint64 `gorm:"column:from_id;default:0;NOT NULL"`     // ID [主]
	ToId       uint64 `gorm:"column:to_id;default:0;NOT NULL"`       // ID  [从]
	GroupPower uint32 `gorm:"column:group_power;default:0;NOT NULL"` // 群权限（0 普通 1管理员 2创始人）
	Level      uint32 `gorm:"column:level;default:0;NOT NULL"`       // 我在本群等级
	Remark     string `gorm:"column:remark;NOT NULL"`                // 群聊备注
	Nickname   string `gorm:"column:nickname;NOT NULL"`              // 我在本群昵称
	IsTop      uint32 `gorm:"column:is_top;default:0;NOT NULL"`      // 是否置顶 0否1是
	IsHidden   uint32 `gorm:"column:is_hidden;default:0"`            // 是否隐藏 0否1是
	IsQuiet    uint32 `gorm:"column:is_quiet;default:0"`             // 是否免打扰 0否1是
	JoinTime   int64  `gorm:"column:join_time;default:0;NOT NULL"`   // 入群时间
}

func (m *ContactGroup) TableName() string {
	return "qqim.contact_group"
}

// 6、消息表
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

// 7、消息表-未读
type MessageUnread struct {
	Uid        uint64 `gorm:"column:uid;default:0"`         // 目标UID
	MsgId      string `gorm:"column:msg_id;NOT NULL"`       // ID [主]
	CreateTime int64  `gorm:"column:create_time;default:0"` // 创建时间
}

func (m *MessageUnread) TableName() string {
	return "qqim.message_unread"
}

// 8.申请联系人表
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
