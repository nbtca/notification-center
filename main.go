package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nbtca/webhook-delivery-center/consolefixfunc"
)

func main() {
	if runtime.GOOS == "windows" {
		err := consolefixfunc.EnableANSIConsole()
		if err != nil {
			fmt.Println("Error enabling ANSI console:", err)
			os.Exit(1)
		}
	}
	gin.SetMode(gin.ReleaseMode)
	err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	r := gin.Default()
	r.Use(cors.Default())             //跨域
	r.GET("/", func(c *gin.Context) { //测试
		c.String(http.StatusOK, "200 ok")
	})
	initWebhook(r)
	initWs(r)
	if cfg.UseCert {
		r.RunTLS(cfg.Bind, cfg.CertFile, cfg.KeyFile) //启动服务
	} else {
		r.Run(cfg.Bind) //启动服务
	}
}
