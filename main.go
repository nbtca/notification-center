package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nbtca/notification-center/router"
	"github.com/nbtca/notification-center/service/zmqclient"
	"github.com/nbtca/notification-center/util"
	"github.com/nbtca/notification-center/util/consolefixfunc"
)

func main() {
	// util.InitDialer()
	// consumer.InitConsumer()
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

	// 初始化ZeroMQ
	err = zmqclient.InitZeroMQ()
	if err != nil {
		fmt.Println("Error initializing ZeroMQ:", err)
		// 这里不退出程序，让程序继续运行，只是ZeroMQ功能不可用
	}
	// 确保在程序退出时关闭ZeroMQ连接
	defer zmqclient.CloseZeroMQ()
	// http服务器
	r := gin.Default()
	r.Use(cors.Default())             //跨域
	r.GET("/", func(c *gin.Context) { //测试
		c.String(http.StatusOK, "200 ok")
	})
	fmt.Println("Server started on ", util.Cfg.Bind)
	router.InitWebhook(r)
	router.InitWs(r)
	if util.Cfg.UseCert {
		r.RunTLS(util.Cfg.Bind, util.Cfg.CertFile, util.Cfg.KeyFile) //启动服务
	} else {
		r.Run(util.Cfg.Bind) //启动服务
	}
}
