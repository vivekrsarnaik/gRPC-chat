package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "grpc-chat/grpc-chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter message: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.SendMessage(ctx, &pb.ChatMessage{
		Username: username,
		Content:  content,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent!")
}
