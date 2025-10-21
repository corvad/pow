package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/corvad/pow/pkg/pow"
)

type ChallengeResponse struct {
	ChallengeID string `json:"challenge_id"`
	Difficulty  int32  `json:"difficulty"`
	Data        string `json:"data"`
	Timestamp   int64  `json:"timestamp"`
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

func main() {
	gatewayURL := "http://localhost:8080"

	// Step 1: Get a challenge
	log.Println("Step 1: Requesting a challenge...")
	challengeReq := map[string]interface{}{
		"difficulty": 4,
		"data":       "my-proof-of-work",
	}
	
	reqBody, _ := json.Marshal(challengeReq)
	resp, err := http.Post(gatewayURL+"/v1/challenge", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Failed to get challenge: %v", err)
	}
	defer resp.Body.Close()

	var challenge ChallengeResponse
	if err := json.NewDecoder(resp.Body).Decode(&challenge); err != nil {
		log.Fatalf("Failed to decode challenge response: %v", err)
	}

	log.Printf("Received challenge: ID=%s, Difficulty=%d, Data=%s", 
		challenge.ChallengeID, challenge.Difficulty, challenge.Data)

	// Step 2: Find a valid nonce (proof-of-work)
	log.Println("\nStep 2: Computing proof-of-work...")
	nonce, hash, found := pow.FindNonce(challenge.ChallengeID, challenge.Data, challenge.Difficulty, 1000000)
	if !found {
		log.Fatal("Failed to find valid nonce within maximum attempts")
	}

	log.Printf("Found valid nonce: %d", nonce)
	log.Printf("Hash: %s", hash)

	// Step 3: Verify the proof
	log.Println("\nStep 3: Verifying proof...")
	verifyReq := map[string]interface{}{
		"challenge_id": challenge.ChallengeID,
		"nonce":        nonce,
		"data":         challenge.Data,
		"difficulty":   challenge.Difficulty,
	}

	reqBody, _ = json.Marshal(verifyReq)
	resp, err = http.Post(gatewayURL+"/v1/verify", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	defer resp.Body.Close()

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		log.Fatalf("Failed to decode verify response: %v", err)
	}

	log.Printf("Verification result: Valid=%v, Message=%s", verifyResp.Valid, verifyResp.Message)
	
	if verifyResp.Valid {
		fmt.Println("\n✓ Proof-of-work successfully verified!")
	} else {
		fmt.Println("\n✗ Proof-of-work verification failed!")
	}
}
