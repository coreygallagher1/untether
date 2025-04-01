FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o user-service ./services/user/cmd

# Expose the port
EXPOSE 50051

# Run the application
CMD ["./user-service"] 