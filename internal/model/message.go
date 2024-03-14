package model

import (
	"imapi/internal/utils"
	"log"
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
