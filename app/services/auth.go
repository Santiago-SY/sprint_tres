package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type AuthLog struct {
	Timestamp string `json:"_time"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Action    string `json:"action"`  // Login, Logout, ResetPassword
	UserID    string `json:"user_id"` // Cardinalidad Media
	IP        string `json:"ip_addr"` // Cardinalidad Alta
	Status    string `json:"status"`  // Success, Failed
	Device    string `json:"device"`  // Mobile, Desktop
}

func RunAuthService(sender *client.LogSender) {
	actions := []string{"LOGIN", "LOGOUT", "REGISTER", "2FA_VERIFY"}
	devices := []string{"iPhone 14", "Pixel 7", "Windows Chrome", "Mac Safari"}

	fmt.Println(" Servicio de Auth: INICIADO")

	for {
		logData := AuthLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Service:   "auth-service",
			Level:     "INFO",
			Action:    actions[rand.Intn(len(actions))],
			UserID:    fmt.Sprintf("user-%d", rand.Intn(10000)), // Simula 10k usuarios
			IP:        fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
			Status:    "SUCCESS",
			Device:    devices[rand.Intn(len(devices))],
		}

		// Simular fallos de login (Seguridad)
		if rand.Float32() < 0.1 { // 10% fallos
			logData.Status = "FAILED"
			logData.Level = "WARN"
		}

		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(50 * time.Millisecond) // Menos tráfico que Pagos
	}
}
