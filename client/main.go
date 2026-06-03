package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

	stream, err := client.Chat(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	go func() {

		for {

			msg, err := stream.Recv()

			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("\n%s: %s\n",
				msg.Username,
				msg.Content,
			)
		}
	}()

	for {

		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		err := stream.Send(&pb.ChatMessage{
			Username: username,
			Content:  text,
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
