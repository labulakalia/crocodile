package tasktype

import (
	"context"
	pb "github.com/labulaka521/crocodile/core/proto"
)

type DataApi struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	PayLoad string            `json:"payload"`
	Header  map[string]string `json:"header"`
}

func (data *DataApi) Run(ctx context.Context) (*pb.TaskResp, error) {
	return nil, nil
}
