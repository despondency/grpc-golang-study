package main

import (
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculator.SumRequest) (*calculator.SumResponse, error) {
	firstArg := req.GetFirstArgument()
	secondArg := req.GetSecondArgument()
	response := &calculator.SumResponse{
		Result: int64(firstArg + secondArg),
	}
	return response, nil
}

func main() {
	log.Printf("Starting gRPC, unary server exercise example")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error occured trying to open tcp connection %v", err)
	}
	s := grpc.NewServer()

	calculator.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error occured trying to start serving gRPC %v" ,err)
	}
}
