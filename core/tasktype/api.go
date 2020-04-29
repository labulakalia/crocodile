package tasktype

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

var _ TaskRuner = DataAPI{}

// DataAPI http req task
type DataAPI struct {
	URL     string            `json:"url" comment:"URL"`
	Method  string            `json:"method" comment:"Method"`
	PayLoad string            `json:"payload" comment:"PayLoad"`
	Header  map[string]string `json:"header" comment:"Header"`
}

// Header
// Body
// Test

// Run implment TaskRun interface
func (da DataAPI) Run(ctx context.Context) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var exitCode = DefaultExitCode
		defer pw.Close()
		defer func() {
			now := time.Now().Local().Format("2006-01-02 15:04:05: ")
			pw.Write([]byte(fmt.Sprintf("\n%sRun Finished,Return Code:%5d", now, exitCode))) // write exitCode,total 5 byte
			// pw.Write([]byte(fmt.Sprintf("%3d", exitCode))) // write exitCode,total 3 byte
		}()
		// go1.13 use NewRequestWithContext

		req, err := http.NewRequestWithContext(ctx, da.Method, da.URL, bytes.NewReader([]byte(da.PayLoad)))
		if err != nil {
			pw.Write([]byte(err.Error()))
			log.Error("NewRequest failed", zap.Error(err))
			return
		}

		for k, v := range da.Header {
			req.Header.Add(k, v)
		}

		client := http.DefaultClient
		doresp, err := client.Do(req)
		if err != nil {
			log.Error("client Do failed", zap.Error(err))
			var customerr bytes.Buffer
			switch ctx.Err() {
			case context.DeadlineExceeded:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxDeadlineExceeded))
			case context.Canceled:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxCanceled))
			default:
				customerr.WriteString(err.Error())
			}
			pw.Write(customerr.Bytes())
			return
		}

		bs, err := ioutil.ReadAll(doresp.Body)
		if err != nil {
			log.Error("Read failed", zap.Error(err))
			return
		}
		pw.Write(bs)

		//var out = make([]byte, 1024)
		//for {
		//	n, err := doresp.Body.Read(out)
		//	if err != nil {
		//		if err == io.EOF {
		//			break
		//		}
		//		log.Error("Read failed", zap.Error(err))
		//		return
		//	}
		//	if n > 0 {
		//		pw.Write(out[:n])
		//	}
		//}
		if doresp.StatusCode > 0 {
			exitCode = doresp.StatusCode
		}
	}()
	return pr
}
