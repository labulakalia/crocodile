package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	midd2 "github.com/labulaka521/crocodile/testcase/grpc/midd"
	"github.com/labulaka521/crocodile/testcase/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type SearchService struct {
	auth *AuthServer
}
type peerKey struct{}

func (s *SearchService) Search(ctx context.Context, req *proto.SearchRequest) (*proto.SearchResponse, error) {
	//p, ok := peer.FromContext(ctx)
	//if ok {
	//	log.Printf("Receive Request From %s Request %s", p.Addr, req.Request)
	//}

	//for i := 0; i < 10; i++  {
	//	if ctx.Err() == context.Canceled {
	//		return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
	//	}
	//	time.Sleep(1 * time.Second)
	//}
	// 进行认证
	if err := s.auth.Check(ctx); err != nil {
		return nil, err
	}
	return &proto.SearchResponse{Response: req.GetRequest() + "Server"}, nil
}

func main() {
	var (
		ca  bool
		err error
		c   credentials.TransportCredentials
	)
	ca = true

	switch ca {
	case true:
		log.Println("start ca auth")
		cert, err := tls.LoadX509KeyPair(
			"conf1/server/server.pem",
			"conf1/server/server.key",
		)
		if err != nil {
			log.Fatalf("credent.LoadX509KeyPair Err: %v", err)
		}
		// 将自签的证书添加至跟证书池
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile("conf1/ca.pem")
		if err != nil {
			log.Fatalf("Read ca Err: %v", err)
		}

		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("certPool.AppendCertsFromPEM Err: %v", err)
		}

		c = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})

	case false:
		c, err = credentials.NewServerTLSFromFile(
			"conf/server.pem",
			"conf/server.key")
		if err != nil {
			log.Fatalf("credentials.NewServerTLSFromFile Err: %v", err)
		}
	}
	grpc.EnableTracing = true
	grpcserver := grpc.NewServer(grpc.Creds(c), grpc_middleware.WithUnaryServerChain(
		midd2.LoggerInterceptor,
		midd2.RecoveryInterceptor,
	),
	)

	proto.RegisterSearchServiceServer(grpcserver, &SearchService{})

	httpserver := GetHTTPServeMux()

	log.Printf("Listen Addr: 9001")

	http.ListenAndServeTLS(
		"127.0.0.1:9001",
		"conf1/server/server.pem",
		"conf1/server/server.key",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcserver.ServeHTTP(w, r)
			} else {
				httpserver.ServeHTTP(w, r)
			}
		}),
	)
	//lis, err := net.Listen("tcp", ":9001")
	//if err != nil {
	//	log.Fatal("Listen Err", err)
	//}

}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("crocodile: test http server\n"))
	})

	return mux
}

type AuthServer struct {
	AppKey    string
	AppSecret string
}

func (a *AuthServer) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Token 获取失败")
	}

	var (
		appkey, appsecret string
	)
	if value, ok := md["app_key"]; ok {
		appkey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appsecret = value[0]
	}

	if appkey != a.GetAppKey() || appsecret != a.GetAppSecret() {
		return status.Errorf(codes.Unauthenticated, "Token 认证失败")
	}
	return nil
}

func (a *AuthServer) GetAppKey() string {
	return "labulaka"
}

func (a *AuthServer) GetAppSecret() string {
	return "crocodile"
}
