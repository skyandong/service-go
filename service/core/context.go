package core

import (
	"context"

	"go.uber.org/zap"
)

// Context for request
type Context struct {
	// Ctx of go
	Ctx context.Context
	// TraceID of request
	TraceID string
	//Url Resource address
	URL string
	//DepositAddress Download directory
	DepositAddress string
	//FileName Download file name
	FileName string
	// Logger obj
	Logger *zap.SugaredLogger
}
