package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nbtca/webhook-delivery-center/consolefixfunc"
	"github.com/nbtca/webhook-delivery-center/consumer"
	"github.com/nbtca/webhook-delivery-center/router"
	"github.com/nbtca/webhook-delivery-center/util"
)

func main() {
	util.InitDialer()
	consumer.InitConsumer()
	if runtime.GOOS == "windows" { //修复控制台上色
		err := consolefixfunc.EnableANSIConsole()
		if err != nil {
			fmt.Println("Error enabling ANSI console:", err)
		}
	}
	gin.SetMode(gin.ReleaseMode)
	err := util.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	r := gin.Default()
	r.Use(cors.Default())             //跨域
	r.GET("/", func(c *gin.Context) { //测试
		c.String(http.StatusOK, "200 ok")
	})
	router.InitWebhook(r)
	router.InitWs(r)
	if util.Cfg.UseCert {
		r.RunTLS(util.Cfg.Bind, util.Cfg.CertFile, util.Cfg.KeyFile) //启动服务
	} else {
		r.Run(util.Cfg.Bind) //启动服务
	}
}
