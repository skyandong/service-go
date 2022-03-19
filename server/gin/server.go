package gin

import (
	"github.com/skyandong/service-go/server"
	"github.com/skyandong/tool/consul"
	"github.com/skyandong/tool/program_controller"
	"github.com/skyandong/tool/server/gin"
	"github.com/skyandong/tool/server/gin/middleware"
)

// New creates a gin server
func New(cfg *server.Config) program_controller.Server {
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
