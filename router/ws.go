package router

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nbtca/notification-center/util"
)

func InitWs(r *gin.Engine) {
	// 注册路由
	r.GET("/ws/*path", handleWs) //ws服务
}

var (
	clients   = make(map[*websocket.Conn]string) //已经链接的ws客户端，value是ws连接的路径
	clientsMu sync.Mutex                         //客户互斥锁
	upgrader  = websocket.Upgrader{              //升级ws请求用
		CheckOrigin: func(r *http.Request) bool { return true },
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			log.Println("WebSocket error:", status, reason)
		},
	}
)

// 从客户端收到消息
func handleWsMessage(message []byte) {
}

// 处理ws路径的请求
func handleWs(c *gin.Context) {
	path := c.Param("path")[1:]
	err := util.Auth(c, nil, &path) //鉴权
	if err != nil {
		log.Println("handleWs : Auth failed:", err)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	go func() {
		clientsMu.Lock()
		clients[conn] = path
		clientsMu.Unlock()
		defer func() { //出作用域删除客户端
			clientsMu.Lock()
			defer clientsMu.Unlock()
			delete(clients, conn)
			conn.Close()
		}()

		for { //循环读取ws消息
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			handleWsMessage(message)
		}

	}()
}

// broadcast message to all clients 广播消息给所有客户端
func broadcastMessage(path *string, message []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		// send to all if path is empty
		// otherwise only send message to clients with same path
		// and client whose path is empty will receive all messages
		// 如果路径为空则发送给所有客户端
		// 否则只发送给相同路径的客户端
		// 路径为空的客户端将接收所有消息
		if *path != "" {
			if clients[client] != "" {
				if clients[client] != *path {
					continue
				}
			}
		}
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil { //disconnect client if failed to send message 发送消息失败则断开客户端
			log.Println("Failed to send WebSocket message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
