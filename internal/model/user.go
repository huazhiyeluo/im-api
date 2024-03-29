package model

import (
	"imapi/internal/utils"
	"log"
)

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
func FindUserByUids(uids []uint64) ([]*User, error) {
	m := &User{}
	var data []*User
	err := utils.DB.Table(m.TableName()).Where("uid in ?", uids).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetUserByUids", err)
	}
	return data, err
}
