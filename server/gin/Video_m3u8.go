package gin

import (
	"fmt"
	"github.com/skyandong/service-go/service/core/downloadM3u8"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyandong/tool/logger"
	"github.com/skyandong/tool/service"

	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/service/core"
)

var lgM3u8 = conf.C.Loggers.Get(conf.LoggerName).GetLogger(logger.InfoLevel)

func getVideoFromM3u8(c *gin.Context) {
	var req api.GetVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	tid := service.GetTraceID(c)
	ctx := &core.Context{
		Logger:         lgM3u8,
		TraceID:        tid,
		URL:            req.URL,
		FileName:       req.FileName,
		DepositAddress: req.DepositAddress,
	}
	fmt.Println("url", req.URL, "DepositAddress", req.DepositAddress, "name", req.FileName)

	worker, err := downloadM3u8.NewTask(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if worker != nil {
		go worker.Start(req.ChanNum)
	}
	c.JSON(http.StatusOK, "ok")
}
