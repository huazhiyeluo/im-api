package model

import (
	"imapi/internal/utils"
	"log"
)

// 获取好友|组
func GetContactGroup(ownUid uint32) ([]*ContactGroup, error) {
	m := &ContactGroup{}
	var data []*ContactGroup
	err := utils.DB.Table(m.TableName()).Where("owner_uid = ? ", ownUid).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetContactGroup", err)
	}
	return data, err
}
