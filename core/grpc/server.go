package grpc

import (
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	Pem        = "grpc/tls/server.pem"
	Key        = "grpc/tls/server.key"
	ServerName = "crocodile"
)

func NewgRPCServer(mode define.RunMode) (*grpc.Server, error) {
	c, err := credentials.NewServerTLSFromFile(Pem, Key)
	if err != nil {
		return nil, errors.Wrap(err, "credentials.NewServerTLSFromFile")
	}
	grpcserver := grpc.NewServer(grpc.Creds(c))
	switch mode {
	case define.Server:
		pb.RegisterHeartbeatServer(grpcserver, &schedule.HeartbeatService{})
	case define.Client:
		pb.RegisterTaskServer(grpcserver, &schedule.TaskService{})
	}
	return grpcserver, nil
}
