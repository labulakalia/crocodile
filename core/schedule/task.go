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
	"strings"
)

// TaskService implementation proto task interface
type TaskService struct {
	Auth Auth
}

// RunTask run task by rpc
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
			RespData: nil,
		}
		return taskresp, err
	}
	taskresp = r.Run(ctx)
	log.Info("TaskResp", zap.Any("resp", taskresp))
	return taskresp, nil
}

// HeartbeatService implementation proto Heartbeat interface
type HeartbeatService struct {
	Auth Auth
}

// RegistryHost client registry
func (hs *HeartbeatService) RegistryHost(ctx context.Context, req *pb.RegistryReq) (*pb.Empty, error) {
	var (
		id string
	)
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Empty{}, errors.New("Registry failed")
	}
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	req.Ip = ip
	addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)
	host, exist, err := model.ExistAddr(ctx, addr)
	if err != nil {
		return &pb.Empty{}, err
	}
	if !exist {
		id, err = model.RegistryNewHost(ctx, req)
		if err != nil {
			return &pb.Empty{}, err
		}

	} else {
		id = host.ID
	}
	hb := pb.HeartbeatReq{
		Ip:   ip,
		Port: req.Port,
	}
	_, err = hs.SendHb(ctx, &hb)
	if err != nil {
		log.Error("Send hearbeat failed", zap.String("error", err.Error()))
		return &pb.Empty{}, err
	}
	if req.Hostgroup != "" {
		hgs, err := model.GetHostGroupName(ctx, req.Hostgroup)
		if err != nil {
			return &pb.Empty{}, err
		}

		if !strings.Contains(strings.Join(hgs.HostsID, ""), id) {
			hgs.HostsID = append(hgs.HostsID, id)
			err = model.ChangeHostGroup(ctx, hgs)
			if err != nil {
				return &pb.Empty{}, err
			}
		}
	}
	log.Info("New Worker Registry Success", zap.String("addr", addr))
	return &pb.Empty{}, err
}

// SendHb client send hearbeat
func (hs *HeartbeatService) SendHb(ctx context.Context, hb *pb.HeartbeatReq) (*pb.Empty, error) {

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Empty{}, errors.New("get peer failed")
	}
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	hb.Ip = ip
	log.Info("recv hearbeat", zap.String("addr", fmt.Sprintf("%s:%d", ip, hb.Port)))
	err := model.UpdateHostHearbeat(ctx, hb)
	return &pb.Empty{}, err
}
