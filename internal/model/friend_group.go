package model

import (
	"qqapi/internal/utils"
	"log"
)

// 获取好友分组
func GetFriendGroup(ownUid uint32) ([]*FriendGroup, error) {
	m := &FriendGroup{}
	var data []*FriendGroup
	err := utils.DB.Table(m.TableName()).Where("owner_uid = ? ", ownUid).Find(&data).Debug().Error
	if err != nil {
		log.Print("GetFriendGroup", err)
	}
	return data, err
}
