package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName           string
	Env               string
	HTTPPort          string
	DBURL             string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	JWTAccessSecret   string
	JWTRefreshSecret  string
	JWTAccessTTLMin   int
	JWTRefreshTTLDays int
	EngineBaseURL     string
	InternalAPIKey    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	cfg := &Config{
		AppName:           getEnv("APP_NAME", "baileys-saas-core"),
		Env:               getEnv("APP_ENV", "development"),
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		DBURL:             getEnv("DATABASE_URL", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           redisDB,
		JWTAccessSecret:   getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:  getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessTTLMin:   mustInt(getEnv("JWT_ACCESS_TTL_MIN", "30")),
		JWTRefreshTTLDays: mustInt(getEnv("JWT_REFRESH_TTL_DAYS", "30")),
		EngineBaseURL:     strings.TrimRight(getEnv("ENGINE_BASE_URL", "http://whatsapp-engine:8090"), "/"),
		InternalAPIKey:    getEnv("INTERNAL_API_KEY", ""),
	}

	if cfg.DBURL == "" || cfg.JWTAccessSecret == "" || cfg.JWTRefreshSecret == "" || cfg.InternalAPIKey == "" {
		return nil, fmt.Errorf("missing required envs: DATABASE_URL, JWT_ACCESS_SECRET, JWT_REFRESH_SECRET, INTERNAL_API_KEY")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	return val
}

func mustInt(raw string) int {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return v
}
