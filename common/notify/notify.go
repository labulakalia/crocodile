package notify

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/labulaka521/crocodile/common/log"
	"go.uber.org/zap"
)

// Sender it send notify to user
type Sender interface {
	Send(to []string, title string, content string) error
}

// alarm notify
// mail
// chat
// dingding
// slack
// telegram
// server jiang

//JSONPost Post req json data to url
func JSONPost(url string, data interface{}, client *http.Client) ([]byte, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("client.Do", zap.Error(err))
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
