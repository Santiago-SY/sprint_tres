package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type AuthLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Action    string `json:"action"` // Login, Logout, Refresh
	UserID    string `json:"user_id"`
	IP        string `json:"ip_address"`
}

func RunAuthService(sender *client.LogSender) {
	// TIER 2: Intencion (6 Hilos)
	concurrency := 6
	fmt.Printf("🔐 AUTH SERVICE: Simulando logins (%d hilos)...\n", concurrency)

	actions := []string{"LOGIN", "LOGIN", "LOGOUT", "REFRESH_TOKEN"}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				// Sleep medio (80ms)
				baseSleep := 80 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------

				logData := AuthLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "auth-service",
					Level:     "INFO",
					Action:    actions[rand.Intn(len(actions))],
					UserID:    fmt.Sprintf("user_%d", rand.Intn(1000)),
					IP:        fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
				}

				if logData.Action == "LOGIN" && rand.Float32() < 0.05 {
					logData.Level = "WARN" // Login fallido
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
