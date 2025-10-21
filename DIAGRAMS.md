# System Diagrams

## High-Level Architecture

```
                                 ┌─────────────────────┐
                                 │   Client (HTTP)     │
                                 │  - Browser          │
                                 │  - curl             │
                                 │  - Example client   │
                                 └──────────┬──────────┘
                                            │
                                            │ HTTP/REST
                                            │ JSON
                                            │
                                 ┌──────────▼──────────┐
                                 │  Gateway Service    │
                                 │  Port: 8080         │
                                 │  - POST /v1/challenge
                                 │  - POST /v1/verify  │
                                 │  - GET /health      │
                                 └──────────┬──────────┘
                                            │
                    ┌───────────────────────┴──────────────────────┐
                    │                                              │
                    │ gRPC                                         │ gRPC
                    │ Protocol Buffers                             │ Protocol Buffers
                    │                                              │
         ┌──────────▼───────────┐                      ┌──────────▼──────────┐
         │ Challenge Service    │                      │ Verifier Service    │
         │ Port: 50051          │                      │ Port: 50052         │
         │ - GetChallenge()     │                      │ - VerifyProof()     │
         └──────────────────────┘                      └─────────────────────┘
```

## Component Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Proof-of-Work System                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────────────────┐  ┌────────────────────┐  ┌───────────────┐ │
│  │  Gateway Service   │  │ Challenge Service  │  │    Verifier   │ │
│  │  ┌──────────────┐  │  │  ┌──────────────┐  │  │ ┌───────────┐ │ │
│  │  │ HTTP Handler │  │  │  │ gRPC Server  │  │  │ │gRPC Server│ │ │
│  │  └──────┬───────┘  │  │  └──────┬───────┘  │  │ └─────┬─────┘ │ │
│  │         │          │  │         │          │  │       │       │ │
│  │  ┌──────▼───────┐  │  │  ┌──────▼───────┐  │  │ ┌─────▼─────┐ │ │
│  │  │ gRPC Clients │  │  │  │  Challenge   │  │  │ │   PoW     │ │ │
│  │  │ - Challenge  │──┼──┼─▶│  Generator   │  │  │ │ Verifier  │ │ │
│  │  │ - Verifier   │──┼──┼─────────────────┼──┼─▶│           │ │ │
│  │  └──────────────┘  │  │  └──────────────┘  │  │ └───────────┘ │ │
│  └────────────────────┘  └────────────────────┘  └───────────────┘ │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                     Shared Library (pkg/pow)                  │  │
│  │  ┌───────────────┐  ┌──────────────┐  ┌──────────────────┐  │  │
│  │  │ ComputeHash() │  │VerifyProof() │  │   FindNonce()    │  │  │
│  │  └───────────────┘  └──────────────┘  └──────────────────┘  │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Sequence Diagram: Complete Flow

```
Client          Gateway         Challenge       Verifier        PoW Library
  │               │                 │               │                │
  │ POST          │                 │               │                │
  │ /v1/challenge │                 │               │                │
  ├──────────────▶│                 │               │                │
  │               │ GetChallenge    │               │                │
  │               │ (gRPC)          │               │                │
  │               ├────────────────▶│               │                │
  │               │                 │ Generate UUID │                │
  │               │                 │ Set difficulty│                │
  │               │                 │ Set timestamp │                │
  │               │    Response     │               │                │
  │               │◀────────────────┤               │                │
  │   JSON        │                 │               │                │
  │◀──────────────┤                 │               │                │
  │               │                 │               │                │
  │ Compute PoW   │                 │               │                │
  │ locally       │                 │               │                │
  ├──────────┐    │                 │               │                │
  │          │    │                 │               │                │
  │ FindNonce│    │                 │               │                │
  │ (local)  │    │                 │               │                │
  │◀─────────┘    │                 │               │                │
  │               │                 │               │                │
  │ POST          │                 │               │                │
  │ /v1/verify    │                 │               │                │
  ├──────────────▶│                 │               │                │
  │               │ VerifyProof     │               │                │
  │               │ (gRPC)          │               │                │
  │               ├────────────────────────────────▶│                │
  │               │                 │               │ ComputeHash()  │
  │               │                 │               ├───────────────▶│
  │               │                 │               │    hash        │
  │               │                 │               │◀───────────────┤
  │               │                 │               │ VerifyProof()  │
  │               │                 │               ├───────────────▶│
  │               │                 │               │   valid/invalid│
  │               │                 │               │◀───────────────┤
  │               │    Response     │               │                │
  │               │◀────────────────────────────────┤                │
  │   JSON        │                 │               │                │
  │◀──────────────┤                 │               │                │
  │               │                 │               │                │
```

## Deployment Diagram: Kubernetes

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Kubernetes Cluster                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                         Namespace: pow                      │    │
│  │                                                             │    │
│  │  ┌────────────────────────────────────────────────────┐    │    │
│  │  │              Ingress / Load Balancer               │    │    │
│  │  │              External IP: x.x.x.x:80               │    │    │
│  │  └────────────────────┬───────────────────────────────┘    │    │
│  │                       │                                     │    │
│  │  ┌────────────────────▼───────────────────────────────┐    │    │
│  │  │           Service: gateway-service                  │    │    │
│  │  │           ClusterIP: 10.0.0.10:8080                │    │    │
│  │  └────────────────────┬───────────────────────────────┘    │    │
│  │                       │                                     │    │
│  │  ┌────────────────────▼───────────────────────────────┐    │    │
│  │  │         Deployment: gateway (3 replicas)           │    │    │
│  │  │  ┌──────┐  ┌──────┐  ┌──────┐                     │    │    │
│  │  │  │ Pod  │  │ Pod  │  │ Pod  │                     │    │    │
│  │  │  └──────┘  └──────┘  └──────┘                     │    │    │
│  │  └────────────────────┬─────────┬─────────────────────┘    │    │
│  │                       │         │                           │    │
│  │         ┌─────────────┘         └──────────────┐            │    │
│  │         │                                      │            │    │
│  │  ┌──────▼───────────┐              ┌──────────▼────────┐   │    │
│  │  │ Service:         │              │ Service:          │   │    │
│  │  │ challenge-svc    │              │ verifier-svc      │   │    │
│  │  │ 10.0.0.11:50051  │              │ 10.0.0.12:50052   │   │    │
│  │  └──────┬───────────┘              └──────────┬────────┘   │    │
│  │         │                                     │            │    │
│  │  ┌──────▼───────────┐              ┌──────────▼────────┐   │    │
│  │  │ Deployment:      │              │ Deployment:       │   │    │
│  │  │ challenge        │              │ verifier          │   │    │
│  │  │ (2 replicas)     │              │ (2 replicas)      │   │    │
│  │  │ ┌────┐  ┌────┐   │              │ ┌────┐  ┌────┐   │   │    │
│  │  │ │Pod │  │Pod │   │              │ │Pod │  │Pod │   │   │    │
│  │  │ └────┘  └────┘   │              │ └────┘  └────┘   │   │    │
│  │  └──────────────────┘              └───────────────────┘   │    │
│  │                                                             │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Data Flow Diagram

```
                                Challenge Flow
                        ┌──────────────────────────┐
                        │   1. Request Challenge   │
                        │   POST /v1/challenge     │
                        │   {difficulty: 4}        │
                        └────────────┬─────────────┘
                                     │
                        ┌────────────▼─────────────┐
                        │   2. Generate Challenge  │
                        │   - UUID                 │
                        │   - Difficulty: 4        │
                        │   - Timestamp            │
                        └────────────┬─────────────┘
                                     │
                        ┌────────────▼─────────────┐
                        │   3. Return Challenge    │
                        │   challenge_id: uuid     │
                        │   difficulty: 4          │
                        │   data: "..."            │
                        └──────────────────────────┘

                            Verification Flow
                        ┌──────────────────────────┐
                        │   1. Client Computes     │
                        │   nonce = FindNonce()    │
                        │   (local computation)    │
                        └────────────┬─────────────┘
                                     │
                        ┌────────────▼─────────────┐
                        │   2. Submit Proof        │
                        │   POST /v1/verify        │
                        │   {challenge_id, nonce}  │
                        └────────────┬─────────────┘
                                     │
                        ┌────────────▼─────────────┐
                        │   3. Verify Proof        │
                        │   hash = ComputeHash()   │
                        │   valid = VerifyProof()  │
                        └────────────┬─────────────┘
                                     │
                        ┌────────────▼─────────────┐
                        │   4. Return Result       │
                        │   valid: true/false      │
                        │   hash: "0000abc..."     │
                        └──────────────────────────┘
```

## Protocol Buffer Structure

```
┌─────────────────────────────────────────────────────────────┐
│                   challenge.proto                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  service ChallengeService {                                 │
│    rpc GetChallenge(GetChallengeRequest)                    │
│        returns (GetChallengeResponse);                      │
│  }                                                           │
│                                                              │
│  message GetChallengeRequest {                              │
│    int32 difficulty = 1;                                    │
│    string data = 2;                                         │
│  }                                                           │
│                                                              │
│  message GetChallengeResponse {                             │
│    string challenge_id = 1;                                 │
│    int32 difficulty = 2;                                    │
│    string data = 3;                                         │
│    int64 timestamp = 4;                                     │
│  }                                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   verifier.proto                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  service VerifierService {                                  │
│    rpc VerifyProof(VerifyProofRequest)                      │
│        returns (VerifyProofResponse);                       │
│  }                                                           │
│                                                              │
│  message VerifyProofRequest {                               │
│    string challenge_id = 1;                                 │
│    int64 nonce = 2;                                         │
│    string data = 3;                                         │
│    int32 difficulty = 4;                                    │
│  }                                                           │
│                                                              │
│  message VerifyProofResponse {                              │
│    bool valid = 1;                                          │
│    string hash = 2;                                         │
│    string message = 3;                                      │
│  }                                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Build Dependency Graph

```
                        ┌─────────────────┐
                        │   WORKSPACE     │
                        │  (Root config)  │
                        └────────┬────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
         ┌──────────▼────────┐    ┌──────────▼────────┐
         │ rules_go          │    │ rules_proto       │
         │ (Go toolchain)    │    │ (Proto toolchain) │
         └──────────┬────────┘    └──────────┬────────┘
                    │                         │
         ┌──────────▼────────┐    ┌──────────▼────────┐
         │   deps.bzl        │    │  api/proto        │
         │ (Go dependencies) │    │  (Proto files)    │
         └──────────┬────────┘    └──────────┬────────┘
                    │                         │
                    │             ┌───────────▼────────┐
                    │             │ challenge_go_proto │
                    │             │ verifier_go_proto  │
                    │             └───────────┬────────┘
                    │                         │
         ┌──────────▼────────┐    ┌───────────▼────────┐
         │  pkg/pow          │    │  services/*        │
         │  (Core library)   │◀───┤  (Microservices)   │
         └───────────────────┘    └────────────────────┘
                    │
         ┌──────────▼────────┐
         │ examples/client   │
         │ (Demo client)     │
         └───────────────────┘
```

## Directory Structure Tree

```
pow/
├── .bazelrc                    # Bazel build configuration
├── .bazelversion              # Bazel version pinning
├── .gitignore                 # Git ignore rules
├── BUILD.bazel                # Root build file
├── WORKSPACE                  # Bazel workspace definition
├── go.mod                     # Go module definition
├── deps.bzl                   # Go dependencies for Bazel
├── Makefile                   # Convenience commands
├── docker-compose.yml         # Container orchestration
│
├── README.md                  # Main documentation
├── QUICKSTART.md              # Quick start guide
├── ARCHITECTURE.md            # Architecture details
├── TESTING.md                 # Testing guide
├── DIAGRAMS.md               # System diagrams
├── IMPLEMENTATION_SUMMARY.md  # Implementation summary
│
├── api/proto/                 # Protocol Buffer definitions
│   ├── BUILD.bazel            # Proto build rules
│   ├── challenge.proto        # Challenge service API
│   └── verifier.proto         # Verifier service API
│
├── pkg/pow/                   # Core proof-of-work library
│   ├── BUILD.bazel            # Library build rules
│   ├── pow.go                 # Implementation
│   └── pow_test.go            # Unit tests
│
├── services/                  # Microservices
│   ├── challenge/             # Challenge service
│   │   ├── BUILD.bazel        # Service build rules
│   │   └── server.go          # gRPC server
│   │
│   ├── verifier/              # Verifier service
│   │   ├── BUILD.bazel        # Service build rules
│   │   └── server.go          # gRPC server
│   │
│   └── gateway/               # Gateway service
│       ├── BUILD.bazel        # Service build rules
│       └── server.go          # HTTP/gRPC gateway
│
├── examples/                  # Example applications
│   └── client/                # Example client
│       ├── BUILD.bazel        # Client build rules
│       └── main.go            # Client implementation
│
├── scripts/                   # Helper scripts
│   └── start-services.sh      # Start all services
│
└── third_party/               # Third-party dependencies
    └── googleapis/            # Google APIs
        └── google/api/        # HTTP annotations
            ├── BUILD.bazel
            ├── annotations.proto
            └── http.proto
```

## Technology Stack Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│  Gateway Service  │  Challenge Service  │  Verifier Service │
│  (HTTP Handler)   │   (gRPC Server)    │   (gRPC Server)   │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Framework Layer                         │
├─────────────────────────────────────────────────────────────┤
│  net/http (HTTP)  │  google.golang.org/grpc (gRPC)         │
│  encoding/json    │  google.golang.org/protobuf            │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Business Logic                          │
├─────────────────────────────────────────────────────────────┤
│  pkg/pow (Proof-of-Work library)                            │
│  - ComputeHash (SHA-256)                                    │
│  - VerifyProof                                              │
│  - FindNonce                                                │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                     Build & Tools Layer                      │
├─────────────────────────────────────────────────────────────┤
│  Bazel (Build System)                                       │
│  - rules_go (Go toolchain)                                  │
│  - rules_proto (Protocol Buffers)                           │
│  - Gazelle (Dependency management)                          │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Runtime Layer                           │
├─────────────────────────────────────────────────────────────┤
│  Go 1.21 Runtime                                            │
│  - Goroutines (Concurrency)                                 │
│  - HTTP/2 (gRPC transport)                                  │
│  - TCP/IP networking                                        │
└─────────────────────────────────────────────────────────────┘
```

These diagrams provide a comprehensive visual representation of the proof-of-work microservices system, showing the architecture, data flows, deployment options, and technology stack.
