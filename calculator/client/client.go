package main

import (
	"context"
	"fmt"
	"grpc-smaarek/calculator/pb"
	"io"
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

	//doUnary(client)
	doServerStreaming(client)

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

func doServerStreaming(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC....")

	req := &pb.PrimeNumberRequest{
		Number: 120,
	}
	resStream, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumber RPC %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumber: %d", msg.GetNumber())
	}
}
