BINARY_NAME = task-queue
BUILD_DIR = bin

WORKERS ?= 4
QUEUE_SIZE ?= 64

.PHONY: all build start run clean test

all: run

build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/task-queue/main.go

start:
	@echo "Starting application with WORKERS=$(WORKERS), QUEUE_SIZE=$(QUEUE_SIZE)..."
	WORKERS=$(WORKERS) QUEUE_SIZE=$(QUEUE_SIZE) $(BUILD_DIR)/$(BINARY_NAME)

run: build start

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@echo "Running tests..."
	@go test -v ./...
