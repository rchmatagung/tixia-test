package config

import (
    "os"
)

type Config struct {
    RedisURL string
}

func Load() *Config {
    return &Config{
        RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}