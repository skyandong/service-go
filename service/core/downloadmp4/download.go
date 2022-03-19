package downloadmp4

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

// Coroutines 协程数
var Coroutines = 1

// Start 启动下载
func (d *Downloader) Start() {
	//download file
	if d.rangeSupport == true {
		d.lg.Info("Byte Ranges Supported, Downloading Please Wait...")
		rangeSize, finalRangeSize := calculateRangeSizeAndLastRange(d.size)
		for i := 0; i < Coroutines; i++ {
			wg.Add(1)
			start := int64(i * rangeSize)
			finish := int64(int(start) + rangeSize)

			//Is this the last chunk?
			if i+1 == Coroutines {
				finish = start + int64(finalRangeSize)
			}
			go d.downloadURLToFileAtRange(start, finish)
		}
		wg.Wait()
	} else {
		d.lg.Info("Byte Ranges Unsupported, Downloading Please Wait...")
		d.downloadURLToFileAtRange(0, d.size)
	}

	//Validate Integrity of file
	md5Result, err := d.validateFileWithMD5()
	if err != nil {
		d.lg.Errorw("Error Calculating md5", "error", err)
	} else {
		d.lg.Infow("File Has Integrity", "md5", md5Result)
	}
}

// downloadURLToFileAtRange 获取数据
func (d *Downloader) downloadURLToFileAtRange(start int64, end int64) {
	defer wg.Done()
	client := &http.Client{}
	bytesString := fmt.Sprintf("bytes=%d-%d", start, end)
	req, _ := http.NewRequest("GET", d.url, nil)
	if d.rangeSupport {
		req.Header.Set("Range", bytesString)
	}
	resp, err := client.Do(req)
	if err != nil {
		d.lg.Infow("request video error", "error", err)
		return
	}
	size, err := d.readFromReaderIntoFileAtOffset(resp.Body, start)
	if err == nil {
		d.lg.Infow("download block ok", "size", size)
	}
}

// validateFileWithMD5 验证数据
func (d *Downloader) validateFileWithMD5() (bool, error) {
	var result []byte
	hash := md5.New()
	if _, err := io.Copy(hash, d.file); err != nil {
		return false, err
	}
	byteHash := hash.Sum(result)
	computedMD5 := fmt.Sprintf("%x", byteHash)

	return computedMD5 == d.md5, nil
}

func calculateRangeSizeAndLastRange(fileSize int64) (mainChunkSize int, lastChunchSize int) {
	mainChunkSize = int(fileSize) / Coroutines
	lastChunchSize = mainChunkSize + int(fileSize)%Coroutines
	return mainChunkSize, lastChunchSize
}
