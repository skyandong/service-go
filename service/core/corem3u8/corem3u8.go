package corem3u8

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skyandong/service-go/service/core"
	"github.com/skyandong/service-go/service/download/downloadm3u8"
	"github.com/skyandong/service-go/service/parse"
	"github.com/skyandong/service-go/service/tool"
)

const (
	tsExt            = ".ts"
	tsFolderName     = "ts"
	tsTempFileSuffix = "_tmp"
)

// Work returns a Task instance
func Work(c *core.Context) error {
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

	 _ = &downloadm3u8.Downloader{
		Folder:          folder,
		TsFolder:        tsFolder,
		MergeTSFilename: c.FileName,
		Result:          result,
		TraceID:         c.TraceID,
		Logger:          c.Logger,
	}
	//d.queue = genSlice(len(result.M3u8.Segments))
	c.Logger.Infow("worker is alredy", "folder", folder, "file_name", c.FileName, "url", result.URL.String(), "traceId", c.TraceID)

	//go d.Start()
	return nil
}
