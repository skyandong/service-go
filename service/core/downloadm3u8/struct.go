package downloadm3u8

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/skyandong/service-go/service/parse"
)

// Downloader 下载器
type Downloader struct {
	lock            sync.Mutex
	queue           []int
	folder          string //下载文件目录
	tsFolder        string //临时文件目录
	result          *parse.Result
	mergeTSFilename string //下载文件名
	traceID         string
	logger          *zap.SugaredLogger
}

// Context for request
type Context struct {
	Ctx             context.Context    // Ctx of go
	TraceID         string             // TraceID of request
	URL             string             //Url Resource address
	DownloadCatalog string             //DownloadCatalog Download directory
	FileName        string             //FileName Download file name
	ChanNum         int                // ChanNum
	Logger          *zap.SugaredLogger // Logger obj
}
