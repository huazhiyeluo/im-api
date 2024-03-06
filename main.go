package main

import (
	"demoapi/models"
	"demoapi/router"
	"demoapi/utils"
	"time"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server")) // listen and serve on 0.0.0.0:8080
}

func InitTimer() {
	utils.Timer(time.Second*3, time.Second*3, models.CleanConnection)
}
