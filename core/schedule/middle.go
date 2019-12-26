package schedule

import (
	"context"
	"errors"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	pb "github.com/labulaka521/crocodile/core/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"time"
)

// grpc middleware

func LoggerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	resp, err := handler(ctx, req)
	if err != nil {
		log.Error("resp failed", zap.Error(err))
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Empty{}, errors.New("Registry failed")
	}

	log.Info("[rpc req]", zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Any("resp", resp),
		zap.Any("reqaddr", p.Addr.String()),
		zap.Duration("latency(ms)", time.Now().Sub(start)*1000),
	)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic Err: %v", err)
		}
	}()
	return handler(ctx, req)

}

func CheckSecretInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "can not get token")
	}
	var secrettoken string
	v, ok := md["secret_token"]
	if ok {
		secrettoken = v[0]
	}
	if secrettoken != config.CoreConf.SecretToken {
		return nil, status.Errorf(codes.Unauthenticated, "secrettoken auth failed")
	}
	return handler(ctx, req)
}
