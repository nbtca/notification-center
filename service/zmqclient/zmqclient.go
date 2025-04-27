package zmqclient

import (
	"fmt"
	"log"
	"sync"

	zmq "github.com/go-zeromq/zmq4"
	"github.com/nbtca/notification-center/util"
	"golang.org/x/net/context"
)

var (
	ctx context.Context
	pub zmq.Socket
	sub zmq.Socket
	mu  sync.Mutex
)

// InitZeroMQ 初始化ZeroMQ
func InitZeroMQ() error {
	if !util.Cfg.ZeroMQ.Enabled {
		log.Println("ZeroMQ is disabled")
		return nil
	}

	pattern := util.Cfg.ZeroMQ.Pattern
	if pattern == "" {
		pattern = "PUB/SUB" // 默认使用PUB/SUB模式
	}

	var err error
	switch pattern {
	case "PUB/SUB":
		err = initPubSub()
	case "REQ/REP":
		err = initReqRep()
	case "PUSH/PULL":
		err = initPushPull()
	default:
		return fmt.Errorf("不支持的ZeroMQ模式: %s", pattern)
	}

	if err != nil {
		return err
	}

	log.Printf("ZeroMQ已初始化，模式: %s", pattern)
	return nil
}

// 初始化PUB/SUB模式
func initPubSub() error {
	var err error
	// 创建发布者
	if util.Cfg.ZeroMQ.PubAddress != "" {
		ctx = context.Background()
		pub = zmq.NewPub(ctx)
		err = pub.Listen(util.Cfg.ZeroMQ.PubAddress)
		if err != nil {
			return fmt.Errorf("绑定PUB地址失败: %v", err)
		}
		log.Printf("ZeroMQ PUB已绑定到 %s", util.Cfg.ZeroMQ.PubAddress)
	}

	// 创建订阅者
	if util.Cfg.ZeroMQ.SubAddress != "" {
		sub = zmq.NewSub(ctx)
		err = sub.Dial(util.Cfg.ZeroMQ.SubAddress)
		if err != nil {
			return fmt.Errorf("连接SUB地址失败: %v", err)
		}
		log.Printf("ZeroMQ SUB已连接到 %s", util.Cfg.ZeroMQ.SubAddress)

		// 启动接收消息的goroutine
		go receiveMessages()
	}

	return nil
}

// 初始化REQ/REP模式
func initReqRep() error {
	// 实现REQ/REP模式的初始化逻辑
	return fmt.Errorf("REQ/REP模式尚未实现")
}

// 初始化PUSH/PULL模式
func initPushPull() error {
	// 实现PUSH/PULL模式的初始化逻辑
	return fmt.Errorf("PUSH/PULL模式尚未实现")
}

// 接收消息的goroutine
func receiveMessages() {
	if sub == nil {
		return
	}

	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Printf("接收ZeroMQ消息错误: %v", err)
			continue
		}

		// 处理收到的消息
		log.Printf("收到ZeroMQ消息: %v", msg)
		// 这里可以添加自定义的消息处理逻辑
	}
}

// PublishMessage 发布消息
func PublishMessage(topic string, message string) error {
	if !util.Cfg.ZeroMQ.Enabled || pub == nil {
		return fmt.Errorf("ZeroMQ未启用或发布者未初始化")
	}

	mu.Lock()
	defer mu.Unlock()

	// 创建消息并发送
	msg := zmq.NewMsgFrom([]byte(topic), []byte(message))
	err := pub.Send(msg)
	if err != nil {
		return fmt.Errorf("发送ZeroMQ消息失败: %v", err)
	}

	return nil
}

// CloseZeroMQ 关闭ZeroMQ连接
func CloseZeroMQ() {
	if pub != nil {
		pub.Close()
	}
	if sub != nil {
		sub.Close()
	}
	log.Println("ZeroMQ已关闭")
}
