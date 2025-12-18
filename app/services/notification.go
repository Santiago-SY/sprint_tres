package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type NotifLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Channel   string `json:"channel"` // EMAIL, SMS, PUSH
	Template  string `json:"template"`
	To        string `json:"recipient_mask"` // x***@gmail.com
}

func RunNotificationService(sender *client.LogSender) {
	channels := []string{"EMAIL", "SMS", "PUSH"}
	fmt.Println("📨 Notification Service: INICIADO")

	for {
		logData := NotifLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Service:   "notification-svc",
			Level:     "INFO",
			Channel:   channels[rand.Intn(len(channels))],
			Template:  "order_confirmation",
			To:        "u***@provider.com",
		}
		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(60 * time.Millisecond)
	}
}
