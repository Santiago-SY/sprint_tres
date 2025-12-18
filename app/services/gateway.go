package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type GatewayLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Method    string `json:"method"`      // GET, POST
	Path      string `json:"path"`        // /api/v1/products
	Status    int    `json:"http_status"` // 200, 404, 500
	LatencyMs int    `json:"latency_ms"`  // Tiempo de respuesta
	UserAgent string `json:"user_agent"`  // Texto largo y repetitivo (Ideal para compresión)
}

func RunGatewayService(sender *client.LogSender) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{"/api/products", "/api/auth/login", "/api/cart", "/api/checkout"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
		"PostmanRuntime/7.32.0",
	}

	fmt.Println(" API Gateway: INICIADO")

	for {
		logData := GatewayLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Service:   "api-gateway",
			Level:     "INFO",
			Method:    methods[rand.Intn(len(methods))],
			Path:      paths[rand.Intn(len(paths))],
			Status:    200,
			LatencyMs: rand.Intn(500), // 0 a 500ms
			UserAgent: userAgents[rand.Intn(len(userAgents))],
		}

		if rand.Float32() < 0.05 { // 5% errores 500
			logData.Status = 500
			logData.Level = "ERROR"
		}

		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(5 * time.Millisecond) // Muy rápido (Alto tráfico)
	}
}
