package log

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logcfg := &logConfig{
		MaxSize:    10,
		Compress:   true,
		LogPath:    "",
		MaxAge:     0,
		MaxBackups: 0,
		LogLevel:   "info",
	}
	err := InitLog(
		LogPath(logcfg.LogPath),
		LogLevel(logcfg.LogLevel),
		Compress(logcfg.Compress),
		MaxSize(logcfg.MaxSize),
		MaxBackups(logcfg.MaxBackups),
		MaxAge(logcfg.MaxAge),
	)
	if err != nil {
		fmt.Printf("InitLog failed: %v", err)
		os.Exit(1)
	}

	Info("TestLog", zap.String("test", "eeyeyyeye"))

	testsss()

}

func testsss() {
	Info("wwww")
}
