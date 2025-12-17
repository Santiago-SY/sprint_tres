package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sprint-tres/client"   // Tu motor
	"sprint-tres/services" // Tus microservicios
)

func main() {
	// --- PASO 1: CONFIGURACIÓN ---
	// Leemos la URL desde las variables de entorno (definidas en docker-compose.yml)
	// Si no existe, usamos localhost (útil para probar fuera de Docker).
	victoriaURL := os.Getenv("VICTORIA_URL")
	if victoriaURL == "" {
		victoriaURL = "http://localhost:9428/insert/jsonline"
	}

	fmt.Printf("\n🚀 INICIANDO LOG GENERATOR\n")
	fmt.Printf("🎯 Objetivo: %s\n", victoriaURL)

	// --- PASO 2: ARRANCAR MOTOR ---
	// Instanciamos el "Camión de Mudanza" y lo encendemos.
	sender := client.NewLogSender(victoriaURL)
	sender.Start()

	// --- PASO 3: ARRANCAR SERVICIOS (CONCURRENCIA) ---
	fmt.Println("🚦 Despertando Microservicios...")

	// Lanzamos el servicio de Pagos en su propia Goroutine (hilo ligero).
	// El 'go' al principio significa: "Ejecuta esto en paralelo y sigue bajando".
	go services.RunPaymentService(sender)

	// (Aquí descomentaremos los otros servicios a medida que los creemos)
	// go services.RunAuthService(sender)
	// go services.RunGatewayService(sender)
	// ...

	// --- PASO 4: ESPERA ACTIVA (GRACEFUL SHUTDOWN) ---
	// Si el programa termina aquí, todo se apaga instantáneamente.
	// Necesitamos bloquear la ejecución hasta que alguien quiera salir.

	// Creamos un canal para escuchar señales del Sistema Operativo (Ctrl+C o Docker Stop)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// El programa se queda "congelado" en esta línea esperando la señal.
	<-c

	fmt.Println("\n🛑 Señal de parada recibida. Apagando sistema...")
	// Aquí el programa termina y Go limpia la memoria.
}
