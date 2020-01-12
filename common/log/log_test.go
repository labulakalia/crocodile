package log

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"go.uber.org/zap"
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
		Path(logcfg.LogPath),
		Level(logcfg.LogLevel),
		Compress(logcfg.Compress),
		MaxSize(logcfg.MaxSize),
		MaxBackups(logcfg.MaxBackups),
		MaxAge(logcfg.MaxAge),
		Format("json"),
	)
	if err != nil {
		fmt.Printf("InitLog failed: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	Info("TestLog", zap.String("test", "eeyeyyeye"))
	Debug("debug", zap.String("debug", "debug"))
	Warn("warn", zap.String("warn", "warn"))
	Error("error", zap.String("error", "error"))
	Panic("panic", zap.String("panic", "panic"))
	Fatal("fatal", zap.String("fatal", "fatal"))
}
