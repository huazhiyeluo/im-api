package model

import (
	"imapi/internal/utils"
	"log"
)

// 获取群公告列表
func GetGroupTipsList(groupId uint64) ([]*GroupTips, error) {
	m := &GroupTips{}
	var data []*GroupTips
	err := utils.DB.Table(m.TableName()).Where("group_id = ?", groupId).Find(&data).Error
	if err != nil {
		log.Print("GetGroupTipsList", err)
	}
	return data, err
}

// 创建群公告
func CreateGroupTips(m *GroupTips) (*GroupTips, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateGroupTips", err)
	}
	return m, err
}

// 更新群公告
func UpdateGroupTips(m *GroupTips) (*GroupTips, error) {
	id := m.Id
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Updates(m).Error
	if err != nil {
		log.Print("UpdateGroupTips", err)
	}
	return m, err
}

// 删除群公告
func DeleteGroupTips(id uint32) (*GroupTips, error) {
	m := &GroupTips{}
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Delete(m).Error
	if err != nil {
		log.Print("DeleteGroupTips", err)
	}
	return m, err
}

// 删除群公告-所有
func DeleteGroupTipsAll(groupId uint64) (*GroupTips, error) {
	m := &GroupTips{}
	err := utils.DB.Table(m.TableName()).Where("group_id = ?", groupId).Delete(m).Error
	if err != nil {
		log.Print("DeleteGroupTipsAll", err)
	}
	return m, err
}
