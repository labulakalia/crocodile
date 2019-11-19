package schedule

import (
	"context"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
)

// run worker
// implementation proto task interface

type TaskService struct {
}

func (ts *TaskService) RunTask(ctx context.Context, t *pb.TaskReq) (*pb.TaskResp, error) {
	r, err := tasktype.GetDataRun(t)
	if err != nil {
		return nil, err
	}

	return r.Run(ctx)
}

// run core server
// implementation proto task interface
type HeartbeatService struct{}

func (hs *HeartbeatService) SendHb(ctx context.Context, hb *pb.HeartbeatReq) (*pb.Empty, error) {
	err := model.UpdateHost(ctx, hb)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
