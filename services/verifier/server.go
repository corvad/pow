package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	verifierpb "github.com/corvad/pow/api/proto/verifier/v1"
	"github.com/corvad/pow/pkg/pow"
)

type verifierServer struct {
	verifierpb.UnimplementedVerifierServiceServer
}

func (s *verifierServer) VerifyProof(ctx context.Context, req *verifierpb.VerifyProofRequest) (*verifierpb.VerifyProofResponse, error) {
	// Compute the hash with the provided nonce
	hash := pow.ComputeHash(req.ChallengeId, req.Data, req.Nonce)
	
	// Verify the proof meets the difficulty requirement
	valid := pow.VerifyProof(hash, req.Difficulty)
	
	message := "Invalid proof"
	if valid {
		message = "Valid proof"
	}
	
	log.Printf("Verified proof: ChallengeID=%s, Nonce=%d, Valid=%v, Hash=%s", 
		req.ChallengeId, req.Nonce, valid, hash)
	
	return &verifierpb.VerifyProofResponse{
		Valid:   valid,
		Hash:    hash,
		Message: message,
	}, nil
}

func main() {
	port := 50052
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	verifierpb.RegisterVerifierServiceServer(grpcServer, &verifierServer{})
	
	// Register reflection service for debugging
	reflection.Register(grpcServer)

	log.Printf("Verifier service listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
