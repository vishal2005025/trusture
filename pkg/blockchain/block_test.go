package blockchain

import (
	"testing"
	"time"
)

func TestNewBlock(t *testing.T) {
	data := map[string]interface{}{
		"amount": 1.5,
		"donor":  "0x123",
	}
	block := NewBlock(1, time.Now(), data, "prev_hash", "donation")

	if block == nil {
		t.Fatal("Expected non-nil block")
	}

	if block.Index != 1 {
		t.Errorf("Expected index 1, got %d", block.Index)
	}

	if block.PreviousHash != "prev_hash" {
		t.Errorf("Expected prev_hash, got %s", block.PreviousHash)
	}

	if block.BlockType != "donation" {
		t.Errorf("Expected donation type, got %s", block.BlockType)
	}

	if block.Hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestBlockValidation(t *testing.T) {
	data := map[string]interface{}{
		"amount": 1.5,
		"donor":  "0x123",
	}
	block := NewBlock(1, time.Now(), data, "prev_hash", "donation")

	// Block should be valid initially
	if !block.IsValid() {
		t.Error("New block should be valid")
	}

	// Test with tampered data
	// Save original hash
	originalHash := block.Hash

	// Tamper with data
	block.Data = map[string]interface{}{
		"amount": 999.9,
		"donor":  "hacker",
	}

	// Block should be invalid after tampering (hash won't match anymore)
	if block.IsValid() {
		t.Error("Block should be invalid after tampering")
	}

	// Verify hash changed
	currentHash := block.calculateHash()
	if currentHash == originalHash {
		t.Error("Block hash should change after data modification")
	}

	// Update hash and validate again
	block.UpdateHash()
	if !block.IsValid() {
		t.Error("Block should be valid after hash update")
	}
}

func TestBlockValidators(t *testing.T) {
	block := NewBlock(1, time.Now(), nil, "prev_hash", "donation")

	// Add validators
	block.AddValidator("auditor1", "sig1", "audit")
	block.AddValidator("system", "sig2", "auto")

	if len(block.Validators) != 2 {
		t.Errorf("Expected 2 validators, got %d", len(block.Validators))
	}

	// Test getting validators by type
	auditValidators := block.GetValidatorsByType("audit")
	if len(auditValidators) != 1 {
		t.Error("Expected 1 audit validator")
	}

	systemValidators := block.GetValidatorsByType("auto")
	if len(systemValidators) != 1 {
		t.Error("Expected 1 system validator")
	}

	// Verify validator details
	if auditValidators[0].ValidatorID != "auditor1" || auditValidators[0].Signature != "sig1" {
		t.Error("Validator details don't match")
	}
}

func TestBlockTimestamp(t *testing.T) {
	now := time.Now()
	block := NewBlock(1, now, nil, "prev_hash", "donation")

	if !block.Timestamp.Equal(now) {
		t.Error("Block timestamp doesn't match creation time")
	}

	// Test future timestamp
	futureBlock := NewBlock(2, now.Add(24*time.Hour), nil, block.Hash, "donation")
	if futureBlock.IsValid() {
		t.Error("Block with future timestamp should be invalid")
	}
}

func TestBlockHash(t *testing.T) {
	data := map[string]interface{}{
		"amount": 1.5,
		"donor":  "0x123",
	}
	block1 := NewBlock(1, time.Now(), data, "prev_hash", "donation")
	block2 := NewBlock(1, time.Now(), data, "prev_hash", "donation")

	// Even with same data, hashes should be different due to timestamp
	if block1.Hash == block2.Hash {
		t.Error("Different blocks should have different hashes")
	}

	// Store original hash
	originalHash := block1.Hash

	// Modify data
	block1.Data = map[string]interface{}{
		"amount": 2.0,
		"donor":  "0x123",
	}

	// Update hash and verify it changed
	block1.UpdateHash()
	if block1.Hash == originalHash {
		t.Error("Hash should change when data is modified")
	}

	// Block should be valid after hash update
	if !block1.IsValid() {
		t.Error("Block should be valid after hash update")
	}
}
