package main

import (
	"context"
	"log"
	"net"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type Note struct {
	Title       string
	Description string
}

type noteService struct {
	rabbitMQConn *amqp.Connection
}

func (s *noteService) AddNote(ctx context.Context, req *Note) (*AddNoteResponse, error) {
	note := req

	// Convert note to RabbitMQ message format
	message := &RabbitMQMessage{
		// Populate RabbitMQ message fields with note data
	}

	// Publish message to RabbitMQ exchange
	err := s.publishMessage(message)
	if err != nil {
		return nil, err
	}

	return &AddNoteResponse{
		Success: true,
	}, nil
}

func (s *noteService) publishMessage(message *RabbitMQMessage) error {
	channel, err := s.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declare and publish to the exchange
	// Use appropriate RabbitMQ exchange, routing key, and message encoding

	return nil
}

func main() {
	rabbitMQConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Initialize noteService and set the RabbitMQ connection
	noteSvc := &noteService{
		rabbitMQConn: rabbitMQConn,
	}

	// Register noteService with the gRPC server
	RegisterNoteServiceServer(grpcServer, noteSvc)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port 50051")
	log.Fatal(grpcServer.Serve(listener))
}
