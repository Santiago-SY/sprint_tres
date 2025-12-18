package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client"
)

type RiskLog struct {
	Timestamp     string   `json:"_time"`
	Service       string   `json:"service"`
	Level         string   `json:"level"`
	RiskScore     float64  `json:"risk_score"`    // 0.0 a 1.0
	Decision      string   `json:"decision"`      // APPROVE, REJECT, CHALLENGE
	RulesChecked  []string `json:"rules_checked"` // Array: ["ip_geo", "velocity"]
	TransactionID string   `json:"tx_id"`
}

func RunRiskService(sender *client.LogSender) {
	fmt.Println("🕵️  Risk Engine: INICIADO")

	for {
		score := rand.Float64()
		decision := "APPROVE"
		if score > 0.8 {
			decision = "REJECT"
		} else if score > 0.5 {
			decision = "CHALLENGE"
		}

		logData := RiskLog{
			Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
			Service:       "risk-engine",
			Level:         "INFO",
			RiskScore:     score,
			Decision:      decision,
			RulesChecked:  []string{"ip_velocity", "geo_fencing", "blacklists"}, // Array complejo
			TransactionID: fmt.Sprintf("tx-%d", rand.Int63()),
		}

		if decision == "REJECT" {
			logData.Level = "WARN"
		}

		jsonData, _ := json.Marshal(logData)
		sender.Enqueue(jsonData)
		time.Sleep(20 * time.Millisecond)
	}
}
