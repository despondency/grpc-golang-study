package main

import (
	"context"
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"sync"
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
	c := calculator.NewFindMaximumServiceClient(clientConn)

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error occurred calling gRPC find maximum %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		for {
			max, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("EOF received, stopping")
					wg.Done()
					break
				} else {
					log.Fatalf("Error receiving from gRPC server find maximum")
				}
			}
			log.Printf("Current max is %v", max.GetMaxNumber())
		}
	}(&wg)
	go func() {
		for i:=0;i<100;i++ {
			rand.Seed(time.Now().UTC().UnixNano())
			err := stream.Send(&calculator.FindMaximumRequest{
				Number: rand.Int63(),
			})
			if err != nil {
				log.Fatalf("Trouble sending FindMaximumRequest with %v", err)
			}
			time.Sleep(time.Second * 5)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error closing client connection %v", err)
		}
	}()
	wg.Wait()
}