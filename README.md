# NanoRediz

A production-ready distributed Redis-like key-value store implemented in Go with Raft consensus algorithm.

## Features

- **Distributed Architecture**: Uses Raft consensus for high availability and consistency
- **Redis-Compatible Commands**: Supports GET, SET, DEL, APPEND, STRLN, PING operations
- **Production Ready**: Includes graceful shutdown, configuration management, and comprehensive error handling
- **Multiple Interfaces**: gRPC API, CLI client, and web interface
- **Docker Support**: Full containerization with Docker and Docker Compose
- **Monitoring**: Built-in health checks and metrics
- **Configurable**: Environment variables and command-line configuration

## Quick Start

### Using Pre-built Binaries

1. **Build the project:**
   ```bash
   make build
   ```

2. **Start a single node:**
   ```bash
   ./bin/nanorediz-server -host 127.0.0.1 -port 8080
   ```

3. **Start a cluster:**
   ```bash
   # Terminal 1: Start first node
   ./bin/nanorediz-server -host 127.0.0.1 -port 8080
   
   # Terminal 2: Start second node 
   ./bin/nanorediz-server -host 127.0.0.1 -port 8081 -contact-host 127.0.0.1 -contact-port 8080
   
   # Terminal 3: Start third node
   ./bin/nanorediz-server -host 127.0.0.1 -port 8082 -contact-host 127.0.0.1 -contact-port 8080
   ```

4. **Connect with client:**
   ```bash
   ./bin/nanorediz-client -host 127.0.0.1 -port 8080
   ```

### Using Docker

1. **Single node:**
   ```bash
   docker build -t nanorediz .
   docker run -p 8080:8080 nanorediz
   ```

2. **Full cluster:**
   ```bash
   docker-compose up
   ```

### Using Make (Development)

```bash
# Start development cluster
make dev-cluster

# Stop development cluster
make stop-cluster

# Run tests
make test

# Run with coverage
make test-coverage
```

## Commands

The client supports the following Redis-compatible commands:

- `ping` - Test connectivity
- `get <key>` - Get value by key
- `set <key> <value>` - Set key-value pair
- `del <key>` - Delete key
- `append <key> <value>` - Append value to existing key
- `strln <key>` - Get string length of value
- `request_log` - Get server request log
- `quit` - Exit client

Example session:
```bash
127.0.0.1:8080> set mykey "hello world"
OK
127.0.0.1:8080> get mykey
"hello world"
127.0.0.1:8080> strln mykey
11
127.0.0.1:8080> append mykey "!"
OK
127.0.0.1:8080> get mykey
"hello world!"
```

## Configuration

### Command Line Options

**Server:**
```bash
./bin/nanorediz-server [OPTIONS]

Options:
  -host string           Server host address (default "127.0.0.1")
  -port string           Server port (default "8080")
  -contact-host string   Contact node host (for joining cluster)
  -contact-port string   Contact node port (for joining cluster)
  -help                  Show help message
```

**Client:**
```bash
./bin/nanorediz-client [OPTIONS]

Options:
  -host string    Server host (default "127.0.0.1")
  -port string    Server port (default "8080")
  -time           Enable response time measurement
```

### Environment Variables

- `NANOREDIZ_HOST` - Server host (default: 0.0.0.0)
- `NANOREDIZ_PORT` - Server port (default: 8080)
- `NANOREDIZ_LOG_LEVEL` - Log level: debug, info, warn, error (default: info)
- `NANOREDIZ_LOG_FORMAT` - Log format: json, text (default: json)
- `NANOREDIZ_GRPC_TIMEOUT` - gRPC request timeout (default: 30s)
- `NANOREDIZ_SHUTDOWN_TIMEOUT` - Graceful shutdown timeout (default: 10s)

## Architecture

NanoRediz is built with a modular architecture:

```
├── cmd/                 # Application entry points
│   ├── client/         # CLI client
│   └── server/         # Server application
├── lib/                # Core libraries
│   ├── app/           # Key-value store implementation
│   ├── client/        # Client library
│   ├── config/        # Configuration management
│   ├── logger/        # Logging utilities
│   ├── pb/            # Protocol Buffer definitions
│   ├── raft/          # Raft consensus implementation
│   ├── server/        # gRPC server
│   ├── service/       # Business logic services
│   └── util/          # Utility functions
├── web/               # Web interface
└── test/              # Tests
```

### Key Components

- **Raft Consensus**: Ensures data consistency across distributed nodes
- **gRPC Communication**: High-performance RPC for inter-node communication
- **Key-Value Store**: In-memory data structure with persistence support
- **Web Interface**: HTTP/HTML interface for browser access
- **Configuration System**: Environment and file-based configuration

## Development

### Prerequisites

- Go 1.20 or later
- Make
- Docker (optional)
- Protocol Buffers compiler (for development)

### Building from Source

```bash
# Clone repository
git clone https://github.com/0xzre/nanorediz.git
cd nanorediz

# Install dependencies
make deps

# Build all binaries
make build

# Run tests
make test

# Format code
make fmt

# Install development tools
make install-tools

# Run linter
make lint
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

## Production Deployment

### Docker Deployment

1. **Build production image:**
   ```bash
   docker build -t nanorediz:latest .
   ```

2. **Deploy cluster:**
   ```bash
   docker-compose up -d
   ```

3. **Scale cluster:**
   ```bash
   docker-compose up -d --scale nanorediz-node-2=2
   ```

### Health Monitoring

The server provides health check endpoints and supports graceful shutdown:

- Health check via ping command
- Graceful shutdown on SIGTERM/SIGINT
- Configurable timeouts
- Request/response logging

### Performance Considerations

- Adjust `NANOREDIZ_GRPC_TIMEOUT` based on network latency
- Configure appropriate Raft timeouts for your network
- Monitor memory usage for large datasets
- Use load balancers for client connections

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make changes and add tests
4. Run tests: `make test`
5. Format code: `make fmt`
6. Run linter: `make lint`
7. Commit changes: `git commit -am 'Add new feature'`
8. Push to branch: `git push origin feature/new-feature`
9. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Authors

| Name                           | NIM      |
| ------------------------------ | -------- |
| Rizky Abdillah Rasyid         | 13521109 |
| Muhammad Rizky Syaban         | 13521119 |
| M. Abdul Aziz Ghazali         | 13521128 |
| M. Dimas Sakti Widyaatmaja    | 13521160 |

## Acknowledgments

- Built as part of IF3230 (Distributed Systems) coursework
- Inspired by Redis and etcd architectures
- Uses Raft consensus algorithm for distributed coordination
