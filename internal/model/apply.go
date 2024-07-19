package model

import (
	"log"
	"qqapi/internal/utils"
)

// 创建申请表
func CreateApply(m *Apply) (*Apply, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateApply", err)
	}
	return m, err
}

// 更新申请表
func UpdateApply(m *Apply) (*Apply, error) {
	id := m.Id
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Updates(m).Error
	if err != nil {
		log.Print("UpdateApply", err)
	}
	return m, err
}

// 查找申请 - id
func FindApplyById(id uint32) (*Apply, error) {
	m := &Apply{}
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Find(m).Error
	if err != nil {
		log.Print("FindApplyById", err)
	}
	return m, err
}

func FindApplyByTwoId(fromId uint64, toId uint64, ttype uint32) (*Apply, error) {
	m := &Apply{}
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and to_id = ?  and type = ? and status = ?", fromId, toId, ttype, 0).Find(m).Error
	if err != nil {
		log.Print("FindApplyByTwoId", err)
	}
	return m, err
}

// 获取好友申请信息
func GetFriendApplyList(uid uint64) ([]*Apply, error) {
	m := &Apply{}
	var data []*Apply
	err := utils.DB.Table(m.TableName()).Where("((from_id = ? or to_id = ?) and type = 1)", uid, uid).Find(&data).Error
	if err != nil {
		log.Print("GetApplyList", err)
	}
	return data, err
}

// 获取群申请信息
func GetGroupApplyList(uid uint64, groupIds []uint64) ([]*Apply, error) {
	m := &Apply{}
	var data []*Apply
	err := utils.DB.Table(m.TableName()).Where("((to_id in ? and type = 2) or (from_id = ? and type = 2))", groupIds, uid).Find(&data).Error
	if err != nil {
		log.Print("GetApplyList", err)
	}
	return data, err
}
