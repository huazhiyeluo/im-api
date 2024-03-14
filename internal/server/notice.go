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
		CreateNoticeMsg(&Message{FromId: v.FromId, ToId: v.ToId, MsgType: 4, MsgMedia: msgMedia, Content: "用户状态"})
	}
}

// 2、用户好友
func UserFriendNoticeMsg(fromId uint64, toId uint64, msgMedia uint32) {
	CreateNoticeMsg(&Message{FromId: fromId, ToId: toId, MsgType: 4, MsgMedia: msgMedia, Content: "用户好友"})
}

// 3、用户群
func UserGroupNoticeMsg(groupId uint64, msgMedia uint32) {
	contacts, err := model.GetGroupContactList(groupId, 2)
	if err != nil {
		log.Logger.Info(fmt.Sprintf("%v ", err))
	}
	for _, v := range contacts {
		CreateNoticeMsg(&Message{FromId: v.FromId, ToId: v.FromId, MsgType: 4, MsgMedia: msgMedia, Content: "用户群"})
	}
}
