package service

import (
	"demoapi/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Upload(c *gin.Context) {

	file, _ := c.FormFile("file")

	ext := filepath.Ext(file.Filename)

	nowtime := time.Now().Unix()
	tempPath := utils.GenMd5(fmt.Sprintf("%s-%d", file.Filename, nowtime))

	dst := fmt.Sprintf("./static/images/%s%s", tempPath, ext)
	url := fmt.Sprintf("%s/images/%s%s", viper.GetString("cdn.url"), tempPath, ext)

	c.SaveUploadedFile(file, dst)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": url,
	})
}
