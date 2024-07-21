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

type ClientInfo struct {
	path    string
	headers map[string][]string
}

var (
	clients   = make(map[*websocket.Conn]*ClientInfo) //已经链接的ws客户端，value是ws连接的路径
	clientsMu sync.Mutex                              //客户互斥锁
	upgrader  = websocket.Upgrader{                   //升级ws请求用
		CheckOrigin: func(r *http.Request) bool { return true },
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			log.Println("WebSocket error:", status, reason)
		},
	}
)

// 从客户端收到消息
func handleWsMessage(conn *websocket.Conn, message []byte) {
	// 转发消息给所有客户端
	info := clients[conn]
	if info == nil {
		log.Println("handleWsMessage: client info not found"+conn.RemoteAddr().String(), conn)
		return
	}
	broadcastMessage(&info.path, message, conn)
	//print message
	log.Println("ws received message from ", info, ":", string(message))
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
	headers := c.Request.Header
	delete(headers, "Authorization")
	delete(headers, "Upgrade")
	delete(headers, "Connection")
	delete(headers, "Sec-Websocket-Key")
	delete(headers, "Sec-Websocket-Version")
	info := &ClientInfo{path: path, headers: headers}
	go func() {
		clientsMu.Lock()
		clients[conn] = info
		clientsMu.Unlock()
		broadcastActiveClientsChange(&info.path)
		defer func() { //出作用域删除客户端
			disconnectClient(conn, &info.path)
		}()
		for { //循环读取ws消息
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			handleWsMessage(conn, message)
		}
	}()
}

// broadcast message to all clients 广播消息给所有客户端
func broadcastMessage(path *string, message []byte, excluedeConn *websocket.Conn) {
	for client, info := range clients {
		if client == excluedeConn {
			continue
		}
		if notSameCategory(path, info) {
			continue
		}
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil { //disconnect client if failed to send message 发送消息失败则断开客户端
			log.Println("Failed to send WebSocket message:", err)
			disconnectClient(client, &info.path)
		}
	}
}
func notSameCategory(path *string, info *ClientInfo) bool {
	// send to all if path is empty
	// otherwise only send message to clients with same path
	// and client whose path is empty will receive all messages
	// 如果路径为空则发送给所有客户端
	// 否则只发送给相同路径的客户端
	// 路径为空的客户端将接收所有消息
	return *path != "" && (info == nil || (info.path != *path && info.path != ""))
}

type PacketSourceInfo struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Version     string `json:"version"`
}
type ActiveBroadcastPacketData struct {
	Clients []ActiveBroadcastPacketDataClient `json:"clients"`
}
type ActiveBroadcastPacketDataClient struct {
	Address string              `json:"address"`
	Headers map[string][]string `json:"headers"`
}

type ActiveBroadcastPacket struct {
	Type   string                    `json:"type"`
	Source PacketSourceInfo          `json:"source"`
	Data   ActiveBroadcastPacketData `json:"data"`
}

func broadcastActiveClientsChange(path *string) {
	pkt := &ActiveBroadcastPacket{
		Type: "active_clients_change",
		Source: PacketSourceInfo{
			"notification-center",
			"notification-center",
			"1.0.0",
		},
		Data: ActiveBroadcastPacketData{
			Clients: make([]ActiveBroadcastPacketDataClient, len(clients)),
		},
	}
	index := 0
	for client, info := range clients {
		if notSameCategory(path, info) {
			continue
		}
		if *path != "" && info != nil && info.path != "" && info.path != *path {
			continue
		}
		pkt.Data.Clients[index] = ActiveBroadcastPacketDataClient{
			Address: client.RemoteAddr().String(),
			Headers: info.headers,
		}
		index++
	}
	log.Println("client count", len(clients), " for path ", *path)
	for client, info := range clients {
		if notSameCategory(path, info) {
			continue
		}
		err := client.WriteJSON(pkt)
		if err != nil {
			log.Println("Failed to send WebSocket message:", err)
			disconnectClient(client, &info.path)
		}
	}
}

func disconnectClient(client *websocket.Conn, path *string) {
	func() {
		log.Println("disconnectClient", client.RemoteAddr().String(), *path)
		clientsMu.Lock()
		defer clientsMu.Unlock()
		delete(clients, client)
		client.Close()
	}()
	broadcastActiveClientsChange(path)
}
