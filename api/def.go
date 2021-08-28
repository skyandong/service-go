package api

// GetVideoRequest 下载视频请求参数
type GetVideoRequest struct {
	URL            string `json:"url"`             // 视频地址
	DepositAddress string `json:"deposit_address"` // 存放地址
	FileName       string `json:"file_name"`       // 视频名称
	ChanNum        int    `json:"chan_num"`        // 协程数
}
