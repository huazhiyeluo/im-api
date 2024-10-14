package login

import (
	"context"
	"errors"
	"fmt"
	"log"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/utils"
	"time"
)

const (
	Const_VISITOR  = 0
	Const_ACCOUT   = 1
	Const_GOOGLE   = 2
	Const_FACEBOOK = 3
)

// LoginAdapter 接口，定义登录操作
type LoginAdapter interface {
	Verify()
	IsNewUser()
}

func setToken(uid uint64) string {
	nowtime := time.Now().Unix()
	token := utils.GenMd5(fmt.Sprintf("%d%d", uid, nowtime))
	rkey := model.Rktoken(uid)

	utils.RDB.Set(context.TODO(), rkey, token, time.Minute*time.Duration(0))
	utils.RDB.ExpireAt(context.TODO(), rkey, time.Now().Add(time.Minute*60*24*2))
	return token
}

// GetLoginAdapter 根据 platform 获取对应的适配器
func GetLoginAdapter(cin *schema.CommonData, in *schema.LoginData, pv *schema.PublicVar) (LoginAdapter, error) {
	switch in.Platform {
	case "visitor":
		return &VisitorLogin{cin: cin, in: in, pv: pv}, nil
	case "account":
		return &AccountLogin{cin: cin, in: in, pv: pv}, nil
	case "facebook":
		return &FacebookLogin{cin: cin, in: in, pv: pv}, nil
	case "google":
		return &GoogleLogin{cin: cin, in: in, pv: pv}, nil
	default:
		return nil, fmt.Errorf("不支持的平台类型：%s", in.Platform)
	}
}

func Init(cin *schema.CommonData, in *schema.LoginData, pv *schema.PublicVar) {
	if cin.Deviceid == "" {
		pv.Err = errors.New("设备号不能为空")
	}
}

func Action(cin *schema.CommonData, in *schema.LoginData, pv *schema.PublicVar, adapter LoginAdapter) interface{} {
	adapter.IsNewUser()
	nowtime := time.Now().Unix()
	updates := make(map[string]interface{})
	if pv.IsNewUser == 1 {
		updates["nickname"] = pv.Nickname
		updates["avatar"] = pv.Avatar
		updates["devname"] = cin.Devname
		updates["deviceid"] = cin.Deviceid
		updates["reg_time"] = nowtime
		updates["login_time"] = nowtime
	} else {
		updates["devname"] = cin.Devname
		updates["deviceid"] = cin.Deviceid
		updates["login_time"] = nowtime
	}
	if in.Platform == "account" {
		if in.Nickname != "" {
			updates["nickname"] = in.Nickname
		}
		if in.Nickname != "" {
			updates["avatar"] = in.Avatar
		}
	}
	var updateData []*model.Fields
	for key, val := range updates {
		updateData = append(updateData, &model.Fields{Field: key, Otype: 2, Value: val})
	}
	user, err := model.ActUser(pv.Uid, updateData)
	if err != nil {
		pv.Err = errors.New("DB Error")
		return nil
	}
	res := make(map[string]interface{})
	tempUser := schema.GetResUser(user)
	usermapSso, err := model.GetUserMapSsoMix(pv.Uid)
	if err != nil {
		pv.Err = errors.New("DB Error")
		return nil
	}
	if usermapSso.Id != 0 {
		tempUser.Username = usermapSso.Username
		tempUser.Phone = usermapSso.Phone
		tempUser.Email = usermapSso.Email
	}

	if pv.IsNewUser == 1 {

		insertFriendGroupData := &model.FriendGroup{
			OwnerUid:  pv.Uid,
			Name:      "默认分组",
			IsDefault: 1,
		}
		friendGroup, err := model.CreateFriendGroup(insertFriendGroupData)
		if err != nil {
			log.Printf("%v", friendGroup)
			pv.Err = errors.New("DB Error")
			return nil
		}

		var updatesContactFriend []*model.Fields
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "friend_group_id", Otype: 2, Value: friendGroup.FriendGroupId})
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "level", Otype: 2, Value: 1})
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "remark", Otype: 2, Value: ""})
		updatesContactFriend = append(updatesContactFriend, &model.Fields{Field: "join_time", Otype: 2, Value: nowtime})
		toContactFriend, err := model.ActContactFriend(pv.Uid, pv.Uid, updatesContactFriend)
		if err != nil {
			log.Printf("%v", toContactFriend)
			pv.Err = errors.New("DB Error")
			return nil
		}
	}

	res["user"] = tempUser
	res["token"] = setToken(user.Uid)
	return res
}

func Login(cin *schema.CommonData, in *schema.LoginData) (interface{}, error) {
	pv := &schema.PublicVar{}
	adapter, err := GetLoginAdapter(cin, in, pv)
	if err != nil {
		return "", err
	}
	//1、初始化数据
	Init(cin, in, pv)
	if pv.Err != nil {
		return "", pv.Err
	}
	//2、验证数据
	adapter.Verify()
	if pv.Err != nil {
		return "", pv.Err
	}
	//3、登录数据操作
	res := Action(cin, in, pv, adapter)
	if pv.Err != nil {
		return "", pv.Err
	}
	return res, nil
}
