package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"untether/services/plaid/client"
	"untether/services/plaid/internal"
	pb "untether/services/plaid/proto"

	_ "github.com/lib/pq"
)

func main() {
	// Get server ports from environment variables with defaults
	grpcPortStr := getEnvOrDefault("GRPC_PORT", "50053")
	httpPortStr := getEnvOrDefault("HTTP_PORT", "8082")

	// Convert port strings to integers
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_PORT: %v", err)
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatalf("Invalid HTTP_PORT: %v", err)
	}

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

	// Create database connection pool
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

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
	grpcLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPlaidServiceServer(grpcServer, plaidService)
	reflection.Register(grpcServer)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", httpPort),
		Handler: internal.NewHTTPHandler(plaidService),
	}

	// Start servers in goroutines
	go func() {
		log.Printf("gRPC server listening at %v", grpcLis.Addr())
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		log.Printf("HTTP server listening at %v", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("Servers stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
