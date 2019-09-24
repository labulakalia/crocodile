package log_test

import (
	"github.com/labulaka521/crocodile/common/log"
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	importlog := &log.ImportLog{
		MaxSize:    10,
		Compress:   true,
		LogPath:    "",
		MaxAge:     0,
		MaxBackups: 0,
	}
	log.NewLogger(importlog, log.ModuleName("mymodule"), log.LogLevel(log.InfoLevel))

	log.Info("TestLog", zap.String("test", "eeyeyyeye"))

	testsss()
}

func testsss() {
	log.Info("wwww")
}
