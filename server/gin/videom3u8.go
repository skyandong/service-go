package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/skyandong/autils/service"
	"github.com/skyandong/service-go/api"
	"github.com/skyandong/service-go/service/core/downloadm3u8"
)

func getVideoFromM3u8(c *gin.Context) {
	var req api.GetVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ctx := &downloadm3u8.Context{
		Logger:          tracelog,
		URL:             req.URL,
		FileName:        req.FileName,
		DownloadCatalog: req.DepositAddress,
		TraceID:         service.GetTraceID(c),
	}
	if err := downloadm3u8.Work(ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "ok")
}
