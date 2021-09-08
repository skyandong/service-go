package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/skyandong/tool/consul"
	"github.com/skyandong/tool/logger"
	"github.com/skyandong/tool/program_controller"

	"github.com/skyandong/service-go/conf"
	"github.com/skyandong/service-go/server"
	"github.com/skyandong/service-go/server/gin"
)

func main() {
	if os.Getenv("DEBUG_PPROF") == "true" {
		err := program_controller.AddServer("tcp", "0.0.0.0:6060", &http.Server{})
		if err != nil {
			log.Fatalf("add pprof server error: %v", err)
		}
	}

	cPort := conf.C.App.ControlPort
	if cPort > 0 {
		err := program_controller.AddServer("tcp", "127.0.0.1:"+strconv.FormatInt(int64(cPort), 10),
			program_controller.NewControlServer())
		if err != nil {
			log.Fatalf("add control server on port %d error: %v", cPort, err)
		}
	}

	var localIP string
	var err error
	if conf.C.App.Register {
		localIP, err = consul.LocalIP()
		if err != nil {
			log.Fatalf("cannot get local ip: %v", err)
		}
	}
	hPort := conf.C.App.Ports[conf.HTTPPortKey]
	cfg := &server.Config{
		Register:     conf.C.App.Register,
		ServiceName:  conf.C.App.Name,
		Address:      localIP,
		Port:         hPort,
		Logger:       conf.C.Loggers.Get(conf.TraceLoggerName).GetLogger(logger.InfoLevel),
		LatencyLimit: conf.C.ErrorLatency,
	}
	err = program_controller.AddServer("tcp", ":"+strconv.FormatInt(int64(hPort), 10), gin.New(cfg))
	if err != nil {
		log.Fatalf("add http server on port %d error: %v", hPort, err)
	}

	err = program_controller.RunServers(conf.C.App.RestartDelay, conf.C.App.ShutdownWait)
	if err != nil {
		log.Fatalf("run servers error: %v", err)
	}
}
