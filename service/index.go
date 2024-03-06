package service

import (
	"demoapi/models"

	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	models.Chat(c)
}
