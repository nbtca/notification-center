package nsqclient

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nbtca/notification-center/service/mail"
	"github.com/nsqio/go-nsq"
)

type EventActionMessageHandler struct {
}

func (m *EventActionMessageHandler) HandleMessage(msg *nsq.Message) (err error) {
	mail.SendMessageViaMail("New Action", msg)
	fmt.Printf("recv from %v, msg:%v\n", msg.NSQDAddress, string(msg.Body))
	return

}

type LogMessageHandler struct {
}

func (m *LogMessageHandler) HandleMessage(msg *nsq.Message) (err error) {
	mail.SendMessageViaMail("New Action", msg)
	fmt.Printf("recv from %v, msg:%v\n", msg.NSQDAddress, string(msg.Body))
	return
}

func CreateConsumer(topic string, channel string, address string, handler nsq.Handler) (err error) {
	nsq_secret := os.Getenv("NSQ_SECRET")
	config := nsq.NewConfig()
	config.AuthSecret = nsq_secret
	config.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return
	}
	c.AddHandler(handler)

	if err := c.ConnectToNSQD(address); err != nil {
		return err
	}
	return nil
}

func InitConsumer() {
	log_topic := os.Getenv("LOG_TOPIC")
	event_topic := os.Getenv("EVENT_TOPIC")
	channel := os.Getenv("CHANNEL")
	nsq_host := os.Getenv("NSQ_HOST")
	if log_topic == "" {
		log.Fatalln("LOG_TOPIC is not set")
		return
	}
	if event_topic == "" {
		log.Fatalln("EVENT_TOPIC is not set")
		return
	}
	if channel == "" {
		log.Fatalln("CHANNEL is not set")
		return
	}
	if nsq_host == "" {
		log.Fatalln("NSQ_HOST is not set")
		return
	}

	err := CreateConsumer(log_topic, channel, nsq_host, &LogMessageHandler{})
	if err != nil {
		log.Fatalf("Init log consumer failed, err:%v\n", err)
		return
	}

	err = CreateConsumer(event_topic, channel, nsq_host, &EventActionMessageHandler{})
	if err != nil {
		log.Fatalf("Init event consumer failed, err:%v\n", err)
		return
	}
}
