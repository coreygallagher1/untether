package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"untether/services/plaid/client"
	"untether/services/user/internal"
	pb "untether/services/user/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	// Initialize database connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/untether?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize Plaid client
	plaidClient := client.NewPlaidClient("client_id", "secret", "sandbox")

	// Initialize user service
	userService := internal.NewUserService(db, plaidClient)

	// Create gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)

	// Register reflection service for gRPCurl
	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
