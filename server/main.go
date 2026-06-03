package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "grpc-chat/grpc-chat/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer

	mu      sync.Mutex
	clients map[pb.ChatService_ChatServer]bool
}

func (s *server) Chat(stream pb.ChatService_ChatServer) error {

	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		s.mu.Lock()

		for client := range s.clients {
			client.Send(msg)
		}

		s.mu.Unlock()
	}
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(
		grpcServer,
		&server{
			clients: make(map[pb.ChatService_ChatServer]bool),
		},
	)

	log.Println("Chat server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
