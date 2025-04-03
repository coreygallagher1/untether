FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY ../../go.mod ../../go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ../../ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o transaction-service ./services/transaction/cmd

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/transaction-service .

# Expose both gRPC and HTTP ports
EXPOSE 50054 8083

# Run the service
CMD ["./transaction-service"] 