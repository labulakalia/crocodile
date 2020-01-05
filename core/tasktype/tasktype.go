package tasktype

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
)

// TaskRuner run task interface
type TaskRuner interface {
	Run(ctx context.Context) *pb.TaskResp
}

// GetDataRun get task type 
// get api or shell
func GetDataRun(t *pb.TaskReq) (TaskRuner, error) {
	switch define.TaskType(t.TaskType) {
	case define.Shell:
		var shell DataShell
		err := json.Unmarshal(t.TaskData, &shell)
		if err != nil {
			return nil, err
		}
		if len(shell.Args) == 0 {
			shell.Args = []string{}
		}
		return &shell, err

	case define.API:
		var api DataAPI
		err := json.Unmarshal(t.TaskData, &api)
		if err != nil {
			return nil, err
		}
		return &api, err

	default:
		err := fmt.Errorf("Unsupport TaskType %d", t.TaskType)
		return nil, err
	}
}
