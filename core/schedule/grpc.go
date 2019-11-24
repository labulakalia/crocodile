package schedule

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	pb "github.com/labulaka521/crocodile/core/proto"
	texttls "github.com/labulaka521/crocodile/core/tls"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
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
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(texttls.TlsPemContext)) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	c := credentials.NewTLS(&tls.Config{ServerName: texttls.ServerName, RootCAs: cp})

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(c))

	if err != nil {
		log.Error("grpc.Dial", zap.String("error", err.Error()))
		return nil, err
	}
	cachegRPCConnM.addgRPCClientConn(addr, conn)
	return conn, nil
}

// new gRPC server
func NewgRPCServer(mode define.RunMode) (*grpc.Server, error) {
	cert, err := tls.X509KeyPair([]byte(texttls.TlsPemContext), []byte(texttls.TlskeyContent))
	if err != nil {
		return nil, errors.Wrap(err, "tls.X509KeyPair")
	}
	c := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})

	grpcserver := grpc.NewServer(grpc.Creds(c))
	switch mode {
	case define.Server:
		pb.RegisterHeartbeatServer(grpcserver, &HeartbeatService{})
	case define.Client:
		pb.RegisterTaskServer(grpcserver, &TaskService{})
	}
	return grpcserver, nil
}
