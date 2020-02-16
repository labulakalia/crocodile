package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	// CoreConf crocodile conf
	CoreConf *coreConf
)

// Init Config
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
	Cert        Cert
	Server      Server
	Client      Client
	Notify      Notify
}

// Log Config
type Log struct {
	LogPath    string
	MaxSize    int
	Compress   bool
	MaxAge     int
	MaxBackups int
	LogLevel   string
	Format     string
}

// Cert tls cert
type Cert struct {
	CertFile string
	KeyFile  string
}

// Server crocodile server config
type Server struct {
	Port        int
	MaxHTTPTime duration
	DB          db
}

type db struct {
	Drivename    string
	Dsn          string
	MaxIdle      int
	MaxConn      int
	MaxQueryTime duration
}

// Client crocodile client config
type Client struct {
	Port       int
	ServerAddr string
	HostGroup  string
	Weight     int
	Remark     string
}

type duration struct {
	time.Duration
}

// UnmarshalText parse 10s to time.Time
func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// Notify send msg to user
type Notify struct {
	Email    email
	DingDing dingding
	Slack    slack
	Telegram telegram
	WeChat   wechat
	WebHook  webhook
}

type email struct {
	SMTPHost   string
	Port       int
	UserName   string
	Password   string
	From       string
	TLS        bool
	Anonymous  bool
	SkipVerify bool
}

type dingding struct {
	WebHook     string
	SecureLevel int
	Secret      string
}

type slack struct {
	WebHook string
}

type telegram struct {
	BotToken string
}

type wechat struct {
	CropID      string
	AgentID     int
	AgentSecret string
}

type webhook struct {
	WebHookURL string
}
