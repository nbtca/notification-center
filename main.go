package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	loadConfig()
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
