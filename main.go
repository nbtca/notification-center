package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func handleWsMessage(message []byte) {
}

func ws(c *gin.Context) {
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

func webhook(c *gin.Context) {
	go func() {
		var msg map[string]interface{}
		if err := c.ShouldBindJSON(&msg); err != nil {
			log.Println(err)
			return
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}

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
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "test passed")
	})

	r.POST("/webhook", webhook)
	r.GET("/ws", ws)

	r.Run(":8000")
}
