package gin

import (
	"github.com/skyandong/autils/consul"
	"github.com/skyandong/autils/controller"
	"github.com/skyandong/autils/logger"
	"github.com/skyandong/autils/server/gin"
	"github.com/skyandong/autils/server/gin/middleware"
	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/server"
)

var tracelog = conf.C.Loggers.Get(conf.LoggerName).GetLogger(logger.InfoLevel)

// New creates a gin server
func New(cfg *server.Config) controller.Server {
	sc := []*consul.ServiceConf{
		{
			Name:    cfg.ServiceName,
			Address: cfg.Address,
			Port:    cfg.Port,
		},
		{
			Name:    "BusinessExporter",
			Address: cfg.Address,
			Port:    cfg.Port,
			Tags:    []string{"serviceName=" + cfg.ServiceName},
		},
	}
	ms := gin.Middleware(
		middleware.Cors(),
		middleware.BossWithOptions(
			cfg.Logger,
			middleware.LatencyLimit(cfg.LatencyLimit),
			middleware.ServiceName(cfg.ServiceName),
		),
	)
	ss := []gin.Option{ms}
	if cfg.Register {
		ss = append(ss, gin.ServiceConfigs(sc))
	}
	s := gin.NewServer(ss...)
	g := s.Origin()
	g.POST("/get/video", getVideoFromM3u8)
	g.POST("/get/mp4", getVideoFromMp4)
	return s
}
