package tasktype

import (
	"context"
	"github.com/labulaka521/crocodile/common/log"
	pb "github.com/labulaka521/crocodile/core/proto"
	"go.uber.org/zap"
)

type DataShell struct {
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

func (ds *DataShell) Run(ctx context.Context) (*pb.TaskResp, error) {
	log.Info("Start Run Command", zap.String("Name", ds.Name), zap.Any("args", ds.Args))
	resp := &pb.TaskResp{
		Code:     -1,
		RespData: []byte(" 111"),
	}

	return resp, nil
}
