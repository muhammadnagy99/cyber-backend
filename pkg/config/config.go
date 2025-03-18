package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort   string
	DatabasePath string
	BearerToken  string
}

func LoadConfig() *Config {
	config := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8000"),
		DatabasePath: getEnv("DATABASE_PATH", "rbac.db"),
		BearerToken:  getEnv("BEARER_TOKEN", "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6IjY0YjZlYmEwM2RlZWE2ZTVjMjZjMTg1NDQ3ZmE4MDNjIn0.eyJzdWIiOiIyOTcyMTUxOTE5IiwibmFtZSI6IkJJR0JPU1MiLCJpYXQiOjEzMjEyMzEzMjEzMjF9.zN7mG-0pI2EBE2wsXu9jsdfud4uiqBiZDPgxrE0e2mJ4sD_CdesyQPANeEYp6c7log4haM8XbeMVr7P54oO-bQ"),
	}

	log.Printf("Loaded config: %+v\n", config)
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
