package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"untether/services/plaid/client"
	"untether/services/user/internal"
	"untether/services/user/internal/middleware"
	pb "untether/services/user/proto"
)

func main() {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Get server ports from environment variables
	grpcPortStr := os.Getenv("GRPC_PORT")
	httpPortStr := os.Getenv("HTTP_PORT")

	// Convert port strings to integers
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_PORT: %v", err)
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatalf("Invalid HTTP_PORT: %v", err)
	}

	// Construct database connection string
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize database connection
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize Plaid client with environment variables
	plaidClientID := os.Getenv("PLAID_CLIENT_ID")
	plaidSecret := os.Getenv("PLAID_CLIENT_SECRET")
	plaidEnv := os.Getenv("PLAID_ENVIRONMENT")
	plaidClient := client.NewPlaidClient(plaidClientID, plaidSecret, plaidEnv)

	// Initialize user service
	userService := internal.NewUserService(db, plaidClient)

	// Create gRPC server
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	// Register reflection service for gRPCurl
	reflection.Register(grpcServer)

	// Create HTTP handler with authentication middleware
	httpHandler := internal.NewHTTPHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Create HTTP server with middleware
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: authMiddleware.Authenticate(httpHandler),
	}

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("gRPC server listening at %v", grpcLis.Addr())
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	log.Printf("HTTP server listening at :%d", httpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
