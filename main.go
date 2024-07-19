package main

import (
	"qqapi/internal/server"
	"qqapi/internal/utils"
	"qqapi/router"
	"qqapi/third_party/log"
	"time"

	"github.com/spf13/viper"
)

func main() {
	log.InitLogger()
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	log.Logger.Info("Starting")

	go server.App()
	InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server")) // listen and serve on 0.0.0.0:8080
}

func InitTimer() {
	Timer(time.Second*3, time.Second*3, server.CleanConnection)
}

func Timer(delay time.Duration, tick time.Duration, fun func()) {
	go func() {
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				// 定时器触发的处理逻辑
				fun()
				t.Reset(tick)
			}
		}
	}()
}
