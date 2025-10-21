package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ComputeHash computes the SHA-256 hash of the challenge data with the given nonce
func ComputeHash(challengeID, data string, nonce int64) string {
	input := fmt.Sprintf("%s:%s:%d", challengeID, data, nonce)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// VerifyProof verifies that the hash has the required number of leading zeros
func VerifyProof(hash string, difficulty int32) bool {
	if difficulty < 0 {
		return false
	}
	
	prefix := strings.Repeat("0", int(difficulty))
	return strings.HasPrefix(hash, prefix)
}

// FindNonce finds a nonce that produces a valid proof-of-work
// This is a brute-force search and is used for demonstration purposes
func FindNonce(challengeID, data string, difficulty int32, maxAttempts int64) (int64, string, bool) {
	for nonce := int64(0); nonce < maxAttempts; nonce++ {
		hash := ComputeHash(challengeID, data, nonce)
		if VerifyProof(hash, difficulty) {
			return nonce, hash, true
		}
	}
	return 0, "", false
}
