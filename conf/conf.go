package conf

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/skyandong/tool/conf"
	"github.com/skyandong/tool/logger"
)

// Conf 全局参数
type Conf struct {
	// App is the service
	App *struct {
		// Name is the service name
		Name string `mapstructure:"name"`
		// Ports is the service ports
		Ports map[string]int `mapstructure:"ports"`
		// Register self to consul
		Register bool `mapstructure:"register"`
		// RegisterDelay after serve
		RegisterDelay time.Duration `mapstructure:"registerDelay"`
		// ControlPort is local control port
		ControlPort int `mapstructure:"controlPort"`
		// RestartDelay for graceful restart
		RestartDelay time.Duration `mapstructure:"restartDelay"`
		// ShutdownWait for graceful shutdown
		ShutdownWait time.Duration `mapstructure:"shutdownWait"`
	} `mapstructure:"app"`

	// Loggers may contains trace logger
	Loggers logger.Conf `mapstructure:"loggers"`
}

const (
	// HTTPPortKey for HTTP port
	HTTPPortKey = "http"
	// LoggerName common logger
	LoggerName = "logger"
	// TraceLoggerName trace logger
	TraceLoggerName = "trace_logger"
	// RecLoggerName rec logger, not required
	RecLoggerName    = "rec_logger"
	configFolderName = "conf"
)

var (
	// C is the global config
	C Conf

	// Env is one of dev/test/prod
	Env = "dev"

	errNotFound = errors.New("not found")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if env := os.Getenv("ENV"); env != "" {
		Env = env
	}
	loadConfig()
	validateConfig()
}

func validateConfig() {
	for _, key := range []string{HTTPPortKey} {
		if _, ok := C.App.Ports[key]; !ok {
			log.Fatalf("%s port not found\n", key)
		}
	}
	err := C.Loggers.Ensure([]string{TraceLoggerName})
	if err != nil {
		log.Fatalln(err)
	}
}

func loadConfig() {
	configDir, err := getDir(configFolderName)
	if err != nil {
		log.Fatalf("getDir(\"%s\") error: %v", configFolderName, err)
	}
	configFileName := Env + ".conf.yaml"
	configPath := path.Join(configDir, configFileName)
	err = conf.LoadConfig(configPath, &C)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
}

func getDir(dir string) (string, error) {
	for i := 0; i < 3; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, nil
		}
		dir = filepath.Join("..", dir)
	}
	return "", errNotFound
}
