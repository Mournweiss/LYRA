package internal

import (
	"os"
	"fmt"
	"strconv"
)

type Config struct {
	GatewayPort         string
	GatewayHost         string
	GatewayAddress      string
	WhisperServiceHost  string
	WhisperServicePort  string
	WhisperServiceAddr  string
	RedisHost           string
	RedisPort           string
	WorkerConcurrency   int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		GatewayPort:        getenvDefault("API_GATEWAY_PORT", "50051"),
		GatewayHost:        getenvDefault("API_GATEWAY_HOST", "api-gateway"),
		WhisperServicePort: getenvDefault("WHISPER_SERVICE_PORT", "50052"),
		WhisperServiceHost: getenvDefault("WHISPER_SERVICE_HOST", "whisper-service"),
		RedisHost:          getenvDefault("REDIS_HOST", "redis"),
		RedisPort:          getenvDefault("REDIS_PORT", "6379"),
	}

	concurrencyStr := getenvDefault("WORKER_CONCURRENCY", "4")
	concurrency, err := strconv.Atoi(concurrencyStr)

	if err != nil || concurrency < 1 {
		concurrency = 4
	}

	cfg.WorkerConcurrency = concurrency
	cfg.GatewayAddress = fmt.Sprintf("%s:%s", cfg.GatewayHost, cfg.GatewayPort)
	cfg.WhisperServiceAddr = fmt.Sprintf("%s:%s", cfg.WhisperServiceHost, cfg.WhisperServicePort)

	if cfg.GatewayPort == "" {
		return nil, ConfigErrorf("API_GATEWAY_PORT environment variable is required")
	}

	if cfg.GatewayHost == "" {
		return nil, ConfigErrorf("API_GATEWAY_HOST environment variable is required")
	}
	if cfg.WhisperServicePort == "" {
		return nil, ConfigErrorf("WHISPER_SERVICE_PORT environment variable is required")
	}

	if cfg.WhisperServiceHost == "" {
		return nil, ConfigErrorf("WHISPER_SERVICE_HOST environment variable is required")
	}

	if cfg.RedisHost == "" {
		return nil, ConfigErrorf("REDIS_HOST environment variable is required")
	}

	if cfg.RedisPort == "" {
		return nil, ConfigErrorf("REDIS_PORT environment variable is required")
	}

	return cfg, nil
}

func getenvDefault(key, def string) string {
	val := os.Getenv(key)
	
	if val == "" {
		return def
	}

	return val
}
