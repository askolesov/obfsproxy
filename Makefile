BUILD_DIR := build/

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BUILD_DIR) ./...

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
