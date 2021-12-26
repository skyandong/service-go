package downloadM3u8

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skyandong/service-go/service/core"
	"github.com/skyandong/service-go/service/parse"
	"github.com/skyandong/service-go/service/tool"
)

const (
	tsExt            = ".ts"
	tsFolderName     = "ts"
	tsTempFileSuffix = "_tmp"
)

// NewTask returns a Task instance
func NewTask(c *core.Context) (*Downloader, error) {

	//----解析url,获取片段----
	result, err := parse.FromURL(c.URL)
	if err != nil {
		return nil, err
	}

	//----文件目录----
	var folder string
	// If no output folder specified, use current directory
	if c.DownloadCatalog == "" {
		current, err := tool.CurrentDir()
		if err != nil {
			return nil, err
		}
		folder = filepath.Join(current, c.DownloadCatalog)
	} else {
		folder = c.DownloadCatalog
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create storage folder failed: %s", err.Error())
	}

	//----临时文件目录----
	tsFolder := filepath.Join(folder, tsFolderName, c.FileName)
	if err := os.MkdirAll(tsFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create ts folder '[%s]' failed: %s", tsFolder, err.Error())
	}

	d := &Downloader{
		folder:          folder,
		tsFolder:        tsFolder,
		mergeTSFilename: c.FileName,
		result:          result,
		traceID:         c.TraceID,
		logger:          c.Logger,
	}
	d.queue = genSlice(len(result.M3u8.Segments))
	d.logger.Infow("worker is alredy", "folder", d.folder, "file_name", d.mergeTSFilename, "url", d.result.URL.String(), "traceId", d.traceID)
	return d, nil
}
