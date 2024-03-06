package models

import (
	"demoapi/utils"
	"log"
)

// 用户表
type User struct {
	Uid        uint64 `gorm:"column:uid;primary_key;AUTO_INCREMENT"` // UID
	Username   string `gorm:"column:username"`                       // 用户名
	Password   string `gorm:"column:password"`                       // 密码
	Email      string `gorm:"column:email"`                          // 邮箱
	Phone      string `gorm:"column:phone"`                          // 手机号
	Avatar     string `gorm:"column:avatar"`                         // 头像
	Identity   string `gorm:"column:identity"`                       // 验证token
	Info       string `gorm:"column:info"`                           // 描述
	ClientIp   string `gorm:"column:client_ip"`                      // 客户端IP
	ClientPort string `gorm:"column:client_port"`                    // 客户端端口
	IsLogout   uint32 `gorm:"column:is_logout;default:0"`            // 是否退出登录0否 1是
}

func (m *User) TableName() string {
	return "test.user"
}

// 创建用户
func CreateUser(m *User) (*User, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUser", err)
	}
	return m, err
}

// 更新用户
func UpdateUser(m *User) (*User, error) {
	uid := m.Uid
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Updates(m).Error
	if err != nil {
		log.Print("CreateUser", err)
	}
	return m, err
}

// 删除用户
func DeleteUser(uid uint64) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Delete(m).Error
	if err != nil {
		log.Print("DeleteUser", err)
	}
	return m, err
}

// 查找用户
func FindUserByName(username string) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("username = ?", username).Find(m).Debug().Error
	if err != nil {
		log.Print("FindUserByName", err)
	}
	return m, err
}

// 查找用户
func FindUserByUid(uid uint64) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Find(m).Debug().Error
	if err != nil {
		log.Print("FindUserByUid", err)
	}
	return m, err
}

// 指定用户
func GetUserByUids(uids []uint64) ([]*User, error) {
	m := &User{}
	var data []*User
	err := utils.DB.Table(m.TableName()).Where("uid in ?", uids).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetUserByUids", err)
	}
	return data, err
}
