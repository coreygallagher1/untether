package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"untether/services/plaid/client"
	"untether/services/plaid/internal"
	pb "untether/services/plaid/proto"

	_ "github.com/lib/pq"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

func main() {
	flag.Parse()

	// Initialize database connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "untether"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "untether"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "untether"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Initialize Plaid client
	plaidClientID := os.Getenv("PLAID_CLIENT_ID")
	if plaidClientID == "" {
		log.Fatal("PLAID_CLIENT_ID environment variable is required")
	}

	plaidClientSecret := os.Getenv("PLAID_CLIENT_SECRET")
	if plaidClientSecret == "" {
		log.Fatal("PLAID_CLIENT_SECRET environment variable is required")
	}

	plaidEnvironment := os.Getenv("PLAID_ENVIRONMENT")
	if plaidEnvironment == "" {
		plaidEnvironment = "sandbox"
	}

	plaidClient := client.NewPlaidClient(plaidClientID, plaidClientSecret, plaidEnvironment)

	// Initialize plaid service
	plaidService := internal.NewPlaidService(plaidClient, db)

	// Create gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlaidServiceServer(s, plaidService)

	// Register reflection service for gRPCurl
	reflection.Register(s)

	log.Printf("Plaid service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
