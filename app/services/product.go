package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type ProductLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	SKU       string `json:"sku"`
	Category  string `json:"category"`
	QueryTime string `json:"query_time_ms"` // Texto para variar "12ms"
}

func RunProductService(sender *client.LogSender) {
	categories := []string{"Electronics", "Clothing", "Home", "Toys"}
	fmt.Println(" Product Service: INICIADO")

	for {
		logData := ProductLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Service:   "product-service",
			Level:     "INFO",
			SKU:       fmt.Sprintf("SKU-%d", rand.Intn(9999)),
			Category:  categories[rand.Intn(len(categories))],
			QueryTime: fmt.Sprintf("%dms", rand.Intn(100)),
		}
		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(15 * time.Millisecond)
	}
}
