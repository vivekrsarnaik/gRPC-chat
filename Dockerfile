FROM golang:1.26

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o chat-server ./server

EXPOSE 50051

CMD ["./chat-server"]