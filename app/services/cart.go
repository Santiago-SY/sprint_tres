package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"sprint-tres/client"

	"github.com/redis/go-redis/v9" // Driver de Redis/Valkey
)

// CartLog corregido para VictoriaLogs
type CartLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Msg       string `json:"_msg"` // <--- ESTO EVITA LOS ERRORES ROJOS
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Product   string `json:"product_sku"`
	Latency   string `json:"latency_ms"`
}

func RunCartService(sender *client.LogSender) {
	// 1. Conexión a Valkey
	valkeyHost := os.Getenv("VALKEY_HOST")
	if valkeyHost == "" {
		valkeyHost = "valkey:6379"
	}

	fmt.Printf("🛒 CART: Conectando a %s...\n", valkeyHost)

	rdb := redis.NewClient(&redis.Options{
		Addr: valkeyHost,
		DB:   0,
	})

	// Ping de prueba
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		fmt.Printf("❌ CART ERROR: No veo a Valkey: %v\n", err)
	} else {
		fmt.Println("🛒 CART: ¡Conectado a Valkey con éxito! 🚀")
	}

	concurrency := 10
	products := []string{"APL-IP15PM", "SON-PS5-SL", "SAM-S24U", "NIKE-AF1", "DYS-V15"}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// Simulación de tráfico
				cycleSeconds := 30.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)
				baseSleep := 100 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))

				userID := fmt.Sprintf("user_%d", rand.Intn(500))
				product := products[rand.Intn(len(products))]
				quantity := rand.Intn(3) + 1

				// --- OPERACIÓN VALKEY ---
				start := time.Now()
				key := fmt.Sprintf("cart:%s", userID)

				// Guardamos en Hash: cart:user_1 -> SKU -> Cantidad
				err := rdb.HSet(ctx, key, product, quantity).Err()
				rdb.Expire(ctx, key, 24*time.Hour) // TTL de 1 día

				duration := time.Since(start).Microseconds() // ¡Microsegundos!

				level := "INFO"
				msg := fmt.Sprintf("Added %d x %s to cart", quantity, product)
				if err != nil {
					level = "ERROR"
					msg = fmt.Sprintf("Valkey Write Failed: %v", err)
				}

				// Log JSON
				logData := CartLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "cart-service",
					Level:     level,
					Msg:       msg,
					UserID:    userID,
					Action:    "ADD",
					Product:   product,
					Latency:   fmt.Sprintf("%dµs", duration), // µs = ultra rápido
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)

				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
