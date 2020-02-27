package config

import (
	"io/ioutil"
	"os"
	"testing"
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
hostgroup = "crocodile_hostgroup"
weight = 100

# 消息通知配置
[notify]
# 邮箱
[notify.email]
smtphost = "smtp.163.com"
port = 465
username = "...@163.com"
password = "password"
from = "...@163.com"
tls = true
# 匿名发送
anonymous = false
# 如使用自建邮件系统请设置 skipVerify 为 true 以避免证书校验错误
skipverify = false
# 钉钉
[notify.dingding]
webhook = "dingdingurl"
# 安全设置
# 1 自定义关键字
# 2 加签
# 3 IP地址
securelevel = 1
# 如果securelevel 为2 需要填写加签密钥
secret = ""
# slack
[notify.slack]
webhook = "url"
# telegram
[notify.telegram]
bottoken = "bottoken"
# 企业微信
[notify.wechat]
cropid = "cropid"
agentid = 100002
agentsecret = "agentsecret"
[notify.webhook]
enable = true
webhookurl = "http://webhook.test"
`

func TestInit(t *testing.T) {
	testfile := "/tmp/crocodile.toml"
	ioutil.WriteFile(testfile, []byte(conf), 0644)
	Init(testfile)
	t.Logf("%+v", CoreConf.Notify.WebHook)
	os.Remove(testfile)
}
