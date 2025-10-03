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
		GatewayPort:        os.Getenv("API_GATEWAY_PORT"),
		GatewayAddress:     os.Getenv("API_GATEWAY_ADDRESS"),
		WhisperServiceAddr: os.Getenv("WHISPER_SERVICE_ADDRESS"),
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
