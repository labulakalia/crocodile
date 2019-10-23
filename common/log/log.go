package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var (
	_logger *zap.Logger
)

type Level uint

// 日志配置
type logConfig struct {
	LogPath    string
	LogLevel   string
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

func getzapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logpath string, maxsize int, compress bool) io.Writer {
	if logpath == "" {
		return os.Stdout
	} else {
		return &lumberjack.Logger{
			Filename: logpath,
			MaxSize:  maxsize,
			Compress: compress,
		}
	}
}

func newZapEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}
func newLoggerCore(log *logConfig) zapcore.Core {
	hook := newLogWriter(log.LogPath, log.MaxSize, log.Compress)

	encoderConfig := newZapEncoder()

	atomLevel := zap.NewAtomicLevelAt(getzapLevel(log.LogLevel))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)),
		atomLevel,
	)
	return core
}

func newLoggerOptions() []zap.Option {
	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(1)
	// 开发者
	zap.Fields()
	development := zap.Development()
	options := []zap.Option{
		caller,
		callerskip,
		development,
	}
	return options
}

type Option func(*logConfig)

func LogPath(logpath string) Option {
	return func(logcfg *logConfig) {
		logcfg.LogPath = logpath
	}
}

func Compress(compress bool) Option {
	return func(logcfg *logConfig) {
		logcfg.Compress = compress
	}
}

func LogLevel(level string) Option {
	return func(logcfg *logConfig) {
		logcfg.LogLevel = level
	}
}

func MaxSize(size int) Option {
	return func(logcfg *logConfig) {
		logcfg.MaxSize = size
	}
}

func MaxAge(age int) Option {
	return func(logcfg *logConfig) {
		logcfg.MaxAge = age
	}
}

func MaxBackups(backup int) Option {
	return func(logcfg *logConfig) {
		logcfg.MaxBackups = backup
	}
}

func defaultOption() *logConfig {
	return &logConfig{
		LogPath:    "",
		MaxSize:    10,
		Compress:   false,
		MaxAge:     7,
		MaxBackups: 7,
		LogLevel:   "info",
	}
}

func NewLog(opts ...Option) error {

	logcfg := defaultOption()
	for _, opt := range opts {
		opt(logcfg)
	}
	core := newLoggerCore(logcfg)

	zapopts := newLoggerOptions()
	_logger = zap.New(core, zapopts...)
	return nil
}

func Debug(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
