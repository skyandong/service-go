package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/service/core/downloadm3u8"
	"github.com/skyandong/tool/logger"
	"github.com/skyandong/tool/service"
)

var lgM3u8 = conf.C.Loggers.Get(conf.LoggerName).GetLogger(logger.InfoLevel)

func getVideoFromM3u8(c *gin.Context) {
	var req api.GetVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ctx := &downloadm3u8.Context{
		Logger:          lgM3u8,
		TraceID:         service.GetTraceID(c),
		URL:             req.URL,
		FileName:        req.FileName,
		ChanNum:         req.ChanNum,
		DownloadCatalog: req.DepositAddress,
	}
	if err := downloadm3u8.Work(ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "ok")
}
