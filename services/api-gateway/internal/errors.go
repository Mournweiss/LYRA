package internal

import "fmt"

type GatewayError struct {
	Msg string
}

func (e *GatewayError) Error() string {
	return e.Msg
}

type ConfigError struct {
	GatewayError
}

func ConfigErrorf(msg string, args ...interface{}) *ConfigError {
	return &ConfigError{GatewayError{Msg: fmt.Sprintf(msg, args...)}}
}

type UpstreamError struct {
	GatewayError
}

func UpstreamErrorf(msg string, args ...interface{}) *UpstreamError {
	return &UpstreamError{GatewayError{Msg: fmt.Sprintf(msg, args...)}}
}
