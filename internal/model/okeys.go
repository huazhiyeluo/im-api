package model

import "fmt"

// 设置游戏token
func Rktoken(uid uint64) string {
	return fmt.Sprintf("RK_ND_TOKEN%d", uid)
}

// 设置游戏seesionKey
func RksessionKey(uid uint64) string {
	return fmt.Sprintf("RK_ND_SESSION_KEY%d", uid)
}

// 设置用户消息通知
func Rkmsg(fromId uint64, toId uint64) string {
	return fmt.Sprintf("msg_%d_%d", fromId, toId)
}

// 设置用户消息通知
func RkUreadMsg(toId uint64) string {
	return fmt.Sprintf("unread_msg_%d", toId)
}
