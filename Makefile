APP_NAME=ggpoker
BUILD_DIR=build

dev:
	@echo "Running in development mode..."
	clear
	@go run main.go

build:
	@echo "Building production binary..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run:
	@echo "Running production binary..."
	@./$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running the testing..."
	@go test -v ./...

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)