package router

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nbtca/notification-center/util"
)

func InitWebhook(r *gin.Engine) {
	// 注册路由
	r.POST("/*path", handleWebhook) //webhook服务
}

type GithubWebhookPost struct {
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    interface{}         `json:"body"`
}

// 处理webhook请求(POST)
func handleWebhook(c *gin.Context) {
	path := c.Param("path")[1:]
	bodyData, err := c.GetRawData()
	if err != nil {
		err := fmt.Errorf("auth failed, get body failed in path '%s'", path)
		c.AbortWithError(401, err)
		return
	}
	err = util.Auth(c, &bodyData, &path) //鉴权
	if err != nil {
		log.Println("handleWebhook : Auth failed:", err)
		return
	}
	var body interface{}
	if err := json.Unmarshal(bodyData, &body); err != nil { //绑定内容json到结构体
		log.Println(err)
		return
	}
	headers := c.Request.Header
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
		broadcastMessage(&path, jsonData, nil) //广播消息
	}()
}
