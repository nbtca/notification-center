package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) //已经链接的ws客户端
	clientsMu sync.Mutex                       //客户端锁
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// 从客户端收到消息
func handleWsMessage(message []byte) {
}

// 处理ws路径的请求
func handleWs(c *gin.Context) {
	go func() {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}
		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			handleWsMessage(message)
		}
	}()
}

// 处理webhook请求(POST)
func handleWebhook(c *gin.Context) {
	go func() {
		// var msg map[string]interface{}
		// if err := c.ShouldBindJSON(&msg); err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// jsonData, err := json.Marshal(msg)
		jsonData, err := c.GetRawData()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(jsonData))
		clientsMu.Lock()
		defer clientsMu.Unlock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Failed to send WebSocket message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}()
}
func main() {
	r := gin.Default()
	r.Use(cors.Default())             //跨域
	r.GET("/", func(c *gin.Context) { //测试
		c.String(http.StatusOK, "test passed")
	})
	r.POST("/webhook", handleWebhook) //webhook服务
	r.GET("/ws", handleWs)            //ws服务
	r.Run(":8000")                    //启动服务
}
