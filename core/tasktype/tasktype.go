package tasktype

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
)

const (
	// DefaultExitCode default err code if not get run task code
	DefaultExitCode int = -1
)

// TaskRuner run task interface
// Please Implment io.ReadCloser
// reader last 3 byte must be exit code
type TaskRuner interface {
	Run(ctx context.Context) (out io.ReadCloser)
}

// GetDataRun get task type
// get api or code
func GetDataRun(t *pb.TaskReq) (TaskRuner, error) {
	switch define.TaskType(t.TaskType) {
	case define.Code:
		var code DataCode
		err := json.Unmarshal(t.TaskData, &code)
		if err != nil {
			return nil, err
		}
		code.LangDesc = code.Lang.String()
		return &code, err

	case define.API:
		var api DataAPI
		err := json.Unmarshal(t.TaskData, &api)
		if err != nil {
			return nil, err
		}
		if api.Header == nil {
			api.Header = make(map[string]string)
		}
		return &api, err

	default:
		err := fmt.Errorf("Unsupport TaskType %d", t.TaskType)
		return nil, err
	}
}
