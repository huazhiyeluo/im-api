package service

import (
	"demoapi/models"
	"demoapi/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditGroup(c *gin.Context) {
	data := schema.EditGroup{}
	c.Bind(&data)
	group, err := models.FindGroupByName(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}
	if group.GroupId != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已经存在"})
		return
	}
	insertData := &models.Group{
		OwnerUid: data.OwnerUid,
		Type:     data.Type,
		Name:     data.Name,
		Icon:     data.Icon,
		Info:     data.Info,
	}
	group, err = models.CreateGroup(insertData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	insertContactData := &models.Contact{
		Uid:      data.OwnerUid,
		TargetId: group.GroupId,
		Type:     2,
		Desc:     "",
	}
	_, err = models.CreateContact(insertContactData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "操作错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": group,
	})
}
