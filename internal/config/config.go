package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	WebhookURL string
	AuthKey    string
	Interval   time.Duration
	BatchSize  int
	Port       string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "messages_db"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		WebhookURL: getEnv("WEBHOOK_URL", "http://localhost:8080/webhook"),
		AuthKey:    getEnv("AUTH_KEY", "secretkey"),
		Interval:   time.Duration(getEnvAsInt("INTERVAL_SEC", 60)) * time.Second,
		BatchSize:  getEnvAsInt("BATCH_SIZE", 2),
		Port:       getEnv("PORT", "8000"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
