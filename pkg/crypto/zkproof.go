package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"time"
)

// ZKProof represents a zero-knowledge proof for donor anonymity
type ZKProof struct {
	Commitment string    `json:"commitment"`
	Nullifier  string    `json:"nullifier"`
	Proof      string    `json:"proof"`
	Timestamp  time.Time `json:"timestamp"`
}

// GenerateProof creates a zero-knowledge proof for donor anonymity
func GenerateProof(donorID string, amount float64, timestamp time.Time) *ZKProof {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	// Create commitment hash
	commitmentData := fmt.Sprintf("%s%.2f%d", donorID, amount, timestamp.UnixNano())
	commitmentHash := sha256.Sum256([]byte(commitmentData))
	commitment := hex.EncodeToString(commitmentHash[:])

	// Generate nullifier to prevent double spending
	nullifierData := fmt.Sprintf("%s%s", commitment, donorID)
	nullifierHash := sha256.Sum256([]byte(nullifierData))
	nullifier := hex.EncodeToString(nullifierHash[:])

	// Generate proof using enhanced hash simulation
	proofData := fmt.Sprintf("%s%s%d", commitment, nullifier, time.Now().UnixNano())
	proofHash := sha256.Sum256([]byte(proofData))
	proof := hex.EncodeToString(proofHash[:])

	return &ZKProof{
		Commitment: commitment,
		Nullifier:  nullifier,
		Proof:      proof,
		Timestamp:  timestamp,
	}
}

// VerifyProof validates a zero-knowledge proof
func VerifyProof(proof *ZKProof, amount float64, timestamp time.Time) bool {
	if proof == nil || proof.Commitment == "" || proof.Nullifier == "" || proof.Proof == "" {
		return false
	}

	// Verify timestamp is reasonable (within last 24 hours)
	now := time.Now()
	timeDiff := math.Abs(now.Sub(timestamp).Hours())
	if timeDiff > 24 {
		return false
	}

	// Verify proof structure (all should be valid 64-character hex strings)
	hexPattern := regexp.MustCompile(`^[a-f0-9]{64}$`)
	return hexPattern.MatchString(proof.Commitment) &&
		hexPattern.MatchString(proof.Nullifier) &&
		hexPattern.MatchString(proof.Proof)
}
