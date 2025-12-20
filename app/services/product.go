package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"sprint-tres/client"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductLog struct {
	Timestamp string  `json:"_time"`
	Service   string  `json:"service"`
	Level     string  `json:"level"`
	Msg       string  `json:"_msg"`
	SKU       string  `json:"sku"`
	Name      string  `json:"product_name"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	QueryTime string  `json:"query_time_ms"`
}

// Ahora recibe el Pool de Postgres
func RunProductService(sender *client.LogSender, db *pgxpool.Pool) {
	concurrency := 20
	fmt.Printf("📦 PRODUCT SERVICE: SQL Real Activo (%d hilos)\n", concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			ctx := context.Background()
			for {
				// Simulación de tráfico (ola)
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)
				baseSleep := 50 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))

				start := time.Now()

				// --- QUERY REAL A POSTGRES ---
				var sku, name, category string
				var price float64

				// Elegimos un producto al azar de la tabla
				err := db.QueryRow(ctx,
					"SELECT sku, name, category, price FROM products ORDER BY RANDOM() LIMIT 1").
					Scan(&sku, &name, &category, &price)

				duration := time.Since(start).Milliseconds()

				level := "INFO"
				if err != nil {
					level = "ERROR"
					sku = "DB_ERR"
					fmt.Printf("❌ Error DB: %v\n", err)
				}

				logData := ProductLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "product-service",
					Level:     level,
					Msg:       fmt.Sprintf("Product queried: %s (%s)", name, sku),
					SKU:       sku,
					Name:      name,
					Category:  category,
					Price:     price,
					QueryTime: fmt.Sprintf("%dms", duration),
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
