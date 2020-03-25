package main

import (
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) DecomposePrime(req *calculator.PrimeNumberDecompositionRequest, srv calculator.PrimeNumberDecompositionService_DecomposePrimeServer) error {
	number := req.Number
	var i int64 = 2
	for number > 1 {
		if number%i == 0 {
			err := srv.Send(&calculator.PrimeNumberDecompositionResponse{
				PrimeFactor: i,
			})
			if err != nil {
				return err
			}
			number /= i
		} else {
			i++
		}
	}
	return nil
}

func main() {
	log.Printf("Starting gRPC, server streaming server exercise example")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error occured trying to open tcp connection %v", err)
	}
	s := grpc.NewServer()

	calculator.RegisterPrimeNumberDecompositionServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error occured trying to start serving gRPC %v", err)
	}
}
