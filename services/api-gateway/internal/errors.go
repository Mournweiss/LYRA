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

type HandlerError struct {
	GatewayError
	Code string
}

func HandlerErrorf(code, msg string, args ...interface{}) *HandlerError {
	return &HandlerError{
		GatewayError: GatewayError{Msg: fmt.Sprintf(msg, args...)},
		Code:        code,
	}
}

type ValidationError struct {
	GatewayError
	Field string
}

func ValidationErrorf(field, msg string, args ...interface{}) *ValidationError {
	return &ValidationError{
		GatewayError: GatewayError{Msg: fmt.Sprintf(msg, args...)},
		Field:       field,
	}
}
