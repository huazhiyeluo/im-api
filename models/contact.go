package models

import (
	"demoapi/utils"
	"log"
)

// 联系人表
type Contact struct {
	Id       uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT"` // ID
	Uid      uint64 `gorm:"column:uid"`                           // UID
	TargetId uint64 `gorm:"column:target_id;default:0"`           // 目标ID
	Type     uint32 `gorm:"column:type;default:0"`                // 消息类型  1用户 2群
	Desc     string `gorm:"column:desc"`                          // 描述
}

func (m *Contact) TableName() string {
	return "test.contact"
}

// 获取联系人
func GetContactList(uid uint64, ttype uint32) ([]*Contact, error) {
	m := &Contact{}
	var data []*Contact
	err := utils.DB.Table(m.TableName()).Where("uid = ? and type = ? ", uid, ttype).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetContactList", err)
	}
	return data, err
}

// 获取联系人
func GetGroupContactList(targetId uint64, ttype uint32) ([]*Contact, error) {
	m := &Contact{}
	var data []*Contact
	err := utils.DB.Table(m.TableName()).Where("target_id = ? and type = ? ", targetId, ttype).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetContactList", err)
	}
	return data, err
}

// 创建关联
func CreateContact(m *Contact) (*Contact, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateContact", err)
	}
	return m, err
}
