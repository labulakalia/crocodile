package schedule

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/cert"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"os"
	"sync"
	"time"
)


const (
	// RPC connect timeout
	defaultRPCTimeout       = time.Second * 3
	// worker send hearbeat ttl
	defaultHearbeatInterval = time.Second * 5
)


var (
	// rpc conn control
	cachegRPCConnM *cachegRPCConn
)

func init() {
	cachegRPCConnM = &cachegRPCConn{
		conn: make(map[string]*grpc.ClientConn),
	}
}

type cachegRPCConn struct {
	sync.RWMutex
	conn map[string]*grpc.ClientConn
}

// getgRPCClientConn return conn or nil
func (cg *cachegRPCConn) getgRPCClientConn(addr string) *grpc.ClientConn {
	cg.Lock()
	conn, exist := cg.conn[addr]
	cg.Unlock()
	if exist && conn.GetState() == connectivity.Ready {
		return conn
	}
	conn.Close()
	return nil
}

func (cg *cachegRPCConn) addgRPCClientConn(addr string, conn *grpc.ClientConn) {
	cg.Lock()
	cg.conn[addr] = conn
	cg.Unlock()
}

// NewgRPCConn Get Grpc Client Conn
func NewgRPCConn(ctx context.Context,addr string) (*grpc.ClientConn, error) {
	conn := cachegRPCConnM.getgRPCClientConn(addr)
	if conn != nil {
		return conn, nil
	}

	c, err := credentials.NewClientTLSFromFile(config.CoreConf.Cert.CertFile, cert.ServerName)
	if err != nil {
		log.Error("credentials.NewClientTLSFromFile failed", zap.Error(err))
		return nil, err
	}
	rpcctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	conn, err = grpc.DialContext(rpcctx, addr, grpc.WithTransportCredentials(c),
		grpc.WithPerRPCCredentials(
			&Auth{SecretToken: config.CoreConf.SecretToken},
		),
		grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}
	cachegRPCConnM.addgRPCClientConn(addr, conn)
	return conn, nil
}

// NewgRPCServer new gRPC server
func NewgRPCServer(mode define.RunMode) (*grpc.Server, error) {
	c, err := credentials.NewServerTLSFromFile(config.CoreConf.Cert.CertFile, config.CoreConf.Cert.KeyFile)
	if err != nil {
		log.Error("credentials.NewServerTLSFromFile failed", zap.Error(err))
		return nil, err
	}
	auth := Auth{SecretToken: config.CoreConf.SecretToken}
	grpcserver := grpc.NewServer(grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			middleware.RecoveryInterceptor,
			middleware.LoggerInterceptor,
			middleware.CheckSecretInterceptor,
		),
	)
	switch mode {
	case define.Server:
		pb.RegisterHeartbeatServer(grpcserver, &HeartbeatService{Auth: auth})
	case define.Client:
		pb.RegisterTaskServer(grpcserver, &TaskService{Auth: auth})
	}
	return grpcserver, nil
}

// Auth check rpc request valid
type Auth struct {
	SecretToken string
}

// GetRequestMetadata implement PerRPCCredentials interface
func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"secret_token": a.SecretToken,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (a *Auth) RequireTransportSecurity() bool {
	return true
}

// try Get grpc client conn
// TODO
// 调度算法
// - 随机
// - 最少任务数
// - 权重
// - 轮询执行
// get rpc conn
func tryGetRCCConn(ctx context.Context, hg *define.HostGroup) (*grpc.ClientConn, error) {
	queryctx, querycancel := context.WithTimeout(ctx,
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer querycancel()
	i := 0
	for i < len(hg.HostsID) {
		i++
		hostid, err := model.RandHostID(hg)
		if err != nil {
			log.Error("model.RandHostID failed", zap.String("error", err.Error()))
			continue
		}

		host, err := model.GetHostByID(queryctx, hostid)
		if err != nil {
			log.Error("model.GetHostByID failed", zap.String("error", err.Error()))
			continue
		}

		if host.Online == 0 {
			log.Warn("host is offline", zap.String("addr", host.Addr))
			continue
		}
		if host.Stop == 1 {
			log.Warn("host is stop worker", zap.String("addr", host.Addr))
			continue
		}
		conn, err := NewgRPCConn(ctx,host.Addr)
		if err != nil {
			log.Error("GetRpcConn failed", zap.String("error", err.Error()))
			continue
		}
		// when only conn is Ready, direct return this conn,otherse
		if conn.GetState() == connectivity.Ready {
			return conn, nil
		}
		conn.Close()
	}
	return nil, fmt.Errorf("can not get valid grpc conn from hostgroup: %s",hg.Name)
}

// RegistryClient registry client to server
func RegistryClient(version string, port int) error {
	conn, err := NewgRPCConn(context.Background(),config.CoreConf.Client.ServerAddr)
	if err != nil {
		return err
	}
	hbClient := pb.NewHeartbeatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
	defer cancel()
	hostname, _ := os.Hostname()
	regHost := pb.RegistryReq{
		Port:      int32(port),
		Hostname:  hostname,
		Version:   version,
		Hostgroup: config.CoreConf.Client.HostGroup,
	}
	_, err = hbClient.RegistryHost(ctx, &regHost)
	if err != nil {
		return errors.Errorf("can not connect server %s", conn.Target())

	}
	log.Info("host registry success")
	go sendhb(hbClient, port)
	return nil
}

// send hearbt
func sendhb(client pb.HeartbeatClient, port int) {
	log.Info("start send hearbeat to server")
	timer := time.NewTimer(defaultHearbeatInterval)
	for {
		select {
		case <-timer.C:
			ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
			defer cancel()
			hbreq := &pb.HeartbeatReq{Port: int32(port)}
			_, err := client.SendHb(ctx, hbreq)
			if err != nil {
				log.Error("client.SendHb failed", zap.Error(err))
			}
			timer.Reset(defaultHearbeatInterval)
		}
	}
}

func dealRPCErr(err error) int {
	statusErr, ok := status.FromError(err)
	if ok {
		switch statusErr.Code() {
		case codes.DeadlineExceeded:
			return resp.ErrCtxDeadlineExceeded
		case codes.Canceled:
			return resp.ErrCtxCanceled
		case codes.Unauthenticated:
			return resp.ErrRPCUnauthenticated
		case codes.Unavailable:
			return resp.ErrRPCUnavailable
		}
	}
	return resp.ErrRPCUnknow
}
