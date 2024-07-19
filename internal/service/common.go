package service

import (
	"context"
	"fmt"
	"qqapi/internal/model"
	"qqapi/internal/utils"
	"time"
)

func setToken(uid uint64) string {
	nowtime := time.Now().Unix()
	token := utils.GenMd5(fmt.Sprintf("%d%d", uid, nowtime))
	rkey := model.Rktoken(uid)

	utils.RDB.Set(context.TODO(), rkey, token, time.Minute*time.Duration(0))
	utils.RDB.ExpireAt(context.TODO(), rkey, time.Now().Add(time.Minute*60*24*2))
	return token
}
