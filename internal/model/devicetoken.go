package model

import (
	"log"
	"qqapi/internal/utils"

	"gorm.io/gorm"
)

// 查找设备TOKEN-uid
func FindDeviceTokenByUid(uid uint64) (*DeviceToken, error) {
	m := &DeviceToken{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Find(m).Error
	if err != nil {
		log.Print("FindDeviceTokenByUid", err)
		return m, err
	}
	return m, err
}

// ACT 设备TOKEN
func ActDeviceToken(uid uint64, fields []*Fields) (*DeviceToken, error) {
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
	m, err := FindDeviceTokenByUid(uid)
	if err != nil {
		log.Print("FindDeviceTokenByUid", err)
		return m, err
	}
	if m.Uid == 0 {
		updates["uid"] = uid
		err = utils.DB.Table(m.TableName()).Create(updates).Error
		if err != nil {
			log.Print("ActDeviceToken", err)
			return m, err
		}
	} else {
		err = utils.DB.Table(m.TableName()).Where("uid = ?", uid).Updates(updates).Error
		if err != nil {
			log.Print("ActDeviceToken", err)
			return m, err
		}
	}
	m, err = FindDeviceTokenByUid(uid)
	return m, err
}
