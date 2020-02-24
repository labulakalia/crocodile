package version

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
)

var (
	// Version build version
	Version string
	// Commit git commit id
	Commit string
	// BuildDate build date
	BuildDate string
	// new version only notify once
	lastnotifyversion string
)

const (
	crocodileGithub = "https://api.github.com/repos/labulaka521/yuque_sync/releases/latest"
)

type githubapi struct {
	HTMLURL     string `json:"html_url"`
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

// CheckLatest check crocodile latest version from github
func CheckLatest() {
	cron := "00 59 23 * * ? *" // 每天的 23点59分开始运行

	expr := cronexpr.MustParse(cron)
	var (
		last time.Time
		next time.Time
	)
	last = time.Now()
	for {
		next = expr.Next(last)
		select {
		case <-time.After(next.Sub(last)):
			last = next
			go checkverson()
		}
	}

}

func checkverson() {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	getreq, err := http.NewRequestWithContext(ctx, http.MethodGet, crocodileGithub, nil)
	if err != nil {
		log.Error("http.NewRequest failed", zap.Error(err))
		return
	}
	resp, err := client.Do(getreq)
	if err != nil {
		log.Error("client.Do failed", zap.Error(err))
		return

	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Request StatusCode is not 200")
		return
	}
	data := githubapi{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("ioutil.ReadAll", zap.Error(err))
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("json.Unmarshal", zap.Error(err))
		return
	}

	// 如果发布时间与检测当日日期相同就生成一个新版本发布通知
	tparse, err := time.Parse("2006-01-02T15:04:05Z", data.PublishedAt)
	if err != nil {
		log.Error("time.Parse failed", zap.Error(err))
		return
	}

	// check published date equal today
	if time.Now().Format("2006-01-02") != tparse.Format("2006-01-02") {
		log.Debug("published is expired", zap.String("published", data.PublishedAt))
		return
	}
	users, _, err := model.GetUsers(ctx, nil, 0, 0)
	if err != nil {
		log.Error("model.GetUsers failed", zap.Error(err))
	}

	for _, user := range users {
		if user.Role != define.AdminUser {
			continue
		}
		notify := define.Notify{
			NotifyType: define.UpgradeNotify,
			Title:      data.Name,
			NotifyUID:  user.ID,
			Content:    data.Body,
			NotifyTime: tparse.Unix(),
		}
		err = model.SaveNewNotify(ctx, notify)
		if err != nil {
			log.Error("model.SaveNewNotify failed", zap.Error(err))
		}
	}

}
