package internal

import (
	"fmt"
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
		return nil, fmt.Errorf("API_GATEWAY_PORT environment variable is required")
	}

	if cfg.GatewayAddress == "" {
		return nil, fmt.Errorf("API_GATEWAY_ADDRESS environment variable is required")
	}

	if cfg.WhisperServiceAddr == "" {
		return nil, fmt.Errorf("WHISPER_SERVICE_ADDRESS environment variable is required")
	}

	return cfg, nil
}
