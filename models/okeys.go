package models

import "fmt"

// 设置游戏seesion key
func Rkonline(uid uint64) string {
	return fmt.Sprintf("RK_ND_SESSID%d", uid)
}
