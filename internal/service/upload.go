package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"qqapi/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Upload(c *gin.Context) {

	file, _ := c.FormFile("file")

	ext := filepath.Ext(file.Filename)

	nowtime := time.Now().Unix()
	tempPath := utils.GenMd5(fmt.Sprintf("%s-%d", file.Filename, nowtime))

	tdate := utils.GetDate(nowtime)

	dst := fmt.Sprintf("%s/%s/%s%s", viper.GetString("cdn.path"), tdate, tempPath, ext)
	url := fmt.Sprintf("%s/%s/%s%s", viper.GetString("cdn.url"), tdate, tempPath, ext)

	c.SaveUploadedFile(file, dst)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": url,
	})
}
