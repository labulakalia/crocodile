package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"net/http"
	"strconv"
	"time"

	"github.com/labulaka521/crocodile/common/notify"
)


// getsign generate a sign when secure level is needsign
func getsign(secret string, now string) string {
	signstr := now + "\n" + secret
	// HmacSHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signstr))
	hm := h.Sum(nil)
	// Base64 encode
	b := base64.StdEncoding.EncodeToString(hm)
	// urlEncode
	sign := url.QueryEscape(b)
	return sign
}

// Secrue dingding secrue setting
// pls reading https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
type Secrue uint

const (
	// CustomKey Custom keywords
	CustomKey Secrue = iota + 1
	// Sign need sign up
	Sign
	// IPCdir IP addres
	IPCdir
)

// Ding dingding alarm conf
type Ding struct {
	MsgType string // text
	url     string
}

// Result post resp
type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type text struct {
	Content string `json:"content"`
}

type at struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

// SendMsg post json data
type SendMsg struct {
	MsgType string `json:"msgtype"`
	Text    text   `json:"text"`
	At      at     `json:"at"`
}

// NewDing init a Dingding send conf
func NewDing(webhookurl string, sl Secrue, secret string) notify.Sender {
	d := Ding{
		url:     webhookurl,
		MsgType: "text",
	}

	if sl == Sign {
		now := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		sign := getsign(secret, now)
		d.url += fmt.Sprintf("&timestamp=%s&sign=%s", now, sign)
	}
	return &d
}

// Send to notify tos is phone number
func (d *Ding) Send(tos []string, title string, content string) error {
	sendmsg := SendMsg{
		MsgType: "text",
		Text: text{
			Content: title + "\n" + content + "\n",
		},
		At: at{
			AtMobiles: tos,
			IsAtAll:   false,
		},
	}

	resp, err := notify.JSONPost(d.url, sendmsg, http.DefaultClient)
	if err != nil {
		return err
	}
	res := Result{}
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return err
	}
	if res.ErrCode != 0 {
		return fmt.Errorf("errmsg: %s errcode: %d", res.ErrMsg, res.ErrCode)
	}
	return nil
}
