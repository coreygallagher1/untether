package main

import (
	"log"
	"net"

	"untether/services/transaction/internal"
	pb "untether/services/transaction/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create a TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create and register the transaction calculator service
	calculator := internal.NewTransactionCalculator()
	pb.RegisterTransactionCalculatorServer(s, calculator)

	// Register reflection service for development
	reflection.Register(s)

	log.Printf("Transaction service starting on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
