package server

import (
	"go.uber.org/zap"
	"time"
)

// Config for server
type Config struct {
	// Register to consul
	Register bool
	// ServiceName to register
	ServiceName string
	// Address of local IP
	Address string
	// Port to listen
	Port int
	// Logger for trace
	Logger *zap.SugaredLogger
	// LatencyLimit for trace error
	LatencyLimit time.Duration
	// LogFilter func
	LogFilter func(m []interface{})
}
