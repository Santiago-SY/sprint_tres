package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sprint-tres/client"
	"sprint-tres/services"
)

func main() {
	// --- PASO 1: CONFIGURACIÓN ---
	victoriaURL := os.Getenv("VICTORIA_URL")
	if victoriaURL == "" {
		victoriaURL = "http://localhost:9428/insert/jsonline"
	}

	fmt.Printf("\n🚀 INICIANDO LOG GENERATOR (7 Microservicios)\n")
	fmt.Printf("🎯 Objetivo: %s\n", victoriaURL)

	// --- PASO 2: ARRANCAR MOTOR ---
	sender := client.NewLogSender(victoriaURL)
	sender.Start()

	// --- PASO 3: ARRANCAR SERVICIOS (CONCURRENCIA REAL) ---
	fmt.Println("🚦 Despertando flota de servicios...")

	// EXPLICACIÓN TÉCNICA (Para tu defensa):
	// Usamos la keyword 'go' para lanzar cada función en una Goroutine separada.
	// Si quitáramos el 'go', el programa se quedaría atrapado en el bucle infinito
	// de RunGatewayService y nunca arrancaría los demás.
	// Esto demuestra el modelo de concurrencia M:N de Go.

	go services.RunGatewayService(sender)      // Mucho tráfico
	go services.RunAuthService(sender)         // IDs de usuarios
	go services.RunPaymentService(sender)      // Dinero
	go services.RunRiskService(sender)         // JSON Complejo (Schema-less)
	go services.RunCartService(sender)         // Estado
	go services.RunProductService(sender)      // Catálogo
	go services.RunNotificationService(sender) // Emails

	// --- PASO 4: ESPERA ACTIVA ---
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Señal de parada recibida. Apagando sistema...")
}
