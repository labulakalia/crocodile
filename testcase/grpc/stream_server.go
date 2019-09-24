package main

import (
	"github.com/labulaka521/crocodile/testcase/grpc/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct{}

// 服务端流
func (s *StreamService) List(r *proto.StreamRequest, stream proto.StreamService_ListServer) error {
	for n := 0; n < 10; n++ {
		err := stream.Send(&proto.StreamResponse{
			Pt: &proto.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 客户端流
func (s *StreamService) Record(stream proto.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.StreamResponse{
				Pt: &proto.StreamPoint{
					Name:  "grpc Server Stream Record Close",
					Value: 101,
				},
			})
		}
		if err != nil {
			return err
		}
		log.Printf("Stream.Recv: %+v", r.Pt)
	}
	return nil
}

// 双向流
func (s *StreamService) Route(stream proto.StreamService_RouteServer) error {
	var n int32 = 0

	for {
		log.Printf("Start Send %d", n)
		err := stream.Send(&proto.StreamResponse{
			Pt: &proto.StreamPoint{
				Name:  "grpc Stream Route",
				Value: n,
			},
		})
		if err != nil {
			return err
		}
		log.Printf("Start Recv %d", n)
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("resp: %+v", resp.Pt)
	}
	return nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("Listen Err: %v", err)
	}
	log.Fatalln(server.Serve(lis))
}
