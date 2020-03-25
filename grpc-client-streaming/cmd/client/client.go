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
	log.Printf("Starting gRPC, client streaming client exercise example")
	
	clientConn, err := tryInitialize("localhost:50051")
	if err != nil {
		log.Fatalf("Could not initialize client connection %v", err)
	}
	defer clientConn.Close()
	c := calculator.NewComputeAverageServiceClient(clientConn)

	client ,err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error occurred trying to call gRPC Compute Average %v", err)
	}
	for i := 1; i <= 100;i++ {
		client.Send(&calculator.ComputeAverageRequest{Number:int32(i)})
	}
	resp, err := client.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error trying to close client stream and recv %v", err)
	}
	log.Printf("Average is %v", resp.GetAverage())


}