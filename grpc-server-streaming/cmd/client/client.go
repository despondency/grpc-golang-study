package main

import (
	"context"
	"fmt"
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"io"
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
	log.Printf("Starting gRPC, server streaming client exercise example")
	
	clientConn, err := tryInitialize("localhost:50051")
	if err != nil {
		log.Fatalf("Could not initialize client connection %v", err)
	}
	defer clientConn.Close()
	c := calculator.NewPrimeNumberDecompositionServiceClient(clientConn)

	n := int64(120)

	primeDecompositionRequest := &calculator.PrimeNumberDecompositionRequest{
		Number: n,
	}

	response, err := c.DecomposePrime(context.Background(), primeDecompositionRequest)
	if err != nil {
		log.Fatalf("Error trying to gRPC call DecomposePrime %v", err)
	}
	for {
		denominator, err := response.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if err != nil {
			log.Fatalf("Error reading stream %v", err)
		}
		fmt.Printf("Denom of %v is %v\n", n, denominator.GetPrimeFactor())
	}
}