package schedule

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/cert"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
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
	clientstophb chan struct{}
)

func init() {
	cachegRPCConnM = &cachegRPCConn{
		conn: make(map[string]*grpc.ClientConn),
	}
	clientstophb = make(chan struct{}, 1)
}

type cachegRPCConn struct {
	sync.RWMutex
	conn map[string]*grpc.ClientConn
}

// getgRPCClientConn return conn or nil
func (cg *cachegRPCConn) getgRPCClientConn(addr string) *grpc.ClientConn {
	cg.RLock()
	conn, exist := cg.conn[addr]
	cg.RUnlock()
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
	var (
		c   credentials.TransportCredentials
		err error
	)

	dialoptions := []grpc.DialOption{
		grpc.WithPerRPCCredentials(
			&Auth{SecretToken: config.CoreConf.SecretToken},
		),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(16 * 1024 * 1024)), // 16M
		grpc.WithBlock(),
		// grpc.WithKeepaliveParams(keepalive.ClientParameters{}),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{MaxDelay: time.Second * 2}, MinConnectTimeout: time.Second * 2}),
	}

	if config.CoreConf.Cert.Enable {
		c, err = credentials.NewClientTLSFromFile(config.CoreConf.Cert.CertFile, cert.ServerName)
		if err != nil {
			log.Error("credentials.NewClientTLSFromFile failed", zap.Error(err))
			return nil, err
		}
		dialoptions = append(dialoptions, grpc.WithTransportCredentials(c))
	} else {
		dialoptions = append(dialoptions, grpc.WithInsecure())
	}

	rpcctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
	defer cancel()
	//
	conn, err = grpc.DialContext(rpcctx, addr, dialoptions...)
	if err != nil {
		return nil, err
	}
	cachegRPCConnM.addgRPCClientConn(addr, conn)
	return conn, nil
}

// NewgRPCServer new gRPC server
func NewgRPCServer(mode define.RunMode) (*grpc.Server, error) {
	serveroptions := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			middleware.RecoveryInterceptor,
			middleware.LoggerInterceptor,
			middleware.CheckSecretInterceptor,
		),
		grpc.MaxRecvMsgSize(16 * 1024 * 1024),
		// grpc.KeepaliveParams(keepalive.ServerParameters{
		// 	MaxConnectionIdle: 5 * time.Minute, // <--- This fixes it!
		// }),
	}
	if config.CoreConf.Cert.Enable {
		c, err := credentials.NewServerTLSFromFile(config.CoreConf.Cert.CertFile, config.CoreConf.Cert.KeyFile)
		if err != nil {
			log.Error("credentials.NewServerTLSFromFile failed", zap.Error(err))
			return nil, err
		}
		serveroptions = append(serveroptions, grpc.Creds(c))

	}
	auth := Auth{SecretToken: config.CoreConf.SecretToken}
	grpcserver := grpc.NewServer(serveroptions...)
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
	return config.CoreConf.Cert.Enable
}

// try Get grpc client conn
// Scheduling Algorithm
// - random
// - LeastTask
// - Weight
// - roundRobin
// get rpc conn
func tryGetRCCConn(ctx context.Context, next Next) (*grpc.ClientConn, error) {
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
			log.Error("GetRpcConn failed", zap.Error(err))
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
func RegistryClient(version string, port int) {
	rand.Seed(time.Now().UnixNano())
	var (
		// cancel   context.CancelFunc
		// ctx      context.Context
		lastaddr string
	)

	for {
		// ctx, cancel = context.WithTimeout(context.Background(), defaultRPCTimeout)
		addrs := config.CoreConf.Client.ServerAddrs
		if len(addrs) == 0 {
			log.Error("server addrs is empty")
			// cancel()
			return
		}
		// do not get last addr
		for {
			getaddr := addrs[rand.Int()%len(addrs)]
			// do not get failed addr
			if getaddr != lastaddr || len(addrs) == 1 {
				lastaddr = getaddr
				break
			}
		}

		conn, err := getgRPCConn(context.Background(), lastaddr)
		if err != nil {
			log.Error("getgRPCConn failed", zap.Error(err))
			time.Sleep(time.Second)

			continue
		}
		hbClient := pb.NewHeartbeatClient(conn)
		hostname, _ := os.Hostname()
		regHost := pb.RegistryReq{
			Port:      int32(port),
			Hostname:  hostname,
			Version:   version,
			Hostgroup: config.CoreConf.Client.HostGroup,
			Weight:    int32(config.CoreConf.Client.Weight),
			Remark:    config.CoreConf.Client.Remark,
		}
		_, err = hbClient.RegistryHost(context.Background(), &regHost)
		if err != nil {
			log.Error("registry client failed", zap.Error(err))
			time.Sleep(time.Second)
			continue

		}

		log.Info("host registry success", zap.String("server", lastaddr))
		timer := time.NewTimer(defaultHearbeatInterval)
		cannotconn := 0
		for {
			select {
			case <-timer.C:
				ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
				hbreq := &pb.HeartbeatReq{
					Port:        int32(port),
					RunningTask: runningtask.GetRunningTasks(),
				}

				_, err := hbClient.SendHb(ctx, hbreq)
				if err != nil {
					cancel()
					err := DealRPCErr(err)
					if err.Error() == resp.GetMsgErr(resp.ErrRPCUnavailable).Error() {
						if cannotconn > 1 {
							// 断开超过两次重新在别的调度中心注册
							if len(config.CoreConf.Client.ServerAddrs) >= 2 {
								log.Debug("can not conn server,change other server")
								conn.Close()
								goto Next
							}
						} else {
							cannotconn++
						}
					}
					log.Error("client.SendHb failed", zap.Error(err))
					timer.Reset(defaultLastFailHearBeatInterval)
					continue
				}
				cannotconn = 0
				cancel()
				log.Debug("send hearbeat success", zap.String("server", lastaddr))
				timer.Reset(defaultHearbeatInterval)
			case <-clientstophb:
				log.Info("Stop Send HearBeat")
				timer.Stop()
				return
			}
		}
	Next:
		time.Sleep(time.Second)
	}
}

// send client will send hearbt to server, let scheduler center know it is alive
// func sendhb(client pb.HeartbeatClient, port int) error {
// 	log.Info("start send hearbeat to server")
// 	timer := time.NewTimer(time.Millisecond)

// 	cannotconn := 0
// 	for {
// 		select {
// 		case <-timer.C:
// 			ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
// 			hbreq := &pb.HeartbeatReq{
// 				Port:        int32(port),
// 				RunningTask: runningtask.GetRunningTasks(),
// 			}
// 			_, err := client.SendHb(ctx, hbreq)
// 			if err != nil {
// 				cancel()
// 				err := DealRPCErr(err)

// 				if err.Error() == resp.GetMsgErr(resp.ErrRPCUnavailable).Error() {
// 					if cannotconn > 2 {
// 						// 断开超过两次
// 						// 重新在别的调度中心注册
// 						if len(config.CoreConf.Client.ServerAddrs) >= 2 {
// 							log.Debug("can not conn server,change other server")
// 							return err
// 						}
// 					} else {
// 						cannotconn++
// 					}
// 				}
// 				log.Error("client.SendHb failed", zap.Error(err))
// 				timer.Reset(defaultLastFailHearBeatInterval)
// 				continue
// 			}
// 			cannotconn = 0
// 			cancel()
// 			log.Debug("Send HearBeat Success")
// 			timer.Reset(defaultHearbeatInterval)

// 		case <-clentstophb:
// 			log.Info("Stop Send HearBeat")
// 			timer.Stop()
// 			return errStopHearBeat
// 		}
// 	}
// }

// DealRPCErr change rpc error to err code
func DealRPCErr(err error) error {
	statusErr, ok := status.FromError(err)
	if ok {
		switch statusErr.Code() {
		case codes.DeadlineExceeded:
			return resp.GetMsgErr(resp.ErrCtxDeadlineExceeded)
		case codes.Canceled:
			return resp.GetMsgErr(resp.ErrCtxCanceled)
		case codes.Unauthenticated:
			return resp.GetMsgErr(resp.ErrRPCUnauthenticated)
		case codes.Unavailable:
			return resp.GetMsgErr(resp.ErrRPCUnavailable)
		}
	}
	return err
}

// DoStopConn will cancel all running task and close grpc conn
func DoStopConn(mode define.RunMode) {
	if mode == define.Server {
		for id, t := range Cron2.ts {
			// sch.running = false
			Cron2.deletetask(id)
			if t.ctxcancel != nil {
				t.ctxcancel()
			}
		}
	}

	if mode == define.Client {
		close(clientstophb)
		for _, taskcancel := range runningtask.running {
			taskcancel()
		}
	}
	for _, conn := range cachegRPCConnM.conn {
		conn.Close()
	}
}
