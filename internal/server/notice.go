package server

import (
	"fmt"
	"qqapi/internal/model"
	"qqapi/third_party/log"
)

// 1-1、用户状态
func UserStatusNoticeMsg(uid uint64, msgMedia uint32) {
	contacts, err := model.GetContactFriendList(uid)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		if v.ToId != uid {
			CreateMsg(&Message{FromId: v.FromId, ToId: v.ToId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: "用户状态"}})
		}
	}
}

// 1-2、用户信息
func UserInfoNoticeMsg(uid uint64, content string) {
	contacts, err := model.GetContactFriendList(uid)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateMsg(&Message{FromId: v.FromId, ToId: v.ToId, MsgType: MSG_TYPE_NOTICE, MsgMedia: MSG_MEDIA_USERINFO, Content: &MessageContent{Data: content}})
	}
}

// 1-3、群信息
func GroupInfoNoticeMsg(groupId uint64, content string) {
	contacts, err := model.GetGroupUser(groupId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateMsg(&Message{FromId: v.FromId, ToId: v.FromId, MsgType: MSG_TYPE_NOTICE, MsgMedia: MSG_MEDIA_GROUPINFO, Content: &MessageContent{Data: content}})
	}
}

// 2-1、用户好友-消息
func UserFriendMsg(fromId uint64, toId uint64, content string, msgMedia uint32) {
	CreateMsg(&Message{FromId: fromId, ToId: toId, MsgType: MSG_TYPE_SINGLE, MsgMedia: msgMedia, Content: &MessageContent{Data: content}})
}

// 2-2、用户好友-通知
func UserFriendNoticeMsg(fromId uint64, toId uint64, content string, msgMedia uint32) {
	CreateMsg(&Message{FromId: fromId, ToId: toId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: content}})
}

// 3、用户群
func UserGroupNoticeMsg(groupId uint64, content string, msgMedia uint32) {
	contacts, err := model.GetGroupUser(groupId)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateMsg(&Message{FromId: v.FromId, ToId: v.FromId, MsgType: MSG_TYPE_NOTICE, MsgMedia: msgMedia, Content: &MessageContent{Data: content}})
	}
}
