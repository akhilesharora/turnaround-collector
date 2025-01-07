# Variables
BIN_DIR := bin
BINARY_NAME := turnaround-collector
DOCKER_IMAGE_NAME := $(BINARY_NAME):latest

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/$(BIN_DIR)

# Docker compose files
COMPOSE_FILE := docker-compose.yml
COMPOSE_TEST_FILE := docker-compose.test.yml

.PHONY: all build run test clean docker-build docker-run docker-up docker-down logs help test-integration

all: build test

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@docker-compose down -v

build:
	@echo "Building services..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/collector ./cmd/collector
	@go build -o $(BIN_DIR)/camera ./cmd/camera
	@go build -o $(BIN_DIR)/target ./cmd/target

run: build
	@echo "Running services locally..."
	@./$(BIN_DIR)/collector

test:
	@echo "Running unit tests..."
	@go test -v ./internal/...

test-integration: docker-up
	@echo "Running integration tests..."
	@go test -v ./test/integration/...
	@docker-compose down -v

docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@docker-compose up --remove-orphans --build -d

docker-down:
	@echo "Stopping all services..."
	@docker-compose down -v

docker-logs:
	@docker-compose logs -f

help:
	@echo "Usage: make [TARGET]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Build and run tests"
	@echo "  build            Build all service binaries"
	@echo "  run              Build and run services locally"
	@echo "  test             Run unit tests"
	@echo "  test-integration Run integration tests"
	@echo "  clean            Remove binaries and clean Docker"
	@echo "  docker-build     Build Docker images"
	@echo "  docker-up        Start all services with Docker Compose"
	@echo "  docker-down      Stop all services"
	@echo "  docker-logs             Show Docker Compose logs"
	@echo "  help             Show this help message"