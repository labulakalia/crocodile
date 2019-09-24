package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var (
	defaultLogLevel  = InfoLevel
	defaultLogModule = "logger"
	logger           *zap.Logger
)

type ImportLog struct {
	LogPath    string
	MaxSize    int
	Compress   bool
	MaxAge     int
	MaxBackups int
}

type logOption struct {
	*ImportLog
	// 模块名称
	ModuleName string

	LogLevel Level
}

type LogOption interface {
	apply(*logOption)
}

type funcLogOption struct {
	f func(*logOption)
}

func (fdo *funcLogOption) apply(do *logOption) {
	fdo.f(do)
}

func newFuncLogOption(f func(option *logOption)) *funcLogOption {
	return &funcLogOption{f: f}
}

func ModuleName(moduleName string) LogOption {
	return newFuncLogOption(func(o *logOption) {
		o.ModuleName = moduleName
	})
}

func LogLevel(loglevel Level) LogOption {
	return newFuncLogOption(func(o *logOption) {
		o.LogLevel = loglevel
	})
}

// if params is zero, like LogPath ,log will console to os.Stdout
// ServerName is zero,
func defaultOption(log *ImportLog) *logOption {
	return &logOption{
		ImportLog:  log,
		ModuleName: defaultLogModule,
		LogLevel:   defaultLogLevel,
	}
}

func getzapLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case PanicLevel:
		return zap.PanicLevel
	case FatalLevel:
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

func NewLogger(log *ImportLog, opt ...LogOption) {
	if log == nil {
		panic("Want Params Log, But Can get nil")
	}
	opts := defaultOption(log)

	for _, o := range opt {
		o.apply(opts)
	}

	hook := newLogWriter(opts.LogPath, opts.MaxSize, opts.Compress)

	encoderConfig := newZapEncoder()

	atomLevel := zap.NewAtomicLevelAt(getzapLevel(opts.LogLevel))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)),
		atomLevel,
	)
	// 文件行号
	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(1)
	// 开发者
	development := zap.Development()
	// 初始化字段
	filed := zap.Fields(zap.String("ModuleName", opts.ModuleName))

	logger = zap.New(core, caller, callerskip, development, filed)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
