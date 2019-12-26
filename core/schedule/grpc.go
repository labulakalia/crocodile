package schedule

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"os"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/cert"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// rpc conn control
var (
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

func (cg *cachegRPCConn) GetgRPCClientConn(addr string) *grpc.ClientConn {
	cg.Lock()
	conn, exist := cg.conn[addr]
	cg.Unlock()
	if exist {
		return conn
	}
	return nil
}

func (cg *cachegRPCConn) addgRPCClientConn(addr string, conn *grpc.ClientConn) {
	if cg.GetgRPCClientConn(addr) != nil {
		return
	}
	cg.Lock()
	cg.conn[addr] = conn
	cg.Unlock()
}

// Get Grpc Client Conn
func NewgRPCConn(addr string) (*grpc.ClientConn, error) {

	conn := cachegRPCConnM.GetgRPCClientConn(addr)
	if conn != nil {
		return conn, nil
	}
	//cp := x509.NewCertPool()
	//if !cp.AppendCertsFromPEM([]byte(texttls.TlsPemContext)) {
	//	return nil, fmt.Errorf("credentials: failed to append certificates")
	//}
	//c := credentials.NewTLS(&tls.Config{ServerName: texttls.ServerName, RootCAs: cp})

	c, err := credentials.NewClientTLSFromFile(config.CoreConf.Pem.CertFile, cert.ServerName)
	if err != nil {
		log.Error("credentials.NewClientTLSFromFile failed", zap.Error(err))
		return nil, err
	}

	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(c),
		grpc.WithPerRPCCredentials(
			&Auth{SecretToken: config.CoreConf.SecretToken}))

	if err != nil {
		log.Error("grpc.Dial", zap.String("error", err.Error()))
		return nil, err
	}
	cachegRPCConnM.addgRPCClientConn(addr, conn)
	return conn, nil
}

// new gRPC server
func NewgRPCServer(mode define.RunMode) (*grpc.Server, error) {
	c, err := credentials.NewServerTLSFromFile(config.CoreConf.Pem.CertFile, config.CoreConf.Pem.KeyFile)
	if err != nil {
		log.Error("credentials.NewServerTLSFromFile failed", zap.Error(err))
		return nil, err
	}
	auth := Auth{SecretToken: config.CoreConf.SecretToken}
	grpcserver := grpc.NewServer(grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggerInterceptor,
			CheckSecretInterceptor,
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

type Auth struct {
	SecretToken string
}

// implement PerRPCCredentials interface
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

// get rpc conn
func tryGetRpcConn(ctx context.Context, hg *define.HostGroup) (*grpc.ClientConn, error) {
	i := 0
	for i < len(hg.HostsID) {
		i++
		hostid, err := model.RandHostId(hg)
		if err != nil {
			log.Error("model.RandHostId failed", zap.String("error", err.Error()))
			continue
		}

		host, err := model.GetHostById(ctx, hostid)
		if err != nil {
			log.Error("model.GetHostById failed", zap.String("error", err.Error()))
			continue
		}

		conn, err := NewgRPCConn(host.Addr)
		if err != nil {
			log.Error("GetRpcConn failed", zap.String("error", err.Error()))
			continue
		}
		// idle
		if conn.GetState() <= connectivity.Ready {
			return conn, nil
		}
		conn.Close()
	}
	return nil, errors.New("can not get valid grpc conn")
}

// registry client to server
func RegistryClient(version string, port int) error {
	conn, err := NewgRPCConn(config.CoreConf.Client.ServerAddr)
	if err != nil {
		return err
	}
	hbClient := pb.NewHeartbeatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultRpcTimeout)
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
			ctx, cancel := context.WithTimeout(context.Background(), defaultRpcTimeout)
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

func dealRpcErr(err error) int {
	statusErr, ok := status.FromError(err)
	if ok {
		switch statusErr.Code() {
		case codes.DeadlineExceeded:
			return resp.ErrRpcDeadlineExceeded
		case codes.Canceled:
			return resp.ErrRpcCanceled
		case codes.Unauthenticated:
			return resp.ErrRpcUnauthenticated
		case codes.Unavailable:
			return resp.ErrRpcUnavailable
		}
	}
	return resp.ErrRpcUnknow
}
