.PHONY: all build test proto clean docker-up docker-down migrate-up migrate-down

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

# Migration commands
MIGRATE=migrate
MIGRATE_PATH=migrations
MIGRATE_DATABASE_URL=postgres://untether:untether@localhost:5432/untether?sslmode=disable

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/untether

test:
	$(GOTEST) -v ./...

proto:
	mkdir -p $(GO_OUT_DIR)
	$(PROTOC) -I$(PROTO_DIR) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
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

# Migration commands
migrate-up:
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE_URL)" down

# Development setup
dev-setup: deps proto docker-up migrate-up

# Development cleanup
dev-cleanup: docker-down clean 