package pow

import (
	"testing"
)

func TestComputeHash(t *testing.T) {
	tests := []struct {
		name        string
		challengeID string
		data        string
		nonce       int64
	}{
		{
			name:        "simple hash",
			challengeID: "test-123",
			data:        "hello",
			nonce:       42,
		},
		{
			name:        "empty data",
			challengeID: "test-456",
			data:        "",
			nonce:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := ComputeHash(tt.challengeID, tt.data, tt.nonce)
			if len(hash) != 64 { // SHA-256 produces 64 hex characters
				t.Errorf("expected hash length 64, got %d", len(hash))
			}
			
			// Verify hash is deterministic
			hash2 := ComputeHash(tt.challengeID, tt.data, tt.nonce)
			if hash != hash2 {
				t.Errorf("hash should be deterministic")
			}
		})
	}
}

func TestVerifyProof(t *testing.T) {
	tests := []struct {
		name       string
		hash       string
		difficulty int32
		want       bool
	}{
		{
			name:       "valid proof difficulty 1",
			hash:       "0123456789abcdef",
			difficulty: 1,
			want:       true,
		},
		{
			name:       "valid proof difficulty 3",
			hash:       "000123456789abcdef",
			difficulty: 3,
			want:       true,
		},
		{
			name:       "invalid proof",
			hash:       "123456789abcdef",
			difficulty: 1,
			want:       false,
		},
		{
			name:       "difficulty too high",
			hash:       "00123456789abcdef",
			difficulty: 3,
			want:       false,
		},
		{
			name:       "negative difficulty",
			hash:       "0123456789abcdef",
			difficulty: -1,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifyProof(tt.hash, tt.difficulty)
			if got != tt.want {
				t.Errorf("VerifyProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindNonce(t *testing.T) {
	tests := []struct {
		name         string
		challengeID  string
		data         string
		difficulty   int32
		maxAttempts  int64
		expectFound  bool
	}{
		{
			name:         "find easy proof",
			challengeID:  "test-123",
			data:         "hello",
			difficulty:   2,
			maxAttempts:  10000,
			expectFound:  true,
		},
		{
			name:         "too few attempts",
			challengeID:  "test-456",
			data:         "world",
			difficulty:   4,
			maxAttempts:  10,
			expectFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, hash, found := FindNonce(tt.challengeID, tt.data, tt.difficulty, tt.maxAttempts)
			
			if found != tt.expectFound {
				t.Errorf("FindNonce() found = %v, want %v", found, tt.expectFound)
			}
			
			if found {
				// Verify the found nonce actually works
				computedHash := ComputeHash(tt.challengeID, tt.data, nonce)
				if computedHash != hash {
					t.Errorf("hash mismatch")
				}
				if !VerifyProof(hash, tt.difficulty) {
					t.Errorf("found nonce does not satisfy difficulty requirement")
				}
			}
		})
	}
}
