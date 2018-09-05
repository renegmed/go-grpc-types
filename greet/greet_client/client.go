package main

import (
	"fmt"
	"grpc-smaarek/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client: %f\n", client)
}
