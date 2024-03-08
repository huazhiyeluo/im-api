package model

import (
	"demoapi/utils"
	"log"
)

// 获取联系人
func GetContactList(fromId uint64, ttype uint32) ([]*Contact, error) {
	m := &Contact{}
	var data []*Contact
	err := utils.DB.Table(m.TableName()).Where("from_id = ? and type = ? ", fromId, ttype).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetContactList", err)
	}
	return data, err
}

// 获取联系人
func GetGroupContactList(toId uint64, ttype uint32) ([]*Contact, error) {
	m := &Contact{}
	var data []*Contact
	err := utils.DB.Table(m.TableName()).Where("to_id = ? and type = ? ", toId, ttype).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetGroupContactList", err)
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
