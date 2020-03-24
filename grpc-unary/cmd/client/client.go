package main

import (
	"context"
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"log"
	"time"
)

func tryInitialize(target string) (*grpc.ClientConn, error) {
	log.Printf("Starting retry logic")
	maxRetry := 0
	for {
		clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			if maxRetry == 5 {
				log.Fatalf("Max retries exceeded, stopping retry mechanism %v", err)
			}
			log.Printf("Could not connect to the given gRPC endpoint: %v, will retry in 5sec", err)
			time.Sleep(time.Second * 5)
			maxRetry++
			continue
		}
		log.Printf("Connected successfully to %v", clientConn)
		return clientConn, nil
	}
}

func main() {
	log.Printf("Starting gRPC, unary client exercise example")
	
	clientConn, err := tryInitialize("localhost:50051")
	if err != nil {
		log.Fatalf("Could not initialize client connection %v", err)
	}
	defer clientConn.Close()
	c := calculator.NewSumServiceClient(clientConn)

	sumRequest := &calculator.SumRequest{
		FirstArgument:        1337,
		SecondArgument:       63,
	}

	response, err := c.Sum(context.Background(), sumRequest)
	if err != nil {
		log.Fatalf("Error occurred trying to call gRPC Sum %v", err)
	}

	log.Printf("Response is %v", response)
}