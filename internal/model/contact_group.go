package model

import (
	"imapi/internal/utils"
	"log"
)

// 获取群列表
func GetContactGroupList(fromId uint64) ([]*ContactGroup, error) {
	m := &ContactGroup{}
	var data []*ContactGroup
	err := utils.DB.Table(m.TableName()).Where("from_id = ?", fromId).Find(&data).Error
	if err != nil {
		log.Print("GetContactGroupList", err)
	}
	return data, err
}

// 获取组成员-所有
func GetGroupUser(toId uint64) ([]*ContactGroup, error) {
	m := &ContactGroup{}
	var data []*ContactGroup
	err := utils.DB.Table(m.TableName()).Where("to_id = ?", toId).Find(&data).Error
	if err != nil {
		log.Print("GetGroupUser", err)
	}
	return data, err
}

// 获取单个组成员
func GetContactGroupOne(fromId uint64, toId uint64) (*ContactGroup, error) {
	m := &ContactGroup{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Find(&m).Error
	if err != nil {
		log.Print("GetGroupContactOne", err)
	}
	return m, err
}

// 创建群关联
func CreateContactGroup(m *ContactGroup) (*ContactGroup, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateContactGroup", err)
	}
	return m, err
}

// 更新群关联
func UpdateContactGroup(m *ContactGroup) (*ContactGroup, error) {
	fromId := m.FromId
	toId := m.ToId
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Updates(m).Error
	if err != nil {
		log.Print("UpdateContactGroup", err)
	}
	return m, err
}

// 删除群关联
func DeleteContactGroup(fromId uint64, toId uint64) (*ContactGroup, error) {
	m := &ContactGroup{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Delete(m).Error
	if err != nil {
		log.Print("DeleteContactGroup", err)
	}
	return m, err
}

// 删除群关联-所有人
func DeleteContactGroupAll(toId uint64) (*ContactGroup, error) {
	m := &ContactGroup{}
	err := utils.DB.Table(m.TableName()).Where("to_id = ?", toId).Delete(m).Error
	if err != nil {
		log.Print("DeleteContactGroupAll", err)
	}
	return m, err
}
