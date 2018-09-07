package main

import (
	"context"
	"fmt"
	"grpc-smaarek/calculator/pb"
	"io"
	"log"
	"time"

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
	//doServerStreaming(client)

	doClientStreaming(client)

}

func doClientStreaming(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC....")
	requests := []*pb.NumberRequest{
		&pb.NumberRequest{
			Number: 1,
		},
		&pb.NumberRequest{
			Number: 2,
		},
		&pb.NumberRequest{
			Number: 3,
		},
		&pb.NumberRequest{
			Number: 4,
		},
	}
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Average: %v", err)
	}
	// iterate slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Average: %v", err)
	}

	fmt.Printf("Computed Average Response: %v\n", res)
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
