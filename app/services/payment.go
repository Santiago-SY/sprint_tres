package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"sprint-tres/client" // Importamos nuestro motor de envío
)

// PaymentLog define la estructura de un evento de pago.
// Los "json tags" (lo que está entre comillas `...`) definen cómo se verá en la base de datos.
type PaymentLog struct {
	Timestamp string  `json:"_time"`    // Campo especial: VictoriaLogs usa _time por defecto
	Service   string  `json:"service"`  // Para filtrar por servicio: _stream:{service="payment"}
	Level     string  `json:"level"`    // INFO, ERROR, WARN
	Amount    float64 `json:"amount"`   // Monto (numérico para hacer sumas/promedios)
	Currency  string  `json:"currency"` // USD, EUR, UYU
	Gateway   string  `json:"gateway"`  // Stripe, PayPal
	Status    string  `json:"status"`   // SUCCESS, FAILED
	TraceID   string  `json:"trace_id"` // Para rastrear una petición única
}

// RunPaymentService es el bucle infinito que simula la vida del microservicio.
func RunPaymentService(sender *client.LogSender) {
	// Datos semilla para aleatoriedad
	gateways := []string{"Stripe", "PayPal", "MercadoPago"}
	currencies := []string{"USD", "EUR", "UYU"}
	statuses := []string{"SUCCESS", "SUCCESS", "SUCCESS", "FAILED"} // 75% de probabilidad de éxito

	fmt.Println("💳 Servicio de Pagos: INICIADO")

	for {
		// 1. Simulación de Negocio: Crear el dato
		now := time.Now().UTC()
		logData := PaymentLog{
			Timestamp: now.Format(time.RFC3339Nano), // Formato estándar ISO8601
			Service:   "payment-api",
			Level:     "INFO",
			Amount:    float64(rand.Intn(50000)) / 100.0, // Genera monto entre 0.00 y 500.00
			Currency:  currencies[rand.Intn(len(currencies))],
			Gateway:   gateways[rand.Intn(len(gateways))],
			Status:    statuses[rand.Intn(len(statuses))],
			TraceID:   fmt.Sprintf("trace-%d", rand.Int63()), // ID único random
		}

		// Si falló el pago, cambiamos el nivel a ERROR (para que se vea rojo en Grafana)
		if logData.Status == "FAILED" {
			logData.Level = "ERROR"
		}

		// 2. Serialización: Convertir Struct de Go -> JSON Bytes
		// json.Marshal es rápido y seguro.
		jsonData, _ := json.Marshal(logData)

		// 3. Envío NO Bloqueante: ¡Aquí usamos el motor!
		// Tiramos el JSON a la cinta transportadora y nos olvidamos.
		sender.Enqueue(jsonData)

		// 4. Ritmo de Tráfico
		// Dormimos un poco para no saturar tu PC local.
		// En un entorno real, esto dependería de los usuarios reales.
		// 10ms = ~100 logs/segundo solo de este servicio.
		time.Sleep(10 * time.Millisecond)
	}
}
