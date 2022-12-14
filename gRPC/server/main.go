package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/kqns91/til/gRPC/pb/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	logger := log.Default()

	server := grpc.NewServer()
	pb.RegisterExampleServiceServer(server, newExampleServiceServer())

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Printf("faield to listen: %v", err)
	}

	logger.Println("listening...")

	if err := server.Serve(lis); err != nil {
		logger.Printf("failed to serve: %v", err)
	}
}

func newExampleServiceServer() pb.ExampleServiceServer {
	return &exampleServiceServer{
		UnimplementedExampleServiceServer: pb.UnimplementedExampleServiceServer{},
	}
}

type exampleServiceServer struct {
	pb.UnimplementedExampleServiceServer
}

func (s *exampleServiceServer) ClientStream(stream pb.ExampleService_ClientStreamServer) error {
	name := ""
	data := make([]byte, 0, 1024*1024*100)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to receve: %w", err)
		}

		name = req.GetName()
		data = append(data, req.GetData()...)
	}

	fmt.Printf("name: %v\n", name)
	fmt.Printf("string(data): %v\n", string(data))

	stream.SendAndClose(&pb.ExampleResponse{
		Message: "finish",
	})

	return nil
}
