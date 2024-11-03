BUILD_DIR := build/

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BUILD_DIR) ./...

.PHONY: docker
docker:
	docker build -f docker/Dockerfile -t obfsproxy .

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: pre-commit
pre-commit:
	make lint
	make test
	make build
	make docker
	cd integration && make test
