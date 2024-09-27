package model

import (
	"log"
	"qqapi/internal/utils"

	"gorm.io/gorm"
)

// 获取好友列表
func GetContactFriendList(fromId uint64) ([]*ContactFriend, error) {
	m := &ContactFriend{}
	var data []*ContactFriend
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and join_time > ?", fromId, 0).Find(&data).Error
	if err != nil {
		log.Print("GetContactFriendList", err)
	}
	return data, err
}

// 获取单个好友
func GetContactFriendOne(fromId uint64, toId uint64) (*ContactFriend, error) {
	m := &ContactFriend{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ? and join_time > ?", fromId, toId, 0).Find(&m).Error
	if err != nil {
		log.Print("GetContactFriendOne", err)
	}
	return m, err
}

// 查找好友-指定好友
func GetContactFriendByToIds(fromId uint64, toIds []uint64) ([]*ContactFriend, error) {
	m := &ContactFriend{}
	var data []*ContactFriend
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id in ? and join_time > ?", fromId, toIds, 0).Find(&data).Error
	if err != nil {
		log.Print("GetUserByUids", err)
		return data, err
	}
	return data, err
}

// 创建好友关联
func CreateContactFriend(m *ContactFriend) (*ContactFriend, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateContactFriend", err)
	}
	return m, err
}

// 更新好友关联
func UpdateContactFriend(m *ContactFriend) (*ContactFriend, error) {
	fromId := m.FromId
	toId := m.ToId
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Updates(m).Error
	if err != nil {
		log.Print("UpdateContactFriend", err)
	}
	return m, err
}

// 更新好友关联 批量更新联系人组
func UpdateContactFriendByFriendGroupId(friendGroupId uint32, m *ContactFriend) (*ContactFriend, error) {
	fromId := m.FromId
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and friend_group_id = ?", fromId, friendGroupId).Updates(m).Debug().Error
	if err != nil {
		log.Print("UpdateContactFriendByFriendGroupId", err)
	}
	return m, err
}

// 删除关联
func DeleteContactFriend(fromId uint64, toId uint64) (*ContactFriend, error) {
	m := &ContactFriend{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Delete(m).Error
	if err != nil {
		log.Print("DeleteContactFriend", err)
	}
	return m, err
}

// ACT 联系人
func ActContactFriend(fromId uint64, toId uint64, fields []*Fields) (*ContactFriend, error) {
	updates := make(map[string]interface{})
	for _, v := range fields {
		if v.Otype == 0 {
			updates[v.Field] = gorm.Expr(v.Field+" + ?", v.Value)
		}
		if v.Otype == 1 {
			updates[v.Field] = gorm.Expr(v.Field+" - ?", v.Value)
		}
		if v.Otype == 2 {
			updates[v.Field] = v.Value
		}
	}
	m, err := GetContactFriendOne(fromId, toId)
	if err != nil {
		log.Print("GetContactFriendOne", err)
		return m, err
	}
	if m.FromId == 0 {
		updates["from_id"] = fromId
		updates["to_id"] = toId
		err = utils.DB.Table(m.TableName()).Create(updates).Error
		if err != nil {
			log.Print("ActContactFriend", err)
			return m, err
		}
	} else {
		err = utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Updates(updates).Error
		if err != nil {
			log.Print("ActContactFriend", err)
			return m, err
		}
	}
	m, err = GetContactFriendOne(fromId, toId)
	return m, err
}
