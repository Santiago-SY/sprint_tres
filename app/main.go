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
	// 1. Logs (VictoriaLogs)
	victoriaURL := os.Getenv("VICTORIA_URL")
	if victoriaURL == "" {
		victoriaURL = "http://localhost:9428/insert/jsonline"
	}
	sender := client.NewLogSender(victoriaURL)
	go sender.Start()

	// 2. Base de Datos (PostgreSQL)
	fmt.Println(" Conectando a persistencia...")
	dbPool, err := client.InitDB()
	if err != nil {
		panic(err) // Si falla la DB, no podemos seguir
	}
	defer dbPool.Close()

	// 3. Iniciar Servicios
	fmt.Println(" Iniciando servicios...")

	go services.RunGatewayService(sender)
	go services.RunAuthService(sender)
	go services.RunPaymentService(sender)
	go services.RunRiskService(sender)
	go services.RunCartService(sender)
	go services.RunNotificationService(sender)

	// AQUÍ EL CAMBIO: Le pasamos la DB al servicio de productos
	go services.RunProductService(sender, dbPool)

	// 4. Esperar señal de salida (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println(" Apagando...")
}
