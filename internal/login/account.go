package login

import (
	"errors"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/utils"
)

// AccountLogin 账号登录适配器
type AccountLogin struct {
	cin *schema.CommonData
	in  *schema.LoginData
	pv  *schema.PublicVar
}

func (b *AccountLogin) Verify() {
	b.pv.Sid = Const_ACCOUT
	usermapSso := &model.UsermapSso{}
	var err error
	if b.in.Username != "" {
		usermapSso, err = model.GetUsermapSsoByUsername(b.in.Username)
	} else if b.in.Phone != "" {
		usermapSso, err = model.GetUsermapSsoByPhone(b.in.Phone)
	} else if b.in.Email != "" {
		usermapSso, err = model.GetUsermapSsoByEmail(b.in.Email)
	}
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return
	}
	if usermapSso.Id == 0 {
		b.pv.Err = errors.New("账号不存在")
		return
	}
	if utils.GenMd5(b.in.Password) != usermapSso.Password {
		b.pv.Err = errors.New("密码错误")
		return
	}
	b.pv.Siteuid = usermapSso.Siteuid
}

func (b *AccountLogin) IsNewUser() {
	b.pv.Nickname = b.in.Nickname
	b.pv.Avatar = b.in.Avatar
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
}

func (b *AccountLogin) BindVistor() *model.Usermap {
	reply := &model.Usermap{}

	usermapDevice, err := model.GetUsermapDeviceByDeviceid(b.cin.Deviceid)
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

	bindList, err := model.GetUserMapBindList(usermap.Uid, b.pv.Sid)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return reply
	}
	if len(bindList) > 0 {
		return reply
	}

	usermapBind, err := model.CreateUsermapBind(&model.UsermapBind{Uid: usermap.Uid, Siteuid: b.pv.Siteuid, Sid: b.pv.Sid})
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return reply
	}
	reply = &model.Usermap{Uid: usermapBind.Uid, Siteuid: usermapBind.Siteuid, Sid: usermapBind.Sid}
	return reply

}
