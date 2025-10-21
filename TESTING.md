# Testing Guide

## Running Tests with Bazel

### Run All Tests
```bash
bazel test //...
```

### Run Specific Test
```bash
bazel test //pkg/pow:pow_test
```

### Run Tests with Verbose Output
```bash
bazel test //... --test_output=all
```

### Run Tests with Coverage
```bash
bazel coverage //...
```

## Manual Testing

### Prerequisites
Ensure all services are running:

```bash
# Terminal 1: Challenge Service
bazel run //services/challenge

# Terminal 2: Verifier Service  
bazel run //services/verifier

# Terminal 3: Gateway Service
bazel run //services/gateway
```

Or use the convenience script:
```bash
./scripts/start-services.sh
```

### Test 1: Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy"
}
```

### Test 2: Get Challenge

```bash
curl -X POST http://localhost:8080/v1/challenge \
  -H "Content-Type: application/json" \
  -d '{
    "difficulty": 4,
    "data": "test-challenge"
  }'
```

Expected response:
```json
{
  "challenge_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "difficulty": 4,
  "data": "test-challenge",
  "timestamp": 1234567890
}
```

### Test 3: Verify Invalid Proof

```bash
curl -X POST http://localhost:8080/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "challenge_id": "test-id",
    "nonce": 123,
    "data": "test-data",
    "difficulty": 4
  }'
```

Expected response (likely invalid):
```json
{
  "valid": false,
  "hash": "1234567890abcdef...",
  "message": "Invalid proof"
}
```

### Test 4: Complete Flow with Client

Run the example client that performs the full flow:

```bash
bazel run //examples/client
```

Expected output:
```
2024/01/01 12:00:00 Step 1: Requesting a challenge...
2024/01/01 12:00:00 Received challenge: ID=uuid, Difficulty=4, Data=my-proof-of-work

2024/01/01 12:00:00 Step 2: Computing proof-of-work...
2024/01/01 12:00:01 Found valid nonce: 12345
2024/01/01 12:00:01 Hash: 0000abc123def456...

2024/01/01 12:00:01 Step 3: Verifying proof...
2024/01/01 12:00:01 Verification result: Valid=true, Message=Valid proof

✓ Proof-of-work successfully verified!
```

## Testing with grpcurl

If you have grpcurl installed, you can test gRPC services directly:

### Challenge Service

```bash
grpcurl -plaintext \
  -d '{"difficulty": 4, "data": "test"}' \
  localhost:50051 \
  pow.challenge.v1.ChallengeService/GetChallenge
```

### Verifier Service

```bash
grpcurl -plaintext \
  -d '{"challenge_id": "test-id", "nonce": 123, "data": "test", "difficulty": 4}' \
  localhost:50052 \
  pow.verifier.v1.VerifierService/VerifyProof
```

### List Services (with reflection enabled)

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50052 list
```

## Load Testing

### Using Apache Bench (ab)

Test gateway endpoint performance:

```bash
# 1000 requests, 10 concurrent
ab -n 1000 -c 10 -p challenge.json -T application/json \
  http://localhost:8080/v1/challenge
```

Where `challenge.json` contains:
```json
{"difficulty": 4, "data": "load-test"}
```

### Using hey

```bash
# Install hey: go install github.com/rakyll/hey@latest

hey -n 1000 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 4}' \
  http://localhost:8080/v1/challenge
```

## Performance Testing

### Measure Proof-of-Work Computation Time

Create a test script to measure nonce finding time:

```go
package main

import (
    "fmt"
    "time"
    "github.com/corvad/pow/pkg/pow"
)

func main() {
    difficulties := []int32{1, 2, 3, 4, 5}
    
    for _, diff := range difficulties {
        start := time.Now()
        nonce, _, found := pow.FindNonce("test", "data", diff, 10000000)
        duration := time.Since(start)
        
        if found {
            fmt.Printf("Difficulty %d: Found nonce %d in %v\n", 
                diff, nonce, duration)
        } else {
            fmt.Printf("Difficulty %d: Not found within limit\n", diff)
        }
    }
}
```

## Integration Testing

### Test Service Communication

1. Start all services
2. Send challenge request to gateway
3. Verify gateway communicates with challenge service
4. Use returned challenge to compute proof
5. Send verification request to gateway
6. Verify gateway communicates with verifier service
7. Confirm verification result

### Test Error Handling

```bash
# Test with invalid JSON
curl -X POST http://localhost:8080/v1/challenge \
  -H "Content-Type: application/json" \
  -d 'invalid json'

# Test with missing fields
curl -X POST http://localhost:8080/v1/verify \
  -H "Content-Type: application/json" \
  -d '{}'

# Test with wrong HTTP method
curl -X GET http://localhost:8080/v1/challenge
```

## Debugging

### View Service Logs

When running services in separate terminals, logs will appear in each terminal.

### Test gRPC Connectivity

```bash
# Check if services are listening
netstat -an | grep -E '(50051|50052|8080)'

# Test TCP connection
telnet localhost 50051
telnet localhost 50052
telnet localhost 8080
```

### View Bazel Test Logs

```bash
# After running tests, view detailed logs
bazel test //pkg/pow:pow_test --test_output=all

# View test.log for a specific test
cat bazel-testlogs/pkg/pow/pow_test/test.log
```

## Continuous Integration

### Example GitHub Actions Workflow

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Bazel
        uses: bazelbuild/setup-bazelisk@v2
      
      - name: Run tests
        run: bazel test //...
      
      - name: Build all
        run: bazel build //...
```

## Troubleshooting

### Port Already in Use

```bash
# Find and kill process using port
lsof -ti:50051 | xargs kill -9
lsof -ti:50052 | xargs kill -9
lsof -ti:8080 | xargs kill -9
```

### gRPC Connection Refused

- Ensure services are running
- Check firewall settings
- Verify correct ports in gateway configuration

### Bazel Build Failures

```bash
# Clean and rebuild
bazel clean
bazel build //...

# Clean everything including external dependencies
bazel clean --expunge
bazel build //...
```

### Protocol Buffer Errors

```bash
# Regenerate proto files
bazel build //api/proto:challenge_go_proto
bazel build //api/proto:verifier_go_proto
```

## Test Coverage

### Generate Coverage Report

```bash
bazel coverage --combined_report=lcov //...
```

### View Coverage in HTML

```bash
genhtml bazel-out/_coverage/_coverage_report.dat -o coverage_html
open coverage_html/index.html
```

## Best Practices

1. **Always run tests before committing**: `bazel test //...`
2. **Test error cases**: Not just happy path
3. **Test edge cases**: Boundary conditions, empty inputs
4. **Test concurrency**: Multiple simultaneous requests
5. **Monitor performance**: Track response times
6. **Use structured logging**: Makes debugging easier
7. **Keep tests fast**: Unit tests should run in milliseconds
8. **Make tests deterministic**: No random failures
