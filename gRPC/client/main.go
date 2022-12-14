package main

import (
	"bufio"
	"context"
	"log"
	"os"

	pb "github.com/kqns91/til/gRPC/pb/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.Default()

	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Printf("failed to dial: %v\n", conn)

		return
	}
	defer conn.Close()

	client := pb.NewExampleServiceClient(conn)

	stream, err := client.ClientStream(context.Background())
	if err != nil {
		logger.Printf("failed to client stream: %v", err)

		return
	}

	file, err := os.Open("./test.txt")
	if err != nil {
		logger.Printf("failed to open: %v", err)

		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		err = stream.Send(&pb.ExampleRequest{
			Name: "test",
			Data: scanner.Bytes(),
		})
		if err != nil {
			logger.Printf("failed to send: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Printf("scanner error: %v", err)

		return
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.Printf("faield to close and receve: %v", err)
	}

	logger.Println(res)
}
