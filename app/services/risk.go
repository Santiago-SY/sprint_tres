package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type RiskLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Score     int    `json:"risk_score"` // 0-100
	Decision  string `json:"decision"`   // APPROVED, REJECTED
	TraceID   string `json:"trace_id"`
}

func RunRiskService(sender *client.LogSender) {
	// TIER 3: Acompaña a pagos (2 Hilos)
	concurrency := 2
	fmt.Printf("🛡️ RISK ENGINE: Analizando fraude (%d hilos)...\n", concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				baseSleep := 200 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------

				score := rand.Intn(100)
				decision := "APPROVED"
				level := "INFO"

				if score > 85 {
					decision = "REJECTED"
					level = "WARN"
				}

				logData := RiskLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "risk-engine",
					Level:     level,
					Score:     score,
					Decision:  decision,
					TraceID:   fmt.Sprintf("trace-%d-%d", id, rand.Int63()),
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
