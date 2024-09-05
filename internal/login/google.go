package login

import (
	"errors"
	"qqapi/internal/model"
	"qqapi/internal/schema"

	"github.com/davecgh/go-spew/spew"
)

// GoogleLogin Google 登录适配器
type GoogleLogin struct {
	in *schema.LoginData
	pv *schema.PublicVar
}

func (b *GoogleLogin) Verify() {
	b.pv.Sid = Const_GOOGLE
	b.pv.Siteuid = b.in.Siteuid

	// usermap, err := model.GetUserMapMix(b.in.Siteuid, b.pv.Sid)
	// if err != nil {
	// 	b.pv.Err = errors.New("DB Error")
	// 	return
	// }
	// if usermap.Uid != 0 {
	// 	b.pv.Uid = usermap.Uid
	// 	return
	// }

	//	服务器在中国，连不上外网
	// url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", b.in.Token)
	// res, err := utils.HttpGet(url)
	// if err != nil {
	// 	b.pv.Err = errors.New("DB Error")
	// 	return
	// }

	// mapData := make(map[string]interface{})
	// err = json.Unmarshal([]byte(res), &mapData)
	// if err != nil {
	// 	b.pv.Err = errors.New("DB Error")
	// 	return
	// }
	// userid := utils.ToString(mapData["sub"])
	// if userid != b.in.Siteuid {
	// 	b.pv.Err = errors.New("DB Error")
	// 	return
	// }
	// b.pv.Nickname = utils.ToString(mapData["name"])
	// b.pv.Avatar = utils.ToString(mapData["picture"])

}

func (b *GoogleLogin) IsNewUser() {
	usermap, err := model.GetUserMapMix(b.pv.Siteuid, b.pv.Sid)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return
	}
	if usermap.Uid != 0 {
		b.pv.Uid = usermap.Uid
		return
	}
	usermap = b.BindVistor()
	if usermap.Uid != 0 {
		b.pv.Uid = usermap.Uid
		return
	}
	m := &model.Usermap{
		Siteuid: b.pv.Siteuid,
		Sid:     b.pv.Sid,
	}
	usermap, err = model.CreateUsermap(m)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return
	}
	b.pv.IsNewUser = 1
	b.pv.Uid = usermap.Uid
	b.pv.Nickname = b.in.Nickname
	b.pv.Avatar = b.in.Avatar
}

func (b *GoogleLogin) BindVistor() *model.Usermap {
	reply := &model.Usermap{}

	usermapDevice, err := model.GetUsermapDeviceByDeviceid(b.in.Deviceid)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return reply
	}
	if usermapDevice.Deviceid == "" {
		return reply
	}
	usermap, err := model.GetUserMapMix(usermapDevice.Siteuid, 0)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return reply
	}
	if usermap.Uid == 0 {
		return reply
	}

	spew.Dump(usermap.Uid, b.pv.Siteuid, b.pv.Sid)

	usermapBind, err := model.CreateUsermapBind(&model.UsermapBind{Uid: usermap.Uid, Siteuid: b.pv.Siteuid, Sid: b.pv.Sid})
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return reply
	}
	reply = &model.Usermap{Uid: usermapBind.Uid, Siteuid: usermapBind.Siteuid, Sid: usermapBind.Sid}
	return reply

}
