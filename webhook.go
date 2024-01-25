package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func initWebhook(r *gin.Engine) {
	// 注册路由
	r.POST("/webhook/*path", handleWebhook) //webhook服务
}

type GithubWebhookPost struct {
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    interface{}         `json:"body"`
}

// 处理webhook请求(POST)
func handleWebhook(c *gin.Context) {
	var body interface{}
	if err := c.ShouldBindJSON(&body); err != nil { //绑定内容json到结构体
		log.Println(err)
		return
	}
	headers := make(map[string][]string)
	for k, v := range c.Request.Header {
		headers[k] = v
	}
	path := c.Param("path")[1:]
	go func() {
		fulldata := GithubWebhookPost{
			Headers: headers, //请求头
			Body:    body,    //原始内容
			Path:    path,    //路径
		}
		jsonData, err := json.Marshal(fulldata)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(jsonData))
		broadcastMessage(&path, jsonData) //广播消息
	}()
}
