package tasktype

import (
	"context"
	pb "github.com/labulaka521/crocodile/core/proto"
)

type DataShell struct {
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

func (ds *DataShell) Run(ctx context.Context) (*pb.TaskResp, error) {
	return nil, nil
}
