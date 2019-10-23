package log

import (
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"os"
)

func InitLog() {
	logcfg := config.CoreConf.Log

	err := log.NewLog(
		log.LogPath(logcfg.LogPath),
		log.LogLevel(logcfg.LogLevel),
		log.Compress(logcfg.Compress),
		log.MaxSize(logcfg.MaxSize),
		log.MaxBackups(logcfg.MaxBackups),
		log.MaxAge(logcfg.MaxAge),
	)
	if err != nil {
		fmt.Printf("InitLog failed: %v", err)
		os.Exit(1)
	}
}
