package downloadMp4

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skyandong/service-go/service/core"
	"github.com/skyandong/service-go/service/tool"
)

// NewTask 初始化
func NewTask(c *core.Context) (*Downloader, error) {
	//获取url信息
	worker, err := getFileInfoFromURL(c.URL)
	if err != nil {
		c.Logger.Errorw("making request error", "error", err)
		return nil, err
	}

	//生成目录文件
	var folder string
	if c.DepositAddress == "" {
		current, err := tool.CurrentDir()
		if err != nil {
			return nil, err
		}
		folder = filepath.Join(current, c.DepositAddress)
	} else {
		folder = c.DepositAddress
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create storage folder failed: %s", err.Error())
	}

	fileName := filepath.Join(folder, c.FileName)
	worker.file, err = os.Create(fileName)
	if err != nil {
		c.Logger.Errorw("Error creating file", "error", err)
		return nil, err
	}
	defer worker.file.Close()
	return worker, nil
}

func getFileInfoFromURL(url string) (*Downloader, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	acceptRangesResult := resp.Header.Get("Accept-Ranges")
	contentLengthResult := resp.Header.Get("Content-Length")
	contentLengthResultInt, err := strconv.Atoi(contentLengthResult)
	md5 := resp.Header.Get("ETag")

	md5 = strings.Replace(md5, "\"", "", -1)

	if err != nil {
		return &Downloader{url: url, size: 0, rangeSupport: false, md5: ""}, err
	}
	return &Downloader{url: url, size: int64(contentLengthResultInt), rangeSupport: "bytes" == acceptRangesResult, md5: md5}, err
}
