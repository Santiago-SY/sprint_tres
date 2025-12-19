package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type GatewayLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"http_status"`
	LatencyMs int    `json:"latency_ms"`
	UserAgent string `json:"user_agent"`
}

func RunGatewayService(sender *client.LogSender) {
	// TIER 1: Entrada Masiva (15 Hilos)
	concurrency := 15
	fmt.Printf("🌐 API GATEWAY: Iniciando simulacion de alto trafico (%d hilos)...\n", concurrency)

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{"/api/products", "/api/auth/login", "/api/cart", "/api/checkout"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
		"PostmanRuntime/7.32.0",
	}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA (SINE WAVE) ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				// Base muy rápida (20ms) para simular miles de requests
				baseSleep := 20 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------------------

				logData := GatewayLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "api-gateway",
					Level:     "INFO",
					Method:    methods[rand.Intn(len(methods))],
					Path:      paths[rand.Intn(len(paths))],
					Status:    200,
					LatencyMs: rand.Intn(500),
					UserAgent: userAgents[rand.Intn(len(userAgents))],
				}

				if rand.Float32() < 0.02 { // 2% de error
					logData.Status = 500
					logData.Level = "ERROR"
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {} // Bloqueo infinito
}
