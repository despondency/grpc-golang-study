package main

import (
	"github.com/despondency/grpc-golang-study/pkg/generated/api/calculator"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"net"
)

type server struct{}

func (s *server) FindMaximum(srv calculator.FindMaximumService_FindMaximumServer) error {
	var mx int64 = math.MinInt64
	for {
		req, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		currNum := req.GetNumber()
		log.Printf("Received number %v", currNum)
		if currNum > mx {
			mx = currNum
		}
		err = srv.Send(&calculator.FindMaximumResponse{MaxNumber: mx})
		if err != nil {
			return err
		}
	}
}

func main() {
	log.Printf("Starting gRPC, bi directional streaming server exercise example")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error occured trying to open tcp connection %v", err)
	}
	s := grpc.NewServer()

	calculator.RegisterFindMaximumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error occured trying to start serving gRPC %v", err)
	}
}
