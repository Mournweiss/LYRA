package internal

import (
	"os"
	"fmt"
)

type Config struct {
	GatewayPort         string
	GatewayDomain       string
	GatewayAddress      string
	WhisperServiceDomain string
	WhisperServicePort   string
	WhisperServiceAddr   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		GatewayPort:        getenvDefault("API_GATEWAY_PORT", "50051"),
		GatewayDomain:      getenvDefault("API_GATEWAY_DOMAIN", "api-gateway"),
		WhisperServicePort: getenvDefault("WHISPER_SERVICE_PORT", "50052"),
		WhisperServiceDomain: getenvDefault("WHISPER_SERVICE_DOMAIN", "whisper-service"),
	}

	cfg.GatewayAddress = fmt.Sprintf("%s:%s", cfg.GatewayDomain, cfg.GatewayPort)
	cfg.WhisperServiceAddr = fmt.Sprintf("%s:%s", cfg.WhisperServiceDomain, cfg.WhisperServicePort)

	if cfg.GatewayPort == "" {
		return nil, ConfigErrorf("API_GATEWAY_PORT environment variable is required")
	}
	if cfg.GatewayDomain == "" {
		return nil, ConfigErrorf("API_GATEWAY_DOMAIN environment variable is required")
	}
	if cfg.WhisperServicePort == "" {
		return nil, ConfigErrorf("WHISPER_SERVICE_PORT environment variable is required")
	}
	if cfg.WhisperServiceDomain == "" {
		return nil, ConfigErrorf("WHISPER_SERVICE_DOMAIN environment variable is required")
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
