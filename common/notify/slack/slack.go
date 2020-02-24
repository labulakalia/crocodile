package slack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labulaka521/crocodile/common/notify"
)

// Slack send conf
type Slack struct {
	webhookurl string
	httpclient *http.Client
}

// SendMsg post json data
type SendMsg struct {
	Text string `json:"text"`
	// Channel string `json:"channel"`
	UserName string `json:"username"`
	// PreText string `json:"pretext"`
}

// NewSlack init
func NewSlack(webhook string) notify.Sender {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	return &Slack{
		webhookurl: webhook,
		httpclient: client,
	}
}

// Send will send send msg to slack channel
func (s *Slack) Send(tos []string, title string, content string) error {
	sendmsg := SendMsg{
		Text: content + "\n" + content,
	}
	resp, err := notify.JSONPost(http.MethodPost, s.webhookurl, sendmsg, s.httpclient)
	if err != nil {
		return err
	}
	if string(resp) != "ok" {
		err := fmt.Errorf("send data to slack failed error:%s", resp)
		return err
	}
	return nil
}
