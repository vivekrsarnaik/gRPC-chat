# gRPC ChatHub

A real-time chat application built using Go, gRPC, Protocol Buffers, and Docker. The application supports direct and broadcast messaging, enabling efficient communication between multiple clients through gRPC streaming and concurrent processing.

## Key Highlights
- Built a distributed chat system using Go and gRPC.
- Implemented real-time messaging with streaming RPCs.
- Used Protocol Buffers for efficient message serialization.
- Leveraged Go goroutines and channels for concurrent client handling.
- Containerized the application using Docker for easy deployment and scalability.

## Tech Stack
Go • gRPC • Protocol Buffers • Docker

## Features
- Real-time messaging
- One-to-one messaging
- Broadcast messaging
- Concurrent client support
- Dockerized deployment

## Getting Started

```bash
git clone https://github.com/yourusername/grpc-chathub.git
cd grpc-chathub
go mod tidy
protoc --go_out=. --go-grpc_out=. proto/chat.proto
go run server/server.go
go run client/client.go
