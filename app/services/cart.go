package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type CartLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	CartID    string `json:"cart_id"`
	ItemCount int    `json:"item_count"`
	Action    string `json:"action"` // ADD, REMOVE, CLEAR
}

func RunCartService(sender *client.LogSender) {
	actions := []string{"ADD_ITEM", "REMOVE_ITEM", "CLEAR_CART"}
	fmt.Println("🛒 Cart Service: INICIADO")

	for {
		logData := CartLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Service:   "cart-service",
			Level:     "INFO",
			CartID:    fmt.Sprintf("cart-%d", rand.Intn(5000)),
			ItemCount: rand.Intn(10) + 1,
			Action:    actions[rand.Intn(len(actions))],
		}
		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(30 * time.Millisecond)
	}
}
