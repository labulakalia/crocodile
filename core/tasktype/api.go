package tasktype

import (
	"context"
	pb "github.com/labulaka521/crocodile/core/proto"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// TODO 获取期望的值，不在于返回码为准备
type DataApi struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	PayLoad string            `json:"payload"`
	Header  map[string]string `json:"header"`
	Timeout int               `json:"timeout"` // s
}

// Header
// Body
// Test

// http req api
// form post
func (da *DataApi) Run(ctx context.Context) (taskresp *pb.TaskResp) {
	taskresp = &pb.TaskResp{}
	req, err := http.NewRequest(da.Method, da.Url, strings.NewReader(da.PayLoad))
	if err != nil {
		taskresp.Code = -1
		taskresp.ErrMsg = []byte(err.Error())
		return
	}
	for hk, hb := range da.Header {
		req.Header.Add(hk, hb)
	}
	client := http.Client{Timeout: time.Second * time.Duration(da.Timeout)}

	resp, err := client.Do(req)
	if err != nil {
		taskresp.Code = -1
		taskresp.ErrMsg = []byte(err.Error())
		return
	}
	taskresp.Code = int32(resp.StatusCode)
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		taskresp.Code = -1
		taskresp.ErrMsg = []byte(err.Error())
		return
	}
	taskresp.RespData = respbody
	return
}
