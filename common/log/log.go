package log

import (
	"crocodile/common/cfg"
	"github.com/labulaka521/logging"
)

func Init() {
	if cfg.LogConfig.Level != "" {
		logging.SetLogLevel(cfg.LogConfig.Level)
	}
	if cfg.LogConfig.Path != "" {
		logging.SetLogPath(cfg.LogConfig.Path)
	}
	if cfg.LogConfig.Size != 0 {
		logging.SetLogSize(cfg.LogConfig.Size)
	}
	logging.Setup()
	logging.Info("Init Logging...")
}
