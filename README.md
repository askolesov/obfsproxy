# obfsproxy

obfsproxy is a simple obfuscating proxy designed to provide a layer of obfuscation for network traffic. It works by inverting the bytes of the data passing through it, making it harder for network monitors to identify the content of the communication.

## Features

- Simple byte inversion obfuscation
- Bidirectional proxying
- Command-line interface for easy configuration
- Docker support for containerized deployment
- Integration tests to ensure proper functionality

## Installation

### Prerequisites

- Go 1.23.2 or later
- Docker (for containerized deployment and integration tests)

### Building from source

1. Clone the repository:
   ```
   git clone https://github.com/askolesov/obfsproxy.git
   cd obfsproxy
   ```

2. Build the project:
   ```
   make build
   ```

The compiled binary will be available in the `build/` directory.

## Usage

Run obfsproxy with the following command:

```
obfsproxy -l [listen_address] -t [target_address]
```

- `-l, --listen`: Address to listen on (default: "localhost:8080")
- `-t, --target`: Address to forward to (default: "localhost:80")

Example:
```
obfsproxy -l 0.0.0.0:8080 -t example.com:80
```

This will start the proxy listening on all interfaces on port 8080 and forward traffic to example.com on port 80.

## Docker

To build and run obfsproxy using Docker:

1. Build the Docker image:
   ```
   make docker
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 obfsproxy:latest obfsproxy -l 0.0.0.0:8080 -t example.com:80
   ```

## Development

### Running tests

To run the unit tests:

```
make test
```

### Running integration tests

To run the integration tests:

```
cd integration
make test
```

This will start the necessary Docker containers and run the integration tests.

### Linting

To run the linter:

```
make lint
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
