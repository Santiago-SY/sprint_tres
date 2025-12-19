package services

import (
	"encoding/json"
	"fmt"
	"math"
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
	QueryTime string `json:"query_time_ms"`
}

func RunProductService(sender *client.LogSender) {
	// TIER 1: Catalog View (15 Hilos) - La gente mira mucho producto
	concurrency := 15
	fmt.Printf("📦 PRODUCT SERVICE: Iniciando simulacion de catalogo (%d hilos)...\n", concurrency)

	categories := []string{"Electronics", "Clothing", "Home", "Toys"}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				baseSleep := 20 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------

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
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
