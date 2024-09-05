package model

import (
	"log"
	"qqapi/internal/utils"
)

// 获取好友分组
func GetFriendGroup(ownUid uint64) ([]*FriendGroup, error) {
	m := &FriendGroup{}
	var data []*FriendGroup
	err := utils.DB.Table(m.TableName()).Where("owner_uid = ? ", ownUid).Find(&data).Error
	if err != nil {
		log.Print("GetFriendGroup", err)
	}
	return data, err
}

// 获取好友分组 ByName
func GetFriendGroupByName(ownUid uint64, name string) (*FriendGroup, error) {
	m := &FriendGroup{}
	err := utils.DB.Table(m.TableName()).Where("owner_uid = ?  and name = ? ", ownUid, name).Find(&m).Error
	if err != nil {
		log.Print("GetFriendGroupByName", err)
	}
	return m, err
}

// 创建好友分组
func CreateFriendGroup(m *FriendGroup) (*FriendGroup, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateFriendGroup", err)
	}
	return m, err
}

// 更新好友分组
func UpdateFriendGroup(m *FriendGroup) (*FriendGroup, error) {
	friendGroupId := m.FriendGroupId
	err := utils.DB.Table(m.TableName()).Where("friendGroupId = ?", friendGroupId).Updates(m).Error
	if err != nil {
		log.Print("UpdateFriendGroup", err)
	}
	return m, err
}

// 删除好友分组
func DeleteFriendGroup(friendGroupId uint32) (*FriendGroup, error) {
	m := &FriendGroup{}
	err := utils.DB.Table(m.TableName()).Where("friend_group_id = ? ", friendGroupId).Delete(m).Error
	if err != nil {
		log.Print("DeleteFriendGroup", err)
	}
	return m, err
}
