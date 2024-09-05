package model

import (
	"log"
	"qqapi/internal/utils"

	"gorm.io/gorm"
)

// 创建用户
func CreateUser(m *User) (*User, error) {
	err := utils.DB.Table(m.TableName()).Create(m).Error
	if err != nil {
		log.Print("CreateUser", err)
	}
	return m, err
}

// 删除用户
func DeleteUser(uid uint64) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Delete(m).Error
	if err != nil {
		log.Print("DeleteUser", err)
	}
	return m, err
}

// 查找用户-username
func FindUserByName(username string) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("username = ?", username).Find(m).Error
	if err != nil {
		log.Print("FindUserByName", err)
	}
	return m, err
}

// 查找用户-uid
func FindUserByUid(uid uint64) (*User, error) {
	m := &User{}
	err := utils.DB.Table(m.TableName()).Where("uid = ?", uid).Find(m).Error
	if err != nil {
		log.Print("FindUserByUid", err)
		return m, err
	}
	return m, err
}

// 查找用户-指定用户
func FindUserByUids(uids []uint64) ([]*User, error) {
	m := &User{}
	var data []*User
	err := utils.DB.Table(m.TableName()).Where("uid in ?", uids).Find(&data).Error
	if err != nil {
		log.Print("GetUserByUids", err)
		return data, err
	}
	return data, err
}

// 查找用户-关键字
func FindUserByKeyword(pageSize uint32, pageNum uint32, keyword string) ([]*User, int64, error) {
	m := &User{}
	var data []*User
	var count int64
	db := utils.DB.Table(m.TableName())

	if keyword != "" {
		db.Where("uid like ? or nickname like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err1 := db.Count(&count).Error
	if err1 != nil {
		return data, count, err1
	}

	offset := int((pageNum - 1) * pageSize)
	size := int(pageSize)

	err := db.Limit(size).Offset(offset).Order("uid asc").Find(&data).Error
	if err != nil {
		log.Print("FindUserByKeyword", err)
		return data, count, err
	}
	return data, count, err
}

// ACT 用户
func ActUser(uid uint64, fields []*Fields) (*User, error) {
	updates := make(map[string]interface{})
	for _, v := range fields {
		if v.Otype == 0 {
			updates[v.Field] = gorm.Expr(v.Field+" + ?", v.Value)
		}
		if v.Otype == 1 {
			updates[v.Field] = gorm.Expr(v.Field+" - ?", v.Value)
		}
		if v.Otype == 2 {
			updates[v.Field] = v.Value
		}
	}
	m, err := FindUserByUid(uid)
	if err != nil {
		log.Print("FindUserByUid", err)
		return m, err
	}
	if m.Uid == 0 {
		updates["uid"] = uid
		err = utils.DB.Table(m.TableName()).Create(updates).Error
		if err != nil {
			log.Print("ActUser", err)
			return m, err
		}
	} else {
		err = utils.DB.Table(m.TableName()).Where("uid = ?", uid).Updates(updates).Error
		if err != nil {
			log.Print("ActUser", err)
			return m, err
		}
	}
	m, err = FindUserByUid(uid)
	return m, err
}
