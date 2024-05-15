package model

import (
	"imapi/internal/utils"
	"log"
)

// 获取好友列表
func GetContactUserList(fromId uint64) ([]*ContactUser, error) {
	m := &ContactUser{}
	var data []*ContactUser
	err := utils.DB.Table(m.TableName()).Where("from_id = ?", fromId).Find(&data).Error
	if err != nil {
		log.Print("GetContactUserList", err)
	}
	return data, err
}

// 获取单个好友
func GetContactUserOne(fromId uint64, toId uint64) (*ContactUser, error) {
	m := &ContactUser{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Find(&m).Debug().Error
	if err != nil {
		log.Print("GetContactUserOne", err)
	}
	return m, err
}

// 创建好友关联
func CreateContactUser(m *ContactUser) (*ContactUser, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateContactUser", err)
	}
	return m, err
}

// 更新好友关联
func UpdateContactUser(m *ContactUser) (*ContactUser, error) {
	fromId := m.FromId
	toId := m.ToId
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Updates(m).Error
	if err != nil {
		log.Print("UpdateContactUser", err)
	}
	return m, err
}

// 删除关联
func DeleteContactUser(fromId uint64, toId uint64) (*ContactUser, error) {
	m := &ContactUser{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Delete(m).Error
	if err != nil {
		log.Print("DeleteContactUser", err)
	}
	return m, err
}
