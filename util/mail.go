package util

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nsqio/go-nsq"
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	host := os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Default().Println("parse mail port failed, err:", err)
	}
	userName := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	dialer = gomail.NewDialer(host, port, userName, password)
}

func SendMail(message *gomail.Message) error {
	message.SetHeader("From", dialer.Username)
	return dialer.DialAndSend(message)
}

func generateTableHTML(jsonStr []byte) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		return "", err
	}
	tableHTML := "<table>"
	for key, value := range data {
		key = html.EscapeString(key)
		valueStr := html.EscapeString(fmt.Sprintf("%v", value))
		tableHTML += fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td></tr>",
			key, valueStr,
		)
	}
	tableHTML += "</table>"

	return tableHTML, nil
}

func SendMessageViaMail(subject string, msg *nsq.Message) error {
	m := gomail.NewMessage()
	receiverAddress := os.Getenv("MAIL_RECEIVER_ADDRESS")
	if receiverAddress == "" {
		return fmt.Errorf("MAIL_RECEIVER_ADDRESS is not set")
	}
	m.SetHeader("To", receiverAddress)
	m.SetHeader("Subject", subject)

	tableHTML, err := generateTableHTML(msg.Body)
	if err != nil {
		tableHTML = fmt.Sprintf(`<div>%s
	  </div>`, msg.Body)
	}
	m.SetBody("text/html", fmt.Sprintf(
		`<div>
		<div>
		<span style="padding-right:10px;">日志事件:</span>
		 %s
	  </div>
</div>`, tableHTML))

	if err := SendMail(m); err != nil {
		return err
	}
	return nil
}
