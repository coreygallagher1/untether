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

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"untether/services/transaction/internal"
	pb "untether/services/transaction/proto"
)

func main() {
	// Get server ports from environment variables with defaults
	grpcPortStr := getEnvOrDefault("GRPC_PORT", "50051")
	httpPortStr := getEnvOrDefault("HTTP_PORT", "8083")

	// Convert port strings to integers
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_PORT: %v", err)
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatalf("Invalid HTTP_PORT: %v", err)
	}

	// Initialize database connection with connection pooling
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create a TCP listener for gRPC
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register the transaction calculator service
	calculator := internal.NewTransactionCalculator(db)
	pb.RegisterTransactionCalculatorServer(grpcServer, calculator)

	// Register reflection service for development
	reflection.Register(grpcServer)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: internal.NewHTTPHandler(calculator),
	}

	// Create a channel to listen for errors coming from the servers
	serverErrors := make(chan error, 2)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("gRPC server listening at %v", grpcLis.Addr())
		serverErrors <- grpcServer.Serve(grpcLis)
	}()

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP server listening at :%d", httpPort)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Create a channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		log.Printf("Error starting server: %v", err)
	case sig := <-shutdown:
		log.Printf("Got signal: %v", sig)
	}

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown of the HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down HTTP server: %v", err)
	}

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initDB() (*sql.DB, error) {
	// Get database connection details from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "postgres")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "untether")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "untether")
	dbName := getEnvOrDefault("DB_NAME", "untether")

	// Construct database connection string
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize database connection
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
