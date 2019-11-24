package schedule

import (
	"context"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/peer"
	"net"
)

// run worker
// implementation proto task interface

type TaskService struct {
}

func (ts *TaskService) RunTask(ctx context.Context, t *pb.TaskReq) (*pb.TaskResp, error) {
	log.Info("runTask", zap.Any("task", t))
	var (
		taskresp *pb.TaskResp
	)
	r, err := tasktype.GetDataRun(t)
	if err != nil {
		taskresp = &pb.TaskResp{
			Code:     -1,
			ErrMsg:   []byte(err.Error()),
			RespData: []byte(err.Error()),
		}
		return taskresp, err
	}
	taskresp, err = r.Run(ctx)
	if err != nil {
		log.Info("Run failed", zap.String("error", err.Error()))
		return taskresp, err
	}
	log.Info("TaskResp", zap.Any("resp", taskresp))
	return taskresp, nil
}

// run core server
// implementation proto task interface
type HeartbeatService struct{}

func (hs *HeartbeatService) RegistryHost(ctx context.Context, req *pb.RegistryReq) (*pb.Empty, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Empty{}, errors.New("Registry failed")
	}
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	req.Ip = ip

	// check hostgroup
	if req.Hostgroup != "" {
		hgs, err := model.GetHostGroupName(ctx, req.Hostgroup)
		if err != nil {
			return &pb.Empty{}, err
		}
		addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)
		exist := false
		for _, hgaddr := range hgs.Addrs {
			if addr == hgaddr {
				exist = true
			}
		}
		if !exist {
			hgs.Addrs = append(hgs.Addrs, addr)
		}

		err = model.ChangeHostGroup(ctx, hgs)
		if err != nil {
			return &pb.Empty{}, err
		}
	}
	log.Info("New Client Registry ", zap.String("ip", req.Ip), zap.Int32("port", req.Port))
	err := model.RegistryNewHost(ctx, req)
	return &pb.Empty{}, err
}

func (hs *HeartbeatService) SendHb(ctx context.Context, hb *pb.HeartbeatReq) (*pb.Empty, error) {
	err := model.UpdateRunningTask(ctx, hb)
	return &pb.Empty{}, err
}
