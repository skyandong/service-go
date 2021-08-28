package api

// GetVideoRequest 下载视频请求参数
type GetVideoRequest struct {
	URL            string `json:"url"`
	DepositAddress string `json:"deposit_address"`
	FileName       string `json:"file_name"`
	ChanNum        int    `json:"chan_num"`
}
