# Quick Start Guide

Get up and running with the Proof-of-Work microservices in 5 minutes!

## Prerequisites

- Bazel 6.4.0+ (will be auto-downloaded by bazelisk)
- Go 1.21+ (optional, managed by Bazel)
- Git

## Installation

1. **Clone the repository:**
```bash
git clone https://github.com/corvad/pow.git
cd pow
```

## Running the System

### Option 1: Using the convenience script

```bash
./scripts/start-services.sh
```

This will start all three services in the background.

### Option 2: Manual start (3 separate terminals)

**Terminal 1 - Challenge Service:**
```bash
bazel run //services/challenge
```

**Terminal 2 - Verifier Service:**
```bash
bazel run //services/verifier
```

**Terminal 3 - Gateway Service:**
```bash
bazel run //services/gateway
```

## Testing the API

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Get a Challenge
```bash
curl -X POST http://localhost:8080/v1/challenge \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 4, "data": "hello-world"}'
```

You'll get a response like:
```json
{
  "challenge_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "difficulty": 4,
  "data": "hello-world",
  "timestamp": 1234567890
}
```

### 3. Run the Complete Example

The easiest way to see everything in action:
```bash
bazel run //examples/client
```

This will:
1. Request a challenge
2. Compute a proof-of-work
3. Verify the proof
4. Show the results

## What's Next?

- Read [README.md](README.md) for complete documentation
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for system design details
- See [TESTING.md](TESTING.md) for testing strategies
- Review the code in:
  - `services/challenge/` - Challenge generation
  - `services/verifier/` - Proof verification
  - `services/gateway/` - HTTP gateway
  - `pkg/pow/` - Core proof-of-work logic

## Common Commands

```bash
# Build everything
bazel build //...

# Run all tests
bazel test //...

# Build specific service
bazel build //services/gateway

# Clean build artifacts
bazel clean

# Update Go dependencies
bazel run //:gazelle-update-repos
```

## Architecture Overview

```
┌─────────┐
│ Client  │
└────┬────┘
     │ HTTP/REST (port 8080)
     ▼
┌─────────────┐
│   Gateway   │
└──────┬──────┘
       │ gRPC
       ├──────────────┐
       │              │
       ▼              ▼
┌──────────┐   ┌──────────┐
│Challenge │   │ Verifier │
│(50051)   │   │ (50052)  │
└──────────┘   └──────────┘
```

## Troubleshooting

### Services won't start

Check if ports are already in use:
```bash
lsof -i :50051
lsof -i :50052
lsof -i :8080
```

### Bazel build fails

Try cleaning and rebuilding:
```bash
bazel clean
bazel build //...
```

### Need help?

Check the [TESTING.md](TESTING.md) file for debugging tips.

## Example: Computing Proof-of-Work

Here's what happens under the hood:

1. **Challenge**: `challenge_id="abc", data="hello", difficulty=3`
2. **Compute Hash**: `SHA256("abc:hello:nonce")`
3. **Find Nonce**: Try nonce=0,1,2,... until hash starts with "000"
4. **Example**: nonce=12345 → hash="000a1b2c3d..."
5. **Verify**: Submit nonce=12345 for verification ✓

The difficulty determines how many leading zeros are required:
- Difficulty 1: ~16 attempts average
- Difficulty 2: ~256 attempts average
- Difficulty 3: ~4,096 attempts average
- Difficulty 4: ~65,536 attempts average
- Difficulty 5: ~1,048,576 attempts average

## Features

✅ Microservices architecture  
✅ gRPC for service-to-service communication  
✅ REST API for client access  
✅ Built with Bazel  
✅ Protocol Buffers for type-safe APIs  
✅ Unit tested  
✅ Example client included  
✅ Docker Compose support  

## Next Steps

Now that you have the system running, try:

1. **Modify the difficulty**: Change the difficulty parameter in challenges
2. **Add features**: Extend the services with new functionality
3. **Deploy**: Use the docker-compose.yml for containerized deployment
4. **Scale**: Run multiple instances of each service behind a load balancer
5. **Monitor**: Add Prometheus metrics and Grafana dashboards

Enjoy exploring the Proof-of-Work microservices! 🚀
