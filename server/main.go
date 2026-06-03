package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-chat/grpc-chat/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) SendMessage(
	ctx context.Context,
	msg *pb.ChatMessage,
) (*pb.Empty, error) {

	fmt.Printf(
		"[%s] %s: %s\n",
		time.Now().Format("15:04:05"),
		msg.Username,
		msg.Content,
	)

	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(grpcServer, &server{})

	fmt.Println("Server running on port 50051...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
