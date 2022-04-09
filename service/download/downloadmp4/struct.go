package downloadmp4

import (
	"context"
	"os"

	"go.uber.org/zap"
)

// Downloader 下载器
type Downloader struct {
	url          string
	size         int64
	rangeSupport bool
	md5          string
	traceID      string
	file         *os.File
	lg           *zap.SugaredLogger
}

// Context for request
type Context struct {
	// Ctx of go
	Ctx context.Context
	// TraceID of request
	TraceID string
	//Url Resource address
	URL string
	//DownloadCatalog Download directory
	DownloadCatalog string
	//FileName Download file name
	FileName string
	// Logger obj
	Logger *zap.SugaredLogger
}
