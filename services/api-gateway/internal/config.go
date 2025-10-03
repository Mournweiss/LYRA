package internal

import (
	"os"
)

type Config struct {
	GatewayPort         string
	GatewayAddress      string
	WhisperServiceAddr  string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		GatewayPort:        getenvDefault("API_GATEWAY_PORT", "50051"),
		GatewayAddress:     getenvDefault("API_GATEWAY_ADDRESS", "api-gateway:50051"),
		WhisperServiceAddr: getenvDefault("WHISPER_SERVICE_ADDRESS", "whisper-service:50052"),
	}

	if cfg.GatewayPort == "" {
		return nil, ConfigErrorf("API_GATEWAY_PORT environment variable is required")
	}

	if cfg.GatewayAddress == "" {
		return nil, ConfigErrorf("API_GATEWAY_ADDRESS environment variable is required")
	}

	if cfg.WhisperServiceAddr == "" {
		return nil, ConfigErrorf("WHISPER_SERVICE_ADDRESS environment variable is required")
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
