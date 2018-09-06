package main

import (
	"context"
	"fmt"
	"grpc-smaarek/calculator/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)

	doUnary(client)
}

func doUnary(c pb.CalculatorServiceClient) {

	fmt.Println("Starting to do a Unary RPC....")

	req := &pb.SumRequest{
		Number1: 15,
		Number2: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Number)
}
