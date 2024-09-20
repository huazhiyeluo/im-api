package model

import (
	"log"
	"qqapi/internal/utils"

	"gorm.io/gorm"
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

// 查找群用户-指定用户
func FindContactGroupByFromIds(fromIds []uint64, toId uint64) ([]*ContactGroup, error) {
	m := &ContactGroup{}
	var data []*ContactGroup
	err := utils.DB.Table(m.TableName()).Where("from_id in ? and to_id = ?", fromIds, toId).Find(&data).Error
	if err != nil {
		log.Print("GetUserByUids", err)
		return data, err
	}
	return data, err
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

// 删除群关联-指定用户
func DeleteContactGroupByFromIds(fromIds []uint64, toId uint64) (*ContactGroup, error) {
	m := &ContactGroup{}
	err := utils.DB.Table(m.TableName()).Where("from_id in ? and to_id = ?", fromIds, toId).Delete(m).Error
	if err != nil {
		log.Print("DeleteContactGroupByFromIds", err)
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

// ACT 联系人群组
func ActContactGroup(fromId uint64, toId uint64, fields []*Fields) (*ContactGroup, error) {
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
	m, err := GetContactGroupOne(fromId, toId)
	if err != nil {
		log.Print("GetContactGroupOne", err)
		return m, err
	}
	if m.FromId == 0 {
		updates["from_id"] = fromId
		updates["to_id"] = toId
		err = utils.DB.Table(m.TableName()).Create(updates).Error
		if err != nil {
			log.Print("ActContactGroup", err)
			return m, err
		}
	} else {
		err = utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?", fromId, toId).Updates(updates).Error
		if err != nil {
			log.Print("ActContactGroup", err)
			return m, err
		}
	}
	m, err = GetContactGroupOne(fromId, toId)
	return m, err
}
