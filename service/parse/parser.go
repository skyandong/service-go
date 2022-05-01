package parse

import (
	"errors"
	"fmt"
	"github.com/skyandong/service-go/util"
	"io/ioutil"
	"net/url"
)

// Result ...
type Result struct {
	URL  *url.URL       // 资源地址
	M3u8 *M3u8          // .m3u8 文件结构
	Keys map[int]string // 加密类型
}

// FromURL ...
func FromURL(link string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()

	//get .m3u8 file
	body, err := util.Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer body.Close()

	//Analysis .m3u8 file
	m3u8, err := parse(body)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(util.ResolveURL(u, sf.URI))
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}

	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == cryptMethodNONE:
			continue
		case key.Method == cryptMethodAES:
			// Request URL to extract decryption key
			keyURL := key.URI
			keyURL = util.ResolveURL(u, keyURL)
			resp, err := util.Get(keyURL)
			if err != nil {
				return nil, fmt.Errorf("extract key failed: %s", err.Error())
			}
			keyByte, err := ioutil.ReadAll(resp)
			_ = resp.Close()
			if err != nil {
				return nil, err
			}
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}
