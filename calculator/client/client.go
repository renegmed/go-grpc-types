package main

import (
	"context"
	"fmt"
	"grpc-smaarek/calculator/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

	// doClientStreaming(client)

	// doBiDiStreaming(client)

	//doErrorUnary(client)
	doErrorUnary(client)

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

func doBiDiStreaming(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bi-Directional Streaming RPC....")
	requests := []*pb.MaximumNumberRequest{
		&pb.MaximumNumberRequest{
			Number: 1,
		},
		&pb.MaximumNumberRequest{
			Number: 5,
		},
		&pb.MaximumNumberRequest{
			Number: 3,
		},
		&pb.MaximumNumberRequest{
			Number: 6,
		},
		&pb.MaximumNumberRequest{
			Number: 2,
		},
		&pb.MaximumNumberRequest{
			Number: 20,
		},
		&pb.MaximumNumberRequest{
			Number: 15,
		},
		&pb.MaximumNumberRequest{
			Number: 8,
		},
	}

	// create a stream by invoking the client
	stream, err := c.MaximumNumber(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})

	// send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Microsecond)
		}
		stream.CloseSend()
	}()
	// receive a bunch of messages from the server (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received Maximum Number: %v\n", res.GetNumber())
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doErrorUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)

}

func doErrorCall(c pb.CalculatorServiceClient, n int64) {

	res, err := c.SquareRoot(context.Background(),
		&pb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
			return
		}
	}
	fmt.Printf("Result of square root of  %v: %v\n", n, res.GetNumberRoot())

}

// func doErrorUnary(c pb.CalculatorServiceClient) {
// 	fmt.Println("Starting to do a Bi-Directional Streaming RPC....")
// 	number := int32(10)
// 	// correct call
// 	res, err := c.SquareRoot(context.Background(),
// 		&pb.SquareRootRequest{Number: number})
// 	if err != nil {
// 		respErr, ok := status.FromError(err)
// 		if ok {
// 			// actual error from gRPC (user error)
// 			fmt.Println(respErr.Message())
// 			fmt.Println(respErr.Code())
// 			if respErr.Code() == codes.InvalidArgument {
// 				fmt.Println("We probably sent a negative number!\n")
// 			}
// 		} else {
// 			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
// 		}
// 	}

// 	// error call
// }
