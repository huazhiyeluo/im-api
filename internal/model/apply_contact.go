package model

import (
	"imapi/internal/utils"
	"log"
)

// 创建申请表
func CreateApplyContact(m *ApplyContact) (*ApplyContact, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateApplyContact", err)
	}
	return m, err
}

// 更新申请表
func UpdateApplyContact(m *ApplyContact) (*ApplyContact, error) {
	id := m.Id
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Updates(m).Error
	if err != nil {
		log.Print("CreateApplyContact", err)
	}
	return m, err
}

func FindApplyById(id uint32) (*ApplyContact, error) {
	m := &ApplyContact{}
	err := utils.DB.Table(m.TableName()).Where("id = ?", id).Find(m).Debug().Error
	if err != nil {
		log.Print("FindApplyById", err)
	}
	return m, err
}
