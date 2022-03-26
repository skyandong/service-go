package downloadm3u8

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skyandong/service-go/service/parse"
	"github.com/skyandong/service-go/service/tool"
)

const (
	tsExt            = ".ts"
	tsFolderName     = "ts"
	tsTempFileSuffix = "_tmp"
)

// Work returns a Task instance
func Work(c *Context) error {
	result, err := parse.FromURL(c.URL)
	if err != nil {
		return err
	}

	var folder string
	// If no output folder specified, use current directory
	if c.DownloadCatalog == "" {
		current, err := tool.CurrentDir()
		if err != nil {
			return err
		}
		folder = filepath.Join(current, c.DownloadCatalog)
	} else {
		folder = c.DownloadCatalog
	}
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("create storage folder failed: %s", err.Error())
	}

	tsFolder := filepath.Join(folder, tsFolderName, c.FileName)
	if err := os.MkdirAll(tsFolder, os.ModePerm); err != nil {
		return fmt.Errorf("create ts folder '[%s]' failed: %s", tsFolder, err.Error())
	}

	d := &Downloader{
		folder:          folder,
		tsFolder:        tsFolder,
		mergeTSFilename: c.FileName,
		Result:          result,
		traceID:         c.TraceID,
		logger:          c.Logger,
	}
	//d.queue = genSlice(len(result.M3u8.Segments))
	d.logger.Infow("worker is alredy", "folder", d.folder, "file_name", d.mergeTSFilename, "url", d.result.URL.String(), "traceId", d.traceID)

	go d.Start()
	return nil
}
