package main

import (
	"demoapi/router"
	"demoapi/utils"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	r := router.Router()
	r.Run(viper.GetString("port.server")) // listen and serve on 0.0.0.0:8080
}
