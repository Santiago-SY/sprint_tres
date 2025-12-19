package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type CartLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Operation string `json:"operation"` // ADD, REMOVE, CLEAR
	ItemID    string `json:"item_id"`
	CartID    string `json:"cart_id"`
}

func RunCartService(sender *client.LogSender) {
	// TIER 2: Intencion (6 Hilos)
	concurrency := 6
	fmt.Printf("🛒 CART SERVICE: Simulando carritos activos (%d hilos)...\n", concurrency)

	ops := []string{"ADD_ITEM", "ADD_ITEM", "ADD_ITEM", "REMOVE_ITEM", "CLEAR_CART"}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				baseSleep := 80 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------

				logData := CartLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "cart-service",
					Level:     "INFO",
					Operation: ops[rand.Intn(len(ops))],
					ItemID:    fmt.Sprintf("prod_%d", rand.Intn(500)),
					CartID:    fmt.Sprintf("cart_%d", rand.Intn(10000)),
				}
				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
