.PHONY: all build test proto clean docker-up docker-down

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOCMD) mod tidy

# Binary names
BINARY_NAME=untether

# Proto commands
PROTOC=protoc
PROTO_DIR=proto
GO_OUT_DIR=internal/proto

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/untether

test:
	$(GOTEST) -v ./...

proto:
	mkdir -p $(GO_OUT_DIR)
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)
	rm -rf $(GO_OUT_DIR)

deps:
	$(GOTIDY)

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Development setup
dev-setup: deps proto docker-up

# Development cleanup
dev-cleanup: docker-down clean 