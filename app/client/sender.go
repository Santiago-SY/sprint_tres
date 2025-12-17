package client

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CONFIGURACIÓN DEL MOTOR
// Ajustamos estos valores para equilibrar latencia vs. rendimiento.
const (
	BatchSize     = 1000            // "La Caja": Enviamos cuando juntamos 1000 logs.
	FlushInterval = 1 * time.Second // "El Reloj": Si no se llena la caja, enviamos cada 1 seg.
)

// LogSender es nuestro "Camión de Mudanza"
type LogSender struct {
	url        string          // Dirección de VictoriaLogs (ej: http://victoria-logs:9428...)
	logChannel chan []byte     // LA CINTA TRANSPORTADORA (Buffer en memoria)
	wg         *sync.WaitGroup // Semáforo para esperar a que los envíos terminen al cerrar
}

// NewLogSender inicializa la fábrica
func NewLogSender(url string) *LogSender {
	return &LogSender{
		url: url,
		// Creamos el Buffer de 10,000 espacios.
		// Esto absorbe los picos de tráfico sin bloquear a los microservicios.
		logChannel: make(chan []byte, 10000),
		wg:         &sync.WaitGroup{},
	}
}

// Start arranca el motor en segundo plano (Goroutine)
// Se llama una sola vez al inicio del programa.
func (s *LogSender) Start() {
	s.wg.Add(1)
	go s.runBatcher() // ¡Aquí nace la concurrencia!
}

// Enqueue es el método público que usan los servicios (Payment, Auth, etc.)
// Es "Non-blocking": Tira el dato y retorna inmediatamente.
func (s *LogSender) Enqueue(logJSON []byte) {
	select {
	case s.logChannel <- logJSON:
		// Éxito: El log entró en la cinta.
	default:
		// Fallo: La cinta está llena (Backpressure / Load Shedding).
		// Preferimos perder un log que colgar el servicio de Pagos.
		// En un sistema real, aquí incrementaríamos una métrica de "dropped_logs".
		fmt.Println("Warning: Buffer lleno, descartando log para proteger el sistema")
	}
}

// runBatcher es el "Trabajador de Almacén"
// Recoge logs de la cinta y arma los paquetes HTTP.
func (s *LogSender) runBatcher() {
	defer s.wg.Done()

	var batch [][]byte                      // Nuestra "Caja" temporal
	ticker := time.NewTicker(FlushInterval) // El reloj que suena cada 1 seg
	defer ticker.Stop()

	for {
		select {
		// CASO A: Llega un log por la cinta
		case log, ok := <-s.logChannel:
			if !ok {
				// Si el canal se cerró (apagado del sistema), enviamos lo que quede.
				if len(batch) > 0 {
					s.sendBatch(batch)
				}
				return
			}
			batch = append(batch, log) // Metemos el log a la caja

			// Si la caja se llenó (1000 logs), la enviamos YA.
			if len(batch) >= BatchSize {
				s.sendBatch(batch)
				batch = nil // Preparamos caja nueva limpia
			}

		// CASO B: Pasó 1 segundo (El reloj sonó)
		case <-ticker.C:
			// Si hay algo en la caja, lo enviamos aunque no esté llena.
			// Garantiza que los logs aparezcan en Grafana rápido.
			if len(batch) > 0 {
				s.sendBatch(batch)
				batch = nil
			}
		}
	}
}

// sendBatch hace el envío HTTP real a VictoriaLogs
func (s *LogSender) sendBatch(batch [][]byte) {
	// 1. Unimos todos los JSONs con saltos de línea (Formato JSON Stream)
	// VictoriaLogs espera: {json1}\n{json2}\n{json3}
	payload := bytes.Join(batch, []byte("\n"))

	// 2. Configuración de Cliente con Timeout
	// Si VictoriaLogs tarda más de 5s, cortamos para no quedarnos colgados.
	client := http.Client{Timeout: 5 * time.Second}

	// 3. Hacemos el POST
	resp, err := client.Post(s.url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error de Red enviando batch: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// (Opcional) Validar respuesta
	// VictoriaLogs devuelve 204 No Content o 200 OK si todo salió bien.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		fmt.Printf("Error del servidor VictoriaLogs: %d\n", resp.StatusCode)
	}
}
