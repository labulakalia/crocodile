package wechat

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labulaka521/crocodile/common/notify"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	defaultMsgType = "text"
)

// Err 微信返回错误
type err struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//AccessToken 微信企业号请求Token
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	err
	ExpiresInTime time.Time
}

//Client 微信企业号应用配置信息
type client struct {
	CropID      string
	AgentID     int
	AgentSecret string
	Token       accessToken
}

//Result 发送消息返回结果
type Result struct {
	err
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"infvalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

//Content 文本消息内容
type Content struct {
	Content string `json:"content"`
}

//Message 消息主体参数
type Message struct {
	ToUser  string  `json:"touser"`
	ToParty string  `json:"toparty"`
	ToTag   string  `json:"totag"`
	MsgType string  `json:"msgtype"`
	AgentID int     `json:"agentid"`
	Text    Content `json:"text"`
}

//NewWeChat init wechat notidy
func NewWeChat(cropID string, agentID int, agentSecret string) notify.Sender {
	return &client{
		CropID:      cropID,
		AgentID:     agentID,
		AgentSecret: agentSecret,
	}
}

// Send format send msg to Message
func (c *client) Send(tos []string, title, content string) error {
	msg := Message{
		ToUser:  strings.Join(tos, "|"),
		MsgType: defaultMsgType,
		Text: Content{
			Content: title + "\n" + content,
		},
		AgentID: c.AgentID,
	}
	err := c.send(msg)
	if err != nil {
		return err

	}
	return nil

}

//Send 发送信息
func (c *client) send(msg Message) error {
	c.generateAccessToken()

	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + c.Token.AccessToken
	resultByte, err := notify.JSONPost(http.MethodPost, url, msg, http.DefaultClient)
	if err != nil {
		err = errors.New("请求微信接口失败: " + err.Error())
		return err
	}
	result := Result{}
	err = json.Unmarshal(resultByte, &result)
	if err != nil {
		err = errors.New("解析微信接口返回数据失败: " + err.Error())
		return err
	}

	if result.ErrCode != 0 {
		err = errors.New("发送消息失败: " + result.ErrMsg)
		return err

	}

	if result.InvalidUser != "" || result.InvalidTag != "" || result.InvalidParty != "" {
		err = fmt.Errorf("消息发送成功, 但是有部分目标无法送达: %s%s%s", result.InvalidUser, result.InvalidParty, result.InvalidTag)
		return err
	}
	return nil
}

//generateAccessToken 生成会话token
func (c *client) generateAccessToken() {
	var err error
	if c.Token.AccessToken == "" || c.Token.ExpiresInTime.Before(time.Now()) {
		c.Token, err = getAccessTokenFromWeixin(c.CropID, c.AgentSecret)
		if err != nil {
			return
		}
		c.Token.ExpiresInTime = time.Now().Add(time.Duration(c.Token.ExpiresIn-1000) * time.Second)
	}
}

//从微信服务器获取token
func getAccessTokenFromWeixin(cropID, secret string) (TokenSession accessToken, err error) {
	WxAccessTokenURL := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + cropID + "&corpsecret=" + secret

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	result, err := client.Get(WxAccessTokenURL)
	if err != nil {
		return
	}

	res, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}

	defer result.Body.Close()

	err = json.Unmarshal(res, &TokenSession)
	if err != nil {
		return
	}

	if TokenSession.ExpiresIn == 0 || TokenSession.AccessToken == "" {
		err = fmt.Errorf("获取微信错误代码: %v, 错误信息: %v", TokenSession.ErrCode, TokenSession.ErrMsg)
		return
	}

	return
}
