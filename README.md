# Proof of Work - Go Microservices with gRPC and Bazel

A proof-of-work server and client implementation using Go microservices architecture with gRPC communication and Bazel build system.

## Architecture

This project implements a microservices-based proof-of-work system with three main services:

### Services

1. **Challenge Service** (port 50051)
   - Generates proof-of-work challenges
   - Returns challenge ID, difficulty, and data
   - gRPC service

2. **Verifier Service** (port 50052)
   - Verifies proof-of-work solutions
   - Validates nonce against challenge requirements
   - gRPC service

3. **Gateway Service** (port 8080)
   - HTTP/REST API gateway
   - Communicates with backend services via gRPC
   - Provides user-friendly REST endpoints

### Communication Flow

```
Client (HTTP) → Gateway Service (HTTP→gRPC) → Challenge/Verifier Services (gRPC)
```

## Prerequisites

- Bazel 6.4.0 or later
- Go 1.21 or later (managed by Bazel)

## Building with Bazel

Build all services:
```bash
bazel build //...
```

Build individual services:
```bash
bazel build //services/challenge
bazel build //services/verifier
bazel build //services/gateway
bazel build //examples/client
```

## Running Tests

Run all tests:
```bash
bazel test //...
```

Run specific test:
```bash
bazel test //pkg/pow:pow_test
```

## Running the Services

### Start all services (in separate terminals):

Terminal 1 - Challenge Service:
```bash
bazel run //services/challenge
```

Terminal 2 - Verifier Service:
```bash
bazel run //services/verifier
```

Terminal 3 - Gateway Service:
```bash
bazel run //services/gateway
```

### Run the example client:
```bash
bazel run //examples/client
```

## API Endpoints

### Get Challenge
```bash
curl -X POST http://localhost:8080/v1/challenge \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 4, "data": "hello-world"}'
```

Response:
```json
{
  "challenge_id": "uuid",
  "difficulty": 4,
  "data": "hello-world",
  "timestamp": 1234567890
}
```

### Verify Proof
```bash
curl -X POST http://localhost:8080/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "challenge_id": "uuid",
    "nonce": 12345,
    "data": "hello-world",
    "difficulty": 4
  }'
```

Response:
```json
{
  "valid": true,
  "hash": "0000abc123...",
  "message": "Valid proof"
}
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Project Structure

```
.
├── api/proto/              # Protocol buffer definitions
│   ├── challenge.proto     # Challenge service proto
│   └── verifier.proto      # Verifier service proto
├── services/               # Microservices
│   ├── challenge/          # Challenge service implementation
│   ├── verifier/           # Verifier service implementation
│   └── gateway/            # HTTP gateway service
├── pkg/pow/                # Proof-of-work core library
├── examples/client/        # Example client application
├── third_party/            # Third-party dependencies
├── WORKSPACE               # Bazel workspace configuration
├── BUILD.bazel             # Root build file
├── deps.bzl                # Go dependencies
└── go.mod                  # Go module file
```

## How Proof-of-Work Works

1. Client requests a challenge with desired difficulty level
2. Challenge service generates a unique challenge ID and returns it with the data
3. Client computes SHA-256 hashes with different nonces until finding one with required leading zeros
4. Client submits the solution (nonce) to the verifier service
5. Verifier service validates the proof and returns the result

Example: For difficulty 4, the hash must start with "0000"

## Development

### Update Go dependencies:
```bash
bazel run //:gazelle-update-repos
```

### Format code:
```bash
bazel run //:gazelle
```

## License

MIT
