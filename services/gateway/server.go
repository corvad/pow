package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	challengepb "github.com/corvad/pow/api/proto/challenge/v1"
	verifierpb "github.com/corvad/pow/api/proto/verifier/v1"
)

type Gateway struct {
	challengeClient challengepb.ChallengeServiceClient
	verifierClient  verifierpb.VerifierServiceClient
}

func NewGateway(challengeAddr, verifierAddr string) (*Gateway, error) {
	// Connect to challenge service
	challengeConn, err := grpc.Dial(challengeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to challenge service: %w", err)
	}

	// Connect to verifier service
	verifierConn, err := grpc.Dial(verifierAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		challengeConn.Close()
		return nil, fmt.Errorf("failed to connect to verifier service: %w", err)
	}

	return &Gateway{
		challengeClient: challengepb.NewChallengeServiceClient(challengeConn),
		verifierClient:  verifierpb.NewVerifierServiceClient(verifierConn),
	}, nil
}

func (g *Gateway) handleGetChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Difficulty int32  `json:"difficulty"`
		Data       string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.challengeClient.GetChallenge(ctx, &challengepb.GetChallengeRequest{
		Difficulty: req.Difficulty,
		Data:       req.Data,
	})
	if err != nil {
		log.Printf("Error getting challenge: %v", err)
		http.Error(w, "Failed to get challenge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleVerifyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChallengeID string `json:"challenge_id"`
		Nonce       int64  `json:"nonce"`
		Data        string `json:"data"`
		Difficulty  int32  `json:"difficulty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.verifierClient.VerifyProof(ctx, &verifierpb.VerifyProofRequest{
		ChallengeId: req.ChallengeID,
		Nonce:       req.Nonce,
		Data:        req.Data,
		Difficulty:  req.Difficulty,
	})
	if err != nil {
		log.Printf("Error verifying proof: %v", err)
		http.Error(w, "Failed to verify proof", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func main() {
	challengeAddr := "localhost:50051"
	verifierAddr := "localhost:50052"
	httpPort := 8080

	gateway, err := NewGateway(challengeAddr, verifierAddr)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	http.HandleFunc("/v1/challenge", gateway.handleGetChallenge)
	http.HandleFunc("/v1/verify", gateway.handleVerifyProof)
	http.HandleFunc("/health", gateway.handleHealth)

	addr := fmt.Sprintf(":%d", httpPort)
	log.Printf("Gateway listening on %s", addr)
	log.Printf("Challenge service: %s", challengeAddr)
	log.Printf("Verifier service: %s", verifierAddr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
