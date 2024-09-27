package model

import (
	"context"
	"log"
	"qqapi/internal/utils"
	"strings"
)

// - UsermapDevice ----------------------------------------------------------------

func GetUsermapDeviceByDeviceid(deviceid string) (*UsermapDevice, error) {
	m := &UsermapDevice{}
	err := utils.DB.Table(m.TableName()).Where("deviceid = ?", deviceid).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapDevice", err)
	}
	return m, err
}

func GetUsermapDeviceBySiteuid(siteuid string) (*UsermapDevice, error) {
	m := &UsermapDevice{}
	err := utils.DB.Table(m.TableName()).Where("siteuid = ?", siteuid).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapDeviceBySiteuid", err)
	}
	return m, err
}

func CreateUsermapDevice(m *UsermapDevice) (*UsermapDevice, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUsermapDevice", err)
	}
	return m, err
}

func UpdateUsermapDevice(m *UsermapDevice) (*UsermapDevice, error) {
	siteuid := m.Siteuid
	err := utils.DB.Table(m.TableName()).Where("siteuid = ?", siteuid).Updates(m).Error
	if err != nil {
		log.Print("UpdateUsermapDevice", err)
	}
	return m, err
}

// - UsermapSso ----------------------------------------------------------------

func GetUsermapSsoBySiteuid(siteuid string) (*UsermapSso, error) {
	m := &UsermapSso{}
	err := utils.DB.Table(m.TableName()).Where("siteuid = ?", siteuid).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapSsoBySiteuid", err)
	}
	return m, err
}

func GetUsermapSsoByUsername(username string) (*UsermapSso, error) {
	m := &UsermapSso{}
	err := utils.DB.Table(m.TableName()).Where("username = ?", username).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapSsoByUsername", err)
	}
	return m, err
}

func GetUsermapSsoByPhone(phone string) (*UsermapSso, error) {
	m := &UsermapSso{}
	err := utils.DB.Table(m.TableName()).Where("phone = ?", phone).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapSsoByPhone", err)
	}
	return m, err
}

func GetUsermapSsoByEmail(email string) (*UsermapSso, error) {
	m := &UsermapSso{}
	err := utils.DB.Table(m.TableName()).Where("email = ?", email).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapSsoByEmail", err)
	}
	return m, err
}

func CreateUsermapSso(m *UsermapSso) (*UsermapSso, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUsermapSso", err)
	}
	return m, err
}

func UpdateUsermapSso(m *UsermapSso) (*UsermapSso, error) {
	siteuid := m.Siteuid
	err := utils.DB.Table(m.TableName()).Where("siteuid = ?", siteuid).Updates(m).Error
	if err != nil {
		log.Print("UpdateUsermapSso", err)
	}
	return m, err
}

// - Usermap ----------------------------------------------------------------

func GetUserMapMix(siteuid string, sid uint32) (*Usermap, error) {
	m := &Usermap{}
	err := utils.DB.Table(m.TableName()).Where("siteuid = ? and sid = ?", siteuid, sid).Limit(1).Find(&m).Error

	if m.Siteuid == "" {
		mbind := &UsermapBind{}
		err := utils.DB.Table(mbind.TableName()).Where("siteuid = ? and sid = ?", siteuid, sid).Limit(1).Find(&mbind).Error
		if err != nil {
			log.Print("GetUserMapMix", err)
			return nil, err
		}
		m = &Usermap{Uid: mbind.Uid, Siteuid: mbind.Siteuid, Sid: mbind.Sid}
	}
	return m, err
}

func GetUsermapMax(ctx context.Context) (*Usermap, error) {
	m := &Usermap{}
	err := utils.DB.Table(m.TableName()).Order("uid desc").Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetMaxUsermap", err)
	}
	return m, err
}

func GetUserMap(siteuid string, sid uint32) (*Usermap, error) {
	m := &Usermap{}
	err := utils.DB.Table(m.TableName()).Where("siteuid = ? and sid = ?", siteuid, sid).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUserMap", err)
	}
	return m, err
}

func GetUserMapList(uid uint64) ([]*Usermap, error) {
	m := &Usermap{}
	var data []*Usermap
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Find(&data).Error
	if err != nil {
		log.Print("GetUserMapList", err)
		return data, err
	}
	return data, err
}

func CreateUsermap(m *Usermap) (*Usermap, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUsermap", err)
	}
	return m, err
}

func DeleteUsermap(uid uint64, sid uint32) (*Usermap, error) {
	m := &Usermap{}
	err := utils.DB.Table(m.TableName()).Where("uid = ? and sid = ?", uid, sid).Delete(m).Error
	if err != nil {
		log.Print("DeleteUsermap", err)
	}
	return m, err
}

// - UsermapBind ----------------------------------------------------------------

func GetUserMapBind(siteuid string, sid uint32) (*UsermapBind, error) {
	m := &UsermapBind{}
	err := utils.DB.Table(m.TableName()).Where("siteuid = ? and sid = ?", siteuid, sid).Limit(1).Find(&m).Error
	if err != nil {
		log.Print("GetUsermapBind", err)
	}
	return m, err
}

func GetUserMapBindList(uid uint64, sid uint32) ([]*UsermapBind, error) {
	m := &UsermapBind{}
	var data []*UsermapBind
	err := utils.DB.Table(m.TableName()).Where("uid = ? and sid = ?", uid, sid).Find(&data).Error
	if err != nil {
		log.Print("GetUserMapList", err)
		return data, err
	}
	return data, err
}

func CreateUsermapBind(m *UsermapBind) (*UsermapBind, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUsermapBind", err)
	}
	return m, err
}

func DeleteUsermapBind(uid uint64, sid uint32) (*UsermapBind, error) {
	m := &UsermapBind{}
	err := utils.DB.Table(m.TableName()).Where("uid = ? and sid = ?", uid, sid).Delete(m).Error
	if err != nil {
		log.Print("DeleteUsermapBind", err)
	}
	return m, err
}

// - Usermap - UsermapBind - sso ----------------------------------------------------------------

func GetUserMapSsoMix(uid uint64) (*UsermapSso, error) {
	ms := &UsermapSso{}

	m := &Usermap{}
	err := utils.DB.Table(m.TableName()).Where("uid = ? and sid = ?", uid, 1).Find(&m).Error
	if err != nil {
		log.Print("GetUserMapMixList", err)
		return ms, err
	}
	if m.Uid == 0 {
		mb := &UsermapBind{}
		err = utils.DB.Table(mb.TableName()).Where("uid = ? and sid = ?", uid, 1).Find(&mb).Error
		if err != nil {
			log.Print("GetUserMapMixList", err)
			return ms, err
		}
		m = &Usermap{Uid: mb.Uid, Siteuid: mb.Siteuid, Sid: mb.Sid}
	}

	if m.Uid != 0 {
		err := utils.DB.Table(ms.TableName()).Where("siteuid = ?", m.Siteuid).Find(&ms).Error
		if err != nil {
			log.Print("GetUserMapMixList", err)
			return ms, err
		}
		if strings.Contains(ms.Username, "u_") {
			ms.Username = ""
		}
		if strings.Contains(ms.Phone, "p_") {
			ms.Phone = ""
		}
		if strings.Contains(ms.Email, "e_") {
			ms.Email = ""
		}
	}

	return ms, err
}
