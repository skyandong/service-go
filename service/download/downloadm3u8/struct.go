package downloadm3u8

import (
	"runtime"
	"sync"

	"go.uber.org/zap"

	"github.com/skyandong/service-go/service/parse"
)

var wg sync.WaitGroup

var gnums = runtime.NumCPU()

const (
	maxTryAgain = 3
)

// Downloader 下载器
type Downloader struct {
	*parse.Result
	TraceID         string
	lock            sync.Mutex
	Folder          string //下载文件目录
	TsFolder        string //临时文件目录
	MergeTSFilename string //下载文件名
	Logger          *zap.SugaredLogger
}
