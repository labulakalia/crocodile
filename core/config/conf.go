package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	CoreConf = &coreConf{}
)

func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type coreConf struct {
	SecretToken string
	Log         Log
	Pem         Pem
	Server      Server
	Client      Client
}

type Log struct {
	LogPath    string
	MaxSize    int
	Compress   bool
	MaxAge     int
	MaxBackups int
	LogLevel   string
	Format     string
}

type Pem struct {
	CertFile string
	KeyFile  string
}

type Server struct {
	Port        int
	MaxHttpTime duration
	DB          db
}

type db struct {
	Drivename    string
	Dsn          string
	MaxIdle      int
	MaxConn      int
	MaxQueryTime duration
}

type Client struct {
	Port       int
	ServerAddr string
	HostGroup  string
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
