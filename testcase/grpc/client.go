package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/labulaka521/crocodile/testcase/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
)

// 自定义认证
type Auth struct {
	AppKey    string
	AppSecret string
}

// 实现 WithPerRPCCredentials PerRPCCredentials 接口
func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return true
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
			"conf1/client/client.pem",
			"conf1/client/client.key",
		)
		if err != nil {
			log.Fatalf("credent.LoadX509KeyPair Err: %v", err)
		}
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile("conf1/ca.pem")
		if err != nil {
			log.Fatalln(err)
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalln(err)
		}
		c = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   "crocodile.me",
			RootCAs:      certPool,
		})

	case false:
		c, err = credentials.NewClientTLSFromFile(
			"conf/server.pem",
			"labulakalia")
		if err != nil {
			log.Fatalf("credentials.NewServerTLSFromFile Err: %v", err)
		}
	}

	auth := Auth{
		AppKey:    "labulaka",
		AppSecret: "crocodile",
	}
	conn, err := grpc.Dial("127.0.0.1:9001",
		grpc.WithTransportCredentials(c),
		grpc.WithPerRPCCredentials(&auth),
	)

	if err != nil {
		log.Fatalf("Dial Err: %v", err)
	}
	//os.Exit(1)

	defer conn.Close()
	client := proto.NewSearchServiceClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(),time.Second * 5)
	//defer cancel()
	resp, err := client.Search(context.Background(), &proto.SearchRequest{Request: "grpc Test"})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("Search Err: %v", err)
	}
	log.Printf("Resp: %v", resp.GetResponse())
}
