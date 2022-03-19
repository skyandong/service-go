package downloadm3u8

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (d *Downloader) merge(segLen int) error {
	// In fact, the number of downloaded segments should be equal to number of m3u8 segments
	missingCount := 0
	for idx := 0; idx < segLen; idx++ {
		tsFilename := tsFilename(idx)
		f := filepath.Join(d.tsFolder, tsFilename)
		if _, err := os.Stat(f); err != nil {
			d.logger.Warnw("file miss", "idx", idx, "file_name", tsFilename, "traceId", d.traceID)
			missingCount++
		}
	}
	if missingCount > 0 {
		d.logger.Warnw("files missing count", "count", missingCount, "traceId", d.traceID)
	}

	// Create a TS file for merging, all segment files will be written to this file.
	mFilePath := filepath.Join(d.folder, d.mergeTSFilename)
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("create main TS file failed：%s", err.Error())
	}

	//noinspection GoUnhandledErrorResult
	defer mFile.Close()

	writer := bufio.NewWriter(mFile)
	mergedCount := 0
	for segIndex := 0; segIndex < segLen; segIndex++ {
		tsFilename := tsFilename(segIndex)
		bytes, err := ioutil.ReadFile(filepath.Join(d.tsFolder, tsFilename))
		_, err = writer.Write(bytes)
		if err != nil {
			d.logger.Warnw("files merge failed", "ts_file", tsFilename, "mfile", mFile)
			continue
		}
		mergedCount++
	}

	_ = writer.Flush()
	// Remove `ts` folder
	_ = os.RemoveAll(d.tsFolder)

	if mergedCount != segLen {
		d.logger.Warnw("files merge failed", "mergedCount", mergedCount, "seg_len", segLen)
	}

	d.logger.Info("merger file over")

	return nil
}
