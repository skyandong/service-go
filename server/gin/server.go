package gin

import (
	"github.com/skyandong/service-go/server"
	"github.com/skyandong/tool"
	"github.com/skyandong/tool/consul"
	"github.com/skyandong/tool/program_controller"
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
	ms := gin2.Middleware(
		middleware.Cors(),
		middleware.BossWithOptions(
			cfg.Logger,
			middleware.LatencyLimit(cfg.LatencyLimit),
			middleware.ServiceName(cfg.ServiceName),
		),
	)
	ss := []gin2.Option{ms}
	if cfg.Register {
		ss = append(ss, gin2.ServiceConfigs(sc))
	}
	s := gin2.NewServer(ss...)
	g := s.Origin()
	g.POST("/get/video", getVideoFromM3u8)
	g.POST("/get/mp4", getVideoFromMp4)
	return s
}
