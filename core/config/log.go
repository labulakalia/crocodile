package config

type LogConfig struct {
	LogPath    string
	MaxSize    int
	Compress   bool
	MaxAge     int
	MaxBackups int
	LogLevel   string
}
