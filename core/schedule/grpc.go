package schedule

import (
	"context"
	"os"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/cert"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
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
)

const (
	// RPC connect timeout
	defaultRPCTimeout = time.Second * 3
	// worker send hearbeat ttl
	defaultHearbeatInterval         = time.Second * 15 // maxWorkerTTL int64 = 20
	defaultLastFailHearBeatInterval = time.Second * 3
	// max retry get host time for func Next
	defaultMaxRetryGetWorkerHost = 3
)

var (
	// grpc conn pool
	cachegRPCConnM *cachegRPCConn
	// stop sent hearbeat to server
	clentstophb chan struct{}
)

func init() {
	cachegRPCConnM = &cachegRPCConn{
		conn: make(map[string]*grpc.ClientConn),
	}
	clentstophb = make(chan struct{}, 1)
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
	if conn != nil {
		conn.Close()
	}

	return nil
}

func (cg *cachegRPCConn) addgRPCClientConn(addr string, conn *grpc.ClientConn) {
	cg.Lock()
	cg.conn[addr] = conn
	cg.Unlock()
}

// getgRPCConn Get Grpc Client Conn
func getgRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
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
		return nil, err
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
			// middleware.LoggerInterceptor,
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
// Scheduling Algorithm
// - random
// - LeastTask
// - Weight
// - roundRobin
// get rpc conn
func tryGetRCCConn(ctx context.Context, next Next) (*grpc.ClientConn, error) {
	// queryctx, querycancel := context.WithTimeout(ctx,
	// 	config.CoreConf.Server.DB.MaxQueryTime.Duration)
	// defer querycancel()
	var (
		err  error
		conn *grpc.ClientConn
	)
	for i := 0; i < defaultMaxRetryGetWorkerHost; i++ {
		host := next()
		if host == nil {
			err = errors.New("Can't Get Valid Worker Host")
			continue
		}
		conn, err = getgRPCConn(ctx, host.Addr)
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
	return nil, err
}

// RegistryClient registry client to server
func RegistryClient(version string, port int) error {
	conn, err := getgRPCConn(context.Background(), config.CoreConf.Client.ServerAddr)
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
		Weight:    int32(config.CoreConf.Client.Weight),
		Remark:    config.CoreConf.Client.Remark,
	}
	_, err = hbClient.RegistryHost(ctx, &regHost)
	if err != nil {
		log.Error("registry client failed", zap.Error(err))
		return errors.Errorf("can not connect server %s", conn.Target())

	}
	log.Info("host registry success")
	go sendhb(hbClient, port)
	return nil
}

// send client will send hearbt to server, let scheduler center know it is alive
func sendhb(client pb.HeartbeatClient, port int) {
	log.Info("start send hearbeat to server")
	timer := time.NewTimer(time.Millisecond)

	for {
		select {
		case <-timer.C:
			ctx, _ := context.WithTimeout(context.Background(), defaultRPCTimeout)
			hbreq := &pb.HeartbeatReq{
				Port:        int32(port),
				RunningTask: runningtask.GetRunningTasks(),
			}
			_, err := client.SendHb(ctx, hbreq)
			if err != nil {
				code := DealRPCErr(err)
				log.Error("client.SendHb failed", zap.String("short msg", resp.GetMsg(code)), zap.Error(err))
				timer.Reset(defaultLastFailHearBeatInterval)
				continue
			}
			log.Debug("Send HearBeat Success")
			timer.Reset(defaultHearbeatInterval)

		case <-clentstophb:
			log.Info("Stop Send HearBeat")
			return
		}
	}
}

// DealRPCErr change rpc error to err code
func DealRPCErr(err error) int {
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

// DoStopConn will cancel all running task and close grpc conn
func DoStopConn(mode define.RunMode) {
	if mode == define.Server {
		for id, sch := range Cron.sch {
			sch.running = false
			Cron.Del(id)
			if sch.ctxcancel != nil {
				sch.ctxcancel()
			}
		}
	}

	if mode == define.Client {
		close(clentstophb)
		for _, taskcancel := range runningtask.running {
			taskcancel()
		}
	}
	for _, conn := range cachegRPCConnM.conn {
		conn.Close()
	}
}
