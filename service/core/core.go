package core

import (
	"context"
	"go.uber.org/zap"
)

// Context for request
type Context struct {
	Ctx             context.Context    // Ctx of go
	TraceID         string             // TraceID of request
	URL             string             //URL Resource address
	DownloadCatalog string             //DownloadCatalog Download directory
	FileName        string             //FileName Download file name
	ChanNum         int                // ChanNum
	Logger          *zap.SugaredLogger // Logger obj
}
