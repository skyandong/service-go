package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// regex pattern for extracting `key=value` parameters from a line
var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

const (
	playlistTypeVOD   playlistType = "VOD"
	playlistTypeEvent playlistType = "EVENT"

	cryptMethodAES  cryptMethod = "AES-128"
	cryptMethodNONE cryptMethod = "NONE"
)

type (
	playlistType string
	cryptMethod  string
)

// Segment ...
type Segment struct {
	URI      string
	KeyIndex int
	Title    string  // #EXTINF: duration,<title>
	Duration float32 // #EXTINF: duration,<title>
	Length   uint64  // #EXT-X-BYTERANGE: length[@offset]
	Offset   uint64  // #EXT-X-BYTERANGE: length[@offset]
}

// MasterPlaylist #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=240000,RESOLUTION=416x234,CODECS="avc1.42e00a,mp4a.40.2"
type MasterPlaylist struct {
	URI        string
	BandWidth  uint32
	Resolution string
	Codecs     string
	ProgramID  uint32
}

// Key #EXT-X-KEY:METHOD=AES-128,URI="key.key"
type Key struct {
	// 'AES-128' or 'NONE'
	// If the encryption method is NONE, the URI and the IV attributes MUST NOT be present
	Method cryptMethod
	URI    string
	IV     string
}

// M3u8 ...
type M3u8 struct {
	Version        int8   // EXT-X-VERSION:version
	MediaSequence  uint64 // Default 0, #EXT-X-MEDIA-SEQUENCE:sequence
	Segments       []*Segment
	MasterPlaylist []*MasterPlaylist
	Keys           map[int]*Key
	EndList        bool         // #EXT-X-ENDLIST
	PlaylistType   playlistType // VOD or EVENT
	TargetDuration float64      // #EXT-X-TARGETDURATION:duration
}

func parse(reader io.Reader) (*M3u8, error) {
	s := bufio.NewScanner(reader)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	m3u8 := &M3u8{
		Keys: make(map[int]*Key),
	}

	var (
		keyIndex = 0
		key      *Key
		seg      *Segment
		extInf   bool
		extByte  bool
	)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// 检查 m3u8 文件头行
		if i == 0 && "#EXTM3U" != line {
			return nil, fmt.Errorf("invalid m3u8, missing #EXTM3U in line 1")
		}

		switch {
		case line == "":
			continue
		//EXT-X-PLAYLIST-TYPE 表明流媒体类型,全局生效
		case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &m3u8.PlaylistType); err != nil {
				return nil, err
			}
			isValid := m3u8.PlaylistType == "" || m3u8.PlaylistType == playlistTypeVOD || m3u8.PlaylistType == playlistTypeEvent
			if !isValid {
				return nil, fmt.Errorf("invalid playlist type: %s, line: %d", m3u8.PlaylistType, i+1)
			}

		// EXT-X-TARGETDURATION 表示每个视频分段最大的时长（单位秒）。
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &m3u8.TargetDuration); err != nil {
				return nil, err
			}

		// EXT-X-VERSION 指示出 playlist 文件的兼容版本
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-VERSION:%d", &m3u8.Version); err != nil {
				return nil, err
			}

		// EXT-X-MEDIA-SEQUENCE 表示播放列表第一个 URL 片段文件的序列号
		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &m3u8.MediaSequence); err != nil {
				return nil, err
			}

		// Parse master playlist，#EXT-X-STREAM-INF：指定一个包含多媒体信息的 media URI 作为PlayList，一般做M3U8的嵌套使用，它只对紧跟后面的URI有
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			mp, err := parseMasterPlaylist(line)
			if err != nil {
				return nil, err
			}
			i++
			mp.URI = lines[i]
			if mp.URI == "" || strings.HasPrefix(mp.URI, "#") {
				return nil, fmt.Errorf("invalid EXT-X-STREAM-INF URI, line: %d", i+1)
			}
			m3u8.MasterPlaylist = append(m3u8.MasterPlaylist, mp)
			continue

		// EXTINF 表示其后 URL 指定的媒体片段时长（单位为秒),每个 URL 媒体片段之前必须指定该标签
		case strings.HasPrefix(line, "#EXTINF:"):
			if extInf {
				return nil, fmt.Errorf("duplicate EXTINF: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var s string
			if _, err := fmt.Sscanf(line, "#EXTINF:%s", &s); err != nil {
				return nil, err
			}
			if strings.Contains(s, ",") {
				split := strings.Split(s, ",")
				seg.Title = split[1]
				s = split[0]
			}
			df, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, err
			}
			seg.Duration = float32(df)
			seg.KeyIndex = keyIndex
			extInf = true

		//#EXT-X-BYTERANGE表示媒体段是一个媒体URI资源中的一段，只对其后的media URI有效
		case strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
			if extByte {
				return nil, fmt.Errorf("duplicate EXT-X-BYTERANGE: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var b string
			if _, err := fmt.Sscanf(line, "#EXT-X-BYTERANGE:%s", &b); err != nil {
				return nil, err
			}
			if b == "" {
				return nil, fmt.Errorf("invalid EXT-X-BYTERANGE, line: %d", i+1)
			}
			if strings.Contains(b, "@") {
				split := strings.Split(b, "@")
				offset, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					return nil, err
				}
				seg.Offset = uint64(offset)
				b = split[0]
			}
			length, err := strconv.ParseUint(b, 10, 64)
			if err != nil {
				return nil, err
			}
			seg.Length = uint64(length)
			extByte = true

		// Parse segments URI
		case !strings.HasPrefix(line, "#"):
			if extInf {
				if seg == nil {
					return nil, fmt.Errorf("invalid line: %s", line)
				}
				seg.URI = line
				extByte = false
				extInf = false
				m3u8.Segments = append(m3u8.Segments, seg)
				seg = nil
				continue
			}

		// Parse key
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			params := parseLineParameters(line)
			if len(params) == 0 {
				return nil, fmt.Errorf("invalid EXT-X-KEY: %s, line: %d", line, i+1)
			}
			method := cryptMethod(params["METHOD"])
			if method != "" && method != cryptMethodAES && method != cryptMethodNONE {
				return nil, fmt.Errorf("invalid EXT-X-KEY method: %s, line: %d", method, i+1)
			}
			keyIndex++
			key = new(Key)
			key.Method = method
			key.URI = params["URI"]
			key.IV = params["IV"]
			m3u8.Keys[keyIndex] = key
		case line == "#EXT-X-ENDLIST":
			m3u8.EndList = true
		default:
			continue
		}
	}
	return m3u8, nil
}

func parseMasterPlaylist(line string) (*MasterPlaylist, error) {
	params := parseLineParameters(line)
	if len(params) == 0 {
		return nil, errors.New("empty parameter")
	}
	mp := new(MasterPlaylist)
	for k, v := range params {
		switch {
		case k == "BANDWIDTH":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.BandWidth = uint32(v)
		case k == "RESOLUTION":
			mp.Resolution = v
		case k == "PROGRAM-ID":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.ProgramID = uint32(v)
		case k == "CODECS":
			mp.Codecs = v
		}
	}
	return mp, nil
}

// parseLineParameters extra parameters in string `line`
func parseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}
