package model

import (
	"log"
	"qqapi/internal/utils"
)

// 创建消息
func CreateMessage(m *Message) (*Message, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateMessage", err)
	}
	return m, err
}

// 获取消息
func GetMessageList(pageSize uint32, pageNum uint32, toId uint64, msgType uint32) ([]*Message, int64, error) {
	m := &Message{}
	var data []*Message
	var count int64

	db := utils.DB.Table(m.TableName())

	db.Where("to_id = ?", toId)
	db.Where("msg_type = ?", msgType)

	err1 := db.Count(&count).Error
	if err1 != nil {
		return data, count, err1
	}

	offset := int((pageNum - 1) * pageSize)
	size := int(pageSize)
	err := db.Limit(size).Offset(offset).Order("create_time desc").Find(&data).Error
	if err != nil {
		log.Print("GetMessageList", err)
	}
	return data, count, err
}

func GetMessageAll(ids []string) ([]*Message, error) {
	m := &Message{}
	var data []*Message

	db := utils.DB.Table(m.TableName())
	db.Where("id in ?", ids)

	err := db.Order("create_time asc").Find(&data).Error
	if err != nil {
		log.Print("GetMessageList", err)
	}
	return data, err
}

//-------------------------------------------------------------------MessageUnread--------------------------------------------------------------------

// 创建未读消息
func CreateMessageUnread(m *MessageUnread) (*MessageUnread, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateMessageUnread", err)
	}
	return m, err
}

// 获取未读消息
func GetMessageUnread(uid uint64) ([]*MessageUnread, error) {
	m := &MessageUnread{}
	var data []*MessageUnread
	db := utils.DB.Table(m.TableName())

	db.Where("uid = ?", uid)

	err := db.Limit(500).Order("create_time asc").Find(&data).Error
	if err != nil {
		log.Print("GetMessageList", err)
	}
	return data, err
}

// 获取消息
func GetMessageUnreadCount(uid uint64) (int64, error) {
	m := &MessageUnread{}
	var count int64

	db := utils.DB.Table(m.TableName())

	db.Where("uid = ?", uid)

	err := db.Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, err
}

func DeleteMessageUnread(uid uint64) (*MessageUnread, error) {
	m := &MessageUnread{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Delete(m).Error
	if err != nil {
		log.Print("DeleteMessageUnread", err)
	}
	return m, err
}
