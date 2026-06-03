package main

import (
	"database/sql"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	pb "grpc-chat/grpc-chat/proto"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer

	mu      sync.Mutex
	clients map[string]pb.ChatService_ChatServer
	db      *sql.DB
}

func (s *server) Chat(stream pb.ChatService_ChatServer) error {

	var username string

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if username == "" {
			username = strings.ToLower(msg.Username)

			s.mu.Lock()
			s.clients[username] = stream
			s.mu.Unlock()

			rows, err := s.db.Query(`
		SELECT username, content
		FROM messages
		ORDER BY id DESC
		LIMIT 20
	`)

			if err == nil {

				var history []pb.ChatMessage

				for rows.Next() {

					var u string
					var c string

					rows.Scan(&u, &c)

					history = append(history, pb.ChatMessage{
						Username: u,
						Content:  c,
						History:  true,
					})
				}

				for i := len(history) - 1; i >= 0; i-- {
					_ = stream.Send(&history[i])
				}

				rows.Close()
			}
		}

		_, err = s.db.Exec(
			"INSERT INTO messages(username, content) VALUES($1, $2)",
			msg.Username,
			msg.Content,
		)

		if err != nil {
			log.Println("DB insert failed:", err)
		}

		if strings.HasPrefix(msg.Content, "@") {

			parts := strings.SplitN(msg.Content, " ", 2)

			if len(parts) < 2 {
				continue
			}

			target := strings.TrimPrefix(parts[0], "@")
			target = strings.ToLower(target)

			message := parts[1]

			s.mu.Lock()
			targetStream, exists := s.clients[target]
			s.mu.Unlock()

			if exists {
				_ = targetStream.Send(&pb.ChatMessage{
					Username: msg.Username + " (DM)",
					Content:  message,
				})
			}

			continue
		}

		s.mu.Lock()

		for _, client := range s.clients {
			client.Send(msg)
		}

		s.mu.Unlock()
	}
}

func main() {

	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=vivek password=password123 dbname=chatdb sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(
		grpcServer,
		&server{
			clients: make(map[string]pb.ChatService_ChatServer),
			db:      db,
		},
	)
	log.Println("Chat server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
