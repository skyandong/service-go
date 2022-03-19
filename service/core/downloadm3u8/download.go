package downloadm3u8

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/skyandong/service-go/service/tool"
)

const (
	maxTryAgain = 3
)

func (d *Downloader) NewStart(concurrency int) {
	if concurrency > len(d.result.M3u8.Segments) {
		concurrency = len(d.result.M3u8.Segments)
	}
	ch := make(chan *Downloader, concurrency)
	for {
		select {
		case ch <- d:
			go w()
		}
	}
}

func w() {

}

// Start runs downloader
func (d *Downloader) Start(concurrency int) {
	var wg sync.WaitGroup

	defer func() {
		err := os.RemoveAll(d.tsFolder)
		if err != nil {
			d.logger.Errorw("remove ts file or ts folder error", "file_name", d.tsFolder, "err", err, "traceID", d.traceID)
		}
	}()

	// 是否停止下载
	var isStop bool

	// struct{} zero size
	limitChan := make(chan struct{}, concurrency)
	for {
		// 获取下一个节点
		if len(d.queue) > 0 {
			tsIdx := d.queue[0]
			d.queue = d.queue[1:]

			// 开启协程
			wg.Add(1)
			go d.chanWork(tsIdx, limitChan, &wg, &isStop)
			limitChan <- struct{}{}

			// 如果一个文件下载失败，则判定整个视频失败
			if isStop {
				break
			}
		}
	}
	wg.Wait()

	if !isStop {
		d.logger.Infow("start merge .ts file", "file_name", d.tsFolder, "traceID", d.traceID)
		if err := d.merge(len(d.result.M3u8.Segments)); err != nil {
			d.logger.Errorw("merge file error", "err", err, "name", d.mergeTSFilename, "traceId", d.traceID)
		}
	} else {
		// 清理 ts 文件,避免污染空间
		os.RemoveAll(d.tsFolder)
	}
	return
}

func (d *Downloader) chanWork(idx int, limitChan chan struct{}, wg *sync.WaitGroup, stop *bool) {
	defer wg.Done()
	if err := d.download(idx); err != nil {
		d.logger.Errorw("download failed", "idx", idx, "err", err, "traceId", d.traceID)
		// Back into the queue, retry request
		var tryAgain = 1
		for ; tryAgain <= maxTryAgain; tryAgain++ {
			if err = d.download(idx); err != nil {
				d.logger.Errorw("try again error", "err", err, "idx", idx, "traceId", d.traceID, "try_times", tryAgain)
			}
		}
		if tryAgain > maxTryAgain {
			*stop = true
		}
	} else {
		d.logger.Infow("download idx ok", "idx", idx, "traceId", d.traceID)
	}
}

func (d *Downloader) download(segIndex int) error {
	//Fragment file name
	tsFilename := tsFilename(segIndex)

	//Fragment address
	tsURL := d.tsURL(segIndex)

	b, e := tool.Get(tsURL)
	if e != nil {
		return fmt.Errorf("request %s, %s", tsURL, e.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer b.Close()

	//创建临时文件
	fPath := filepath.Join(d.tsFolder, tsFilename)
	fTemp := fPath + tsTempFileSuffix
	f, err := os.Create(fTemp)
	if err != nil {
		return fmt.Errorf("create file: %s, %s", tsFilename, err.Error())
	}

	//读取全部字节，如果字节被加密，选择合适解密方式(如果有)
	bytes, err := ioutil.ReadAll(b)
	if err != nil {
		return fmt.Errorf("read bytes: %s, %s", tsURL, err.Error())
	}
	sf := d.result.M3u8.Segments[segIndex]
	if sf == nil {
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	key, ok := d.result.Keys[sf.KeyIndex]
	if ok && key != "" {
		bytes, err = tool.AES128Decrypt(bytes, []byte(key),
			[]byte(d.result.M3u8.Keys[sf.KeyIndex].IV))
		if err != nil {
			return fmt.Errorf("decryt: %s, %s", tsURL, err.Error())
		}
	}
	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) //0x47
	bLen := len(bytes)
	for j := 0; j < bLen; j++ {
		if bytes[j] == syncByte {
			bytes = bytes[j:]
			break
		}
	}
	w := bufio.NewWriter(f)
	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("write to %s: %s", fTemp, err.Error())
	}
	// Release file resource to rename file
	_ = f.Close()
	if err = os.Rename(fTemp, fPath); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return tool.ResolveURL(d.result.URL, seg.URI)
}

func tsFilename(ts int) string {
	return strconv.Itoa(ts) + tsExt
}

func genSlice(len int) []int {
	s := make([]int, 0)
	for i := 0; i < len; i++ {
		s = append(s, i)
	}
	return s
}
