package tasktype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
)

type TaskRuner interface {
	Run(ctx context.Context) *pb.TaskResp
}

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

	case define.Api:
		var api DataApi
		err := json.Unmarshal(t.TaskData, &api)
		if err != nil {
			return nil, err
		}
		return &api, err

	default:
		err := errors.New(fmt.Sprintf("Unsupport TaskType %d", t.TaskType))
		return nil, err
	}
}
