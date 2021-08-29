package downloadM3u8

import (
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
