.PHONY: test
test:
	# Start the containers
	docker compose up -d --build
	sleep 3

	# Run the tests
	bash -c "set -o pipefail; go test -v ./integration_test.go | tee /dev/stderr"; test_exit_code=$$?

	# Stop the containers
	docker compose down --volumes --remove-orphans

	# Exit with the test exit code
	exit $$test_exit_code
