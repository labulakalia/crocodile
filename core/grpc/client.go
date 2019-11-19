package grpc

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
)

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
	c, err := credentials.NewClientTLSFromFile(Pem, ServerName)
	if err != nil {
		return nil, errors.Wrap(err, "credentials.NewClientTLSFromFile")
	}
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(c))
	cachegRPCConnM.addgRPCClientConn(addr, conn)
	return conn, nil
}
