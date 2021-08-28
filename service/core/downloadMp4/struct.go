package downloadMp4

import (
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
