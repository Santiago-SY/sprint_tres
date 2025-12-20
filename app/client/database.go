package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB inicia el pool de conexiones a PostgreSQL
func InitDB() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	// URL de conexión
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPass, dbHost, dbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error config DB: %v", err)
	}

	// Configuración de seguridad para no saturar
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour

	// Reintentos de conexión (por si la DB tarda en arrancar)
	var pool *pgxpool.Pool
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			if errPing := pool.Ping(context.Background()); errPing == nil {
				fmt.Println("🐘 Conexión a PostgreSQL EXITOSA")
				return pool, nil
			}
		}
		fmt.Printf("⏳ Esperando a Postgres... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("no se pudo conectar a Postgres: %v", err)
}
