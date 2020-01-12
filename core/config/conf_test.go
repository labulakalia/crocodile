package config

import (
	"testing"
	"os"
	"io/ioutil"
)

var conf = `# secret token
secrettoken = "weinjuwiwiuwu"

# log
[log]
logpath = ""
maxsize = 10
compress = true
maxage =  7
maxbackups = 10
loglevel  = "info"
format = "text"

# cert file
[cert]
certfile="core/cert/cert.pem"
keyfile="core/cert/key.pem"

# crocodile server
[server]
port = 8080
maxhttptime = "10s" # 秒
  [server.db]
  drivename = "sqlite3"
  dsn = "db/core.db"
  maxidle = 10
  maxconn = 20
  maxquerytime = "10s"

# crocodile client
[client]
port = 8081        # default rand port
serveraddr = "127.0.0.1:8080"
hostgroup = "crocodile_hostgroup"`

func TestInit(t *testing.T) {
	testfile := "/tmp/crocodile.toml"
	ioutil.WriteFile(testfile, []byte(conf), 0644)
	Init(testfile)
	os.Remove(testfile)
}
