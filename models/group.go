package models

import (
	"demoapi/utils"
	"log"
)

// 群管理表
type Group struct {
	GroupId  uint64 `gorm:"column:group_id;primary_key;AUTO_INCREMENT"` // ID
	OwnerUid uint64 `gorm:"column:owner_uid;default:0"`                 // 拥有者
	Type     uint32 `gorm:"column:type;default:0"`                      // 群类型
	Name     string `gorm:"column:name"`                                // 名称
	Icon     string `gorm:"column:icon"`                                // 图标
	Info     string `gorm:"column:info"`                                // 描述
}

func (m *Group) TableName() string {
	return "test.group"
}

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
		log.Print("CreateGroup", err)
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

// 指定群
func GetGroupByGroupIds(groupIds []uint64) ([]*Group, error) {
	m := &Group{}
	var data []*Group
	err := utils.DB.Table(m.TableName()).Where("group_id in ?", groupIds).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetGroupByGroupIds", err)
	}
	return data, err
}

// 查找用户
func FindGroupByName(name string) (*Group, error) {
	m := &Group{}
	err := utils.DB.Table(m.TableName()).Where("name = ?", name).Find(m).Debug().Error
	if err != nil {
		log.Print("FindGroupByName", err)
	}
	return m, err
}
