package model

import (
	"imapi/internal/utils"
	"log"
)

// 创建群
func CreateGroup(m *Group) (*Group, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateGroup", err)
	}
	return m, err
}

// 更新群
func UpdateGroup(m *Group) (*Group, error) {
	groupId := m.GroupId
	err := utils.DB.Table(m.TableName()).Where("group_id = ?", groupId).Updates(m).Error
	if err != nil {
		log.Print("UpdateGroup", err)
	}
	return m, err
}

// 删除群
func DeleteGroup(groupId uint64) (*Group, error) {
	m := &Group{}
	err := utils.DB.Table(m.TableName()).Where("group_id = ?", groupId).Delete(m).Error
	if err != nil {
		log.Print("DeleteGroup", err)
	}
	return m, err
}

// 查找用群- name
func FindGroupByName(name string) (*Group, error) {
	m := &Group{}
	err := utils.DB.Table(m.TableName()).Where("name = ?", name).Find(m).Error
	if err != nil {
		log.Print("FindGroupByName", err)
	}
	return m, err
}

// 查找群 - groupId
func FindGroupByGroupId(groupId uint64) (*Group, error) {
	m := &Group{}
	err := utils.DB.Table(m.TableName()).Where("group_id = ?", groupId).Find(m).Error
	if err != nil {
		log.Print("FindGroupByGroupId", err)
	}
	return m, err
}

// 指定群 - groupIds
func FindGroupByGroupIds(groupIds []uint64) ([]*Group, error) {
	m := &Group{}
	var data []*Group
	err := utils.DB.Table(m.TableName()).Where("group_id in ?", groupIds).Find(&data).Error
	if err != nil {
		log.Print("GetGroupByGroupIds", err)
	}
	return data, err
}

// 指定群 - 拥有者
func GetGroupByOwnerUid(ownerUid uint64) ([]*Group, error) {
	m := &Group{}
	var data []*Group
	err := utils.DB.Table(m.TableName()).Where("owner_uid = ?", ownerUid).Find(&data).Error
	if err != nil {
		log.Print("GetGroupByOwnerUid", err)
	}
	return data, err
}
