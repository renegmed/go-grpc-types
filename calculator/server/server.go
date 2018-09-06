package main

import (
	"context"
	"fmt"
	"grpc-smaarek/calculator/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.ResultResponse, error) {
	number1 := req.GetNumber1()
	number2 := req.GetNumber2()

	fmt.Printf("Number 1: %d  Number 2: %d", number1, number2)
	res := &pb.ResultResponse{
		Number: number1 + number2,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello Server")

	listener, err := net.Listen("tcp", "0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
