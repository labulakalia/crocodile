package main

import (
	"context"
	"github.com/labulaka521/crocodile/testcase/grpc/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

func printLists(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("resp: %+v", resp.Pt)
	}
	return nil
}

func printRecord(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("server resp: %+v", resp.Pt)
	return nil
}

func printRoute(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 10; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %+v", resp.Pt)
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(":9002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial Err: %v", err)
	}

	defer conn.Close()

	client := proto.NewStreamServiceClient(conn)

	err = printLists(client, &proto.StreamRequest{
		Pt: &proto.StreamPoint{
			Name:  "grpc Stream Client List",
			Value: 1024,
		},
	})
	if err != nil {
		log.Fatalf("printLists Err: %v", err)
	}

	err = printRecord(client, &proto.StreamRequest{
		Pt: &proto.StreamPoint{
			Name:  "grpc Stream Client Record",
			Value: 1024,
		},
	})
	if err != nil {
		log.Fatalf("printRecords Err: %v", err)
	}

	err = printRoute(client, &proto.StreamRequest{
		Pt: &proto.StreamPoint{
			Name:  "grpc Stream Client Route",
			Value: 1025,
		},
	})
	if err != nil {
		log.Fatalf("printRoutes Err: %v", err)
	}
}
