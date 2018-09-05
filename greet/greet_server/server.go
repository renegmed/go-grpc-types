package main

import (
	"fmt"
	"grpc-smaarek/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello World")

	listener, err := net.Listen("tcp", "0.0.0:50051") // 50051 default port for grpc
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
