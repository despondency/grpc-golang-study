package main

import (
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (s *server) ComputeAverage(srv calculator.ComputeAverageService_ComputeAverageServer) error {
	var sum float64 = 0
	var sz float64 = 0
	for {
		n, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF received %v", err)
				var av float64 = 0
				if sz > 0 {
					av = sum/sz
				} else {
					av = 0.0
				}
				srv.SendAndClose(&calculator.ComputeAverageResponse{
					Average: av,
				})
				break
			}
		}
		sz++
		sum += float64(n.GetNumber())
	}
	return nil
}

func main() {
	log.Printf("Starting gRPC, client streaming server exercise example")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error occured trying to open tcp connection %v", err)
	}
	s := grpc.NewServer()

	calculator.RegisterComputeAverageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error occured trying to start serving gRPC %v", err)
	}
}
