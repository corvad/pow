package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	challengepb "github.com/corvad/pow/api/proto/challenge/v1"
	"github.com/google/uuid"
)

type challengeServer struct {
	challengepb.UnimplementedChallengeServiceServer
}

func (s *challengeServer) GetChallenge(ctx context.Context, req *challengepb.GetChallengeRequest) (*challengepb.GetChallengeResponse, error) {
	// Generate a unique challenge ID
	challengeID := uuid.New().String()
	
	// Default difficulty if not specified
	difficulty := req.Difficulty
	if difficulty <= 0 {
		difficulty = 4 // Default to 4 leading zeros
	}
	
	// Use provided data or default message
	data := req.Data
	if data == "" {
		data = "proof-of-work-challenge"
	}
	
	log.Printf("Generated challenge: ID=%s, Difficulty=%d", challengeID, difficulty)
	
	return &challengepb.GetChallengeResponse{
		ChallengeId: challengeID,
		Difficulty:  difficulty,
		Data:        data,
		Timestamp:   time.Now().Unix(),
	}, nil
}

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	challengepb.RegisterChallengeServiceServer(grpcServer, &challengeServer{})
	
	// Register reflection service for debugging
	reflection.Register(grpcServer)

	log.Printf("Challenge service listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
