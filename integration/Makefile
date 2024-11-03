.PHONY: all
all: test

.PHONY: test
test:
	# Start the containers
	docker compose up -d --build
	sleep 3

	# Run the tests
	go test -v ./integration_test.go

	# Stop the containers
	docker compose down --volumes --remove-orphans

.PHONY: clean
clean:
	docker compose down --volumes --remove-orphans
