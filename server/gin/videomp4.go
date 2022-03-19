package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/service/core/downloadmp4"
	"github.com/skyandong/tool/logger"
	"github.com/skyandong/tool/service"
)

var lgMp4 = conf.C.Loggers.Get(conf.LoggerName).GetLogger(logger.InfoLevel)

func getVideoFromMp4(c *gin.Context) {
	req := api.GetVideoRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := &downloadmp4.Context{
		Ctx:             c,
		TraceID:         service.GetTraceID(c),
		URL:             req.URL,
		DownloadCatalog: req.DepositAddress,
		FileName:        req.FileName,
		Logger:          lgMp4,
	}
	worker, err := downloadmp4.NewTask(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if worker != nil {
		go worker.Start()
	}
	c.JSON(http.StatusOK, "ok")
}
