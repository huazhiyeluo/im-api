package server

import (
	"fmt"
	"imapi/internal/model"
	"imapi/third_party/log"
)

// 1、用户状态
func UserStatusNoticeMsg(uid uint64, msgMedia uint32) {
	contacts, err := model.GetContactList(uid, 1)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateMsg(&Message{FromId: v.FromId, ToId: v.ToId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: "用户状态"}})
	}
}

// 2、用户好友
func UserFriendNoticeMsg(fromId uint64, toId uint64, content string, msgMedia uint32) {
	CreateMsg(&Message{FromId: fromId, ToId: toId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: content}})
}

// 3、用户群
func UserGroupNoticeMsg(groupId uint64, content string, msgMedia uint32) {
	contacts, err := model.GetGroupContactList(groupId, 2)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateMsg(&Message{FromId: v.FromId, ToId: v.FromId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: content}})
	}
}
