package login

import (
	"errors"
	"fmt"
	"qqapi/internal/model"
	"qqapi/internal/schema"
	"qqapi/internal/utils"
)

// VisitorLogin 游客登录适配器
type VisitorLogin struct {
	in *schema.LoginData
	pv *schema.PublicVar
}

func (b *VisitorLogin) Verify() {
	b.pv.Sid = Const_VISITOR
	usermapDevice, err := model.GetUsermapDeviceByDeviceid(b.in.Deviceid)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return
	}
	if usermapDevice.Deviceid == "" {
		m := &model.UsermapDevice{
			Deviceid: b.in.Deviceid,
			Siteuid:  utils.GenGUID(),
		}
		usermapDevice, err = model.CreateUsermapDevice(m)
		if err != nil {
			b.pv.Err = errors.New("DB Error")
			return
		}
	}
	if usermapDevice.Siteuid == "" {
		b.pv.Err = errors.New("DB Error")
		return
	}
	b.pv.Siteuid = usermapDevice.Siteuid
}

func (b *VisitorLogin) IsNewUser() {
	usermap, err := model.GetUserMapMix(b.pv.Siteuid, b.pv.Sid)
	if err != nil {
		b.pv.Err = errors.New("DB Error")
		return
	}
	if usermap.Uid == 0 {
		b.pv.IsNewUser = 1
		m := &model.Usermap{
			Siteuid: b.pv.Siteuid,
			Sid:     b.pv.Sid,
		}
		usermap, err = model.CreateUsermap(m)
		if err != nil {
			b.pv.Err = errors.New("DB Error")
			return
		}
	}
	b.pv.Uid = usermap.Uid
	b.pv.Nickname = fmt.Sprintf("USER_%d", usermap.Uid)
	b.pv.Avatar = fmt.Sprintf("http://img.siyuwen.com/godata/avatar/%d.jpg", utils.GetRandNum(0, 580))
}
