package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type NotificationLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Channel   string `json:"channel"` // EMAIL, SMS
	Type      string `json:"type"`    // ORDER_CONFIRMED, OTP
	Recipient string `json:"recipient"`
}

func RunNotificationService(sender *client.LogSender) {
	// TIER 3: Acompaña a pagos (2 Hilos)
	concurrency := 2
	fmt.Printf("NOTIFICATION SVC: Enviando alertas (%d hilos)...\n", concurrency)

	channels := []string{"EMAIL", "EMAIL", "SMS"}
	types := []string{"ORDER_CONFIRMED", "PAYMENT_RECEIVED", "OTP"}

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

				logData := NotificationLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "notification-svc",
					Level:     "INFO",
					Channel:   channels[rand.Intn(len(channels))],
					Type:      types[rand.Intn(len(types))],
					Recipient: fmt.Sprintf("user_%d@example.com", rand.Intn(999)),
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
