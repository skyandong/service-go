package gin

import (
	"github.com/skyandong/service-go/service/core/downloadMp4"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyandong/tool/logger"
	"github.com/skyandong/tool/service"

	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/service/core"
)

var lgMp4 = conf.C.Loggers.Get(conf.LoggerName).GetLogger(logger.InfoLevel)

func getVideoFromMp4(c *gin.Context) {
	req := api.GetVideoRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	tid := service.GetTraceID(c)
	ctx := &core.Context{
		Ctx:            c,
		TraceID:        tid,
		Url:            req.Url,
		DepositAddress: req.DepositAddress,
		FileName:       req.FileName,
		Logger:         lgMp4,
	}
	worker, err := downloadMp4.NewTask(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if worker != nil {
		go worker.Start()
	}
}
