package service

import (
	"imapi/internal/model"
	"imapi/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditGroup(c *gin.Context) {
	data := schema.EditGroup{}
	c.Bind(&data)
	group, err := model.FindGroupByName(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经存在"})
		return
	}
	insertData := &model.Group{
		OwnerUid: data.OwnerUid,
		Type:     data.Type,
		Name:     data.Name,
		Icon:     data.Icon,
		Info:     data.Info,
	}
	group, err = model.CreateGroup(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertContactData := &model.Contact{
		FromId: data.OwnerUid,
		ToId:   group.GroupId,
		Type:   2,
		Remark: "",
	}
	_, err = model.CreateContact(insertContactData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}
