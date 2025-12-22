package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type PaymentLog struct {
	Timestamp string  `json:"_time"`
	Service   string  `json:"service"`
	Level     string  `json:"level"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Gateway   string  `json:"gateway"`
	Status    string  `json:"status"`
	TraceID   string  `json:"trace_id"`
}

func RunPaymentService(sender *client.LogSender) {
	// TIER 3: Conversion Real (2 Hilos) - Pocos pero importantes
	concurrency := 2
	fmt.Printf("PAYMENT API: Procesando pagos (%d hilos)...\n", concurrency)

	gateways := []string{"Stripe", "PayPal", "MercadoPago"}
	currencies := []string{"USD", "EUR", "UYU"}
	statuses := []string{"SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "SUCCESS", "FAILED"}

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			for {
				// --- ALGORITMO DE OLA ---
				cycleSeconds := 60.0
				nowUnix := float64(time.Now().UnixNano()) / 1e9
				wave := 1.0 + math.Sin(2*math.Pi*nowUnix/cycleSeconds)

				// Sleep Lento (200ms) para que haya menos pagos que visitas
				baseSleep := 200 * time.Millisecond
				dynamicSleep := time.Duration(float64(baseSleep) / (0.1 + wave))
				// ------------------------

				logData := PaymentLog{
					Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
					Service:   "payment-api",
					Level:     "INFO",
					Amount:    float64(rand.Intn(10000)) / 100.0,
					Currency:  currencies[rand.Intn(len(currencies))],
					Gateway:   gateways[rand.Intn(len(gateways))],
					Status:    statuses[rand.Intn(len(statuses))],
					TraceID:   fmt.Sprintf("trace-%d-%d", id, rand.Int63()),
				}

				if logData.Status == "FAILED" {
					logData.Level = "ERROR"
				}

				jsonData, _ := json.Marshal(logData)
				sender.Enqueue(jsonData)
				time.Sleep(dynamicSleep)
			}
		}(i)
	}
	select {}
}
