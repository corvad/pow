#!/bin/bash

# Script to start all services for local development

set -e

echo "Starting Proof-of-Work Microservices..."
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if running in tmux or screen
if [ -z "$TMUX" ] && [ -z "$STY" ]; then
    echo "Warning: This script works best in tmux or screen for multiple panes"
    echo "Starting services in background..."
    echo ""
fi

# Build all services
echo -e "${BLUE}Building all services...${NC}"
bazel build //services/challenge //services/verifier //services/gateway

# Start challenge service
echo -e "${GREEN}Starting Challenge Service on port 50051...${NC}"
bazel run //services/challenge &
CHALLENGE_PID=$!
sleep 2

# Start verifier service
echo -e "${GREEN}Starting Verifier Service on port 50052...${NC}"
bazel run //services/verifier &
VERIFIER_PID=$!
sleep 2

# Start gateway service
echo -e "${GREEN}Starting Gateway Service on port 8080...${NC}"
bazel run //services/gateway &
GATEWAY_PID=$!
sleep 2

echo ""
echo -e "${GREEN}All services started successfully!${NC}"
echo ""
echo "Services:"
echo "  - Challenge Service: gRPC on port 50051 (PID: $CHALLENGE_PID)"
echo "  - Verifier Service:  gRPC on port 50052 (PID: $VERIFIER_PID)"
echo "  - Gateway Service:   HTTP on port 8080 (PID: $GATEWAY_PID)"
echo ""
echo "Test the API:"
echo "  curl -X POST http://localhost:8080/v1/challenge -H 'Content-Type: application/json' -d '{\"difficulty\": 4}'"
echo ""
echo "Run example client:"
echo "  bazel run //examples/client"
echo ""
echo "Press Ctrl+C to stop all services"

# Trap Ctrl+C and stop all services
trap "echo 'Stopping services...'; kill $CHALLENGE_PID $VERIFIER_PID $GATEWAY_PID 2>/dev/null; exit" INT TERM

# Wait for all processes
wait
