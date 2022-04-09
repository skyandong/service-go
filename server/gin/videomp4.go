package gin

import (
	downloadmp42 "github.com/skyandong/service-go/service/download/downloadmp4"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/skyandong/autils/service"
	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/service/core/downloadmp4"
)

func getVideoFromMp4(c *gin.Context) {
	req := api.GetVideoRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := &downloadmp42.Context{
		Ctx:             c,
		TraceID:         service.GetTraceID(c),
		URL:             req.URL,
		DownloadCatalog: req.DepositAddress,
		FileName:        req.FileName,
		Logger:          tracelog,
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
