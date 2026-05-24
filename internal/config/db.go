package config

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DBConfig struct {
	DB_HOST               string
	DB_PORT               string
	DB_USER               string
	DB_PASSWORD           string
	DB_NAME               string
	DATABASE_URL          string
	DB_MAX_IDLE_CONNS     int
	DB_MAX_OPEN_CONNS     int
	DB_CONN_MAX_LIFETIME  time.Duration
	DB_CONN_MAX_IDLE_TIME time.Duration
	DB_SSL_MODE           string
}

func DefaultDBConfig() DBConfig {
	return DBConfig{
		DB_MAX_IDLE_CONNS:     5,
		DB_MAX_OPEN_CONNS:     20,
		DB_CONN_MAX_LIFETIME:  60 * time.Minute,
		DB_CONN_MAX_IDLE_TIME: 10 * time.Minute,
		DB_SSL_MODE:           "disable", // Use "require" or "verify-full" in production
	}
}

func Load() (DBConfig, error) {
	config := DefaultDBConfig()
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	config.DB_PORT = getEnv("PORT", "8080")
	config.DATABASE_URL = getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/hris_db?sslmode=disable")

	if maxIdle := getEnv("DB_MAX_IDLE_CONNS", ""); maxIdle != "" {
		if val, err := strconv.Atoi(maxIdle); err == nil {
			config.DB_MAX_IDLE_CONNS = val
		}
	}

	if maxOpen := getEnv("DB_MAX_OPEN_CONNS", ""); maxOpen != "" {
		if val, err := strconv.Atoi(maxOpen); err == nil {
			config.DB_MAX_OPEN_CONNS = val
		}
	}

	if sslMode := getEnv("DB_SSL_MODE", ""); sslMode != "" {
		config.DB_SSL_MODE = sslMode
	}

	return config, nil
}

func NewPostgresDB() (*sql.DB, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", config.DATABASE_URL)
	if err != nil {
		return nil, err
	}

	// Set connection pool configuration
	db.SetMaxIdleConns(config.DB_MAX_IDLE_CONNS)
	db.SetMaxOpenConns(config.DB_MAX_OPEN_CONNS)
	db.SetConnMaxLifetime(config.DB_CONN_MAX_LIFETIME)
	db.SetConnMaxIdleTime(config.DB_CONN_MAX_IDLE_TIME)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
