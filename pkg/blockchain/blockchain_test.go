package blockchain

import (
	"testing"
	"time"
)

func TestBlockchainCreation(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)
	if bc == nil {
		t.Fatal("Expected non-nil blockchain")
	}

	if bc.GetChainLength() != 1 {
		t.Errorf("Expected genesis block, got chain length %d", bc.GetChainLength())
	}

	genesis := bc.GetBlockByIndex(0)
	if genesis.PreviousHash != "0" {
		t.Errorf("Expected genesis previous hash '0', got %s", genesis.PreviousHash)
	}

	if !genesis.IsValid() {
		t.Error("Expected genesis block to be valid")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)

	data := map[string]interface{}{
		"amount":    1.5,
		"donor":     "0x123",
		"timestamp": time.Now(),
	}

	// Add first block
	newBlock := NewBlock(1, time.Now(), data, bc.GetLatestBlock().Hash, "donation")
	if !bc.AddBlock(newBlock) {
		t.Fatal("Failed to add valid block")
	}

	// Verify chain length
	if bc.GetChainLength() != 2 {
		t.Errorf("Expected chain length 2, got %d", bc.GetChainLength())
	}

	// Verify block is properly linked
	if newBlock.PreviousHash != bc.GetBlockByIndex(0).Hash {
		t.Error("Block not properly linked to previous block")
	}

	// Try to add block with invalid previous hash
	invalidBlock := NewBlock(2, time.Now(), data, "invalid_hash", "donation")
	if bc.AddBlock(invalidBlock) {
		t.Error("Should not add block with invalid previous hash")
	}
}

func TestConcurrentMining(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)
	done := make(chan bool)

	// Start multiple goroutines trying to add blocks simultaneously
	for i := 0; i < 3; i++ {
		go func(index int) {
			data := map[string]interface{}{
				"amount":    float64(index),
				"donor":     "0x123",
				"timestamp": time.Now(),
			}
			newBlock := NewBlock(1, time.Now(), data, bc.GetLatestBlock().Hash, "donation")
			bc.AddBlock(newBlock)
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify chain is valid
	if !bc.IsChainValid() {
		t.Error("Chain should remain valid after concurrent mining")
	}

	// Verify chain length (should be genesis + 1 successful block)
	expectedLength := 2 // genesis + 1
	if bc.GetChainLength() != expectedLength {
		t.Errorf("Expected chain length %d, got %d", expectedLength, bc.GetChainLength())
	}
}

func TestBlockchainValidation(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)

	// Add some valid blocks
	for i := 0; i < 3; i++ {
		data := map[string]interface{}{
			"amount":    float64(i),
			"donor":     "0x123",
			"timestamp": time.Now(),
		}
		newBlock := NewBlock(bc.GetChainLength(), time.Now(), data, bc.GetLatestBlock().Hash, "donation")
		if !bc.AddBlock(newBlock) {
			t.Fatalf("Failed to add block %d", i)
		}
	}

	// Verify entire chain
	if !bc.IsChainValid() {
		t.Error("Chain should be valid")
	}

	// Try to tamper with a block's data
	middleBlock := bc.GetBlockByIndex(2)
	tamperedData := map[string]interface{}{
		"amount": 999.9,
		"donor":  "hacker",
	}
	middleBlock.Data = tamperedData

	// Chain should now be invalid
	if bc.IsChainValid() {
		t.Error("Chain should be invalid after tampering")
	}
}

func TestBlockChainValidators(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)

	data := map[string]interface{}{
		"amount": 1.5,
		"donor":  "0x123",
	}

	newBlock := NewBlock(1, time.Now(), data, bc.GetLatestBlock().Hash, "donation")

	// Add validators before mining
	newBlock.AddValidator("auditor1", "sig1", "audit")
	newBlock.AddValidator("system", "sig2", "auto")

	if !bc.AddBlock(newBlock) {
		t.Fatal("Failed to add block with validators")
	}

	// Verify validators were preserved
	addedBlock := bc.GetBlockByIndex(1)
	if len(addedBlock.Validators) != 3 { // 2 added + 1 system validator from AddBlock
		t.Errorf("Expected 3 validators, got %d", len(addedBlock.Validators))
	}

	// Verify validator types
	auditValidators := addedBlock.GetValidatorsByType("audit")
	if len(auditValidators) != 1 {
		t.Error("Expected 1 audit validator")
	}

	systemValidators := addedBlock.GetValidatorsByType("auto")
	if len(systemValidators) != 1 {
		t.Error("Expected 1 system validator")
	}
}

func TestBlockchainStats(t *testing.T) {
	bc := NewBlockchain("NGO001", "donation", 2)

	// Add some blocks
	for i := 0; i < 3; i++ {
		data := map[string]interface{}{
			"amount":    float64(i),
			"timestamp": time.Now(),
		}
		newBlock := NewBlock(bc.GetChainLength(), time.Now(), data, bc.GetLatestBlock().Hash, "donation")
		if !bc.AddBlock(newBlock) {
			t.Fatalf("Failed to add block %d", i)
		}
	}

	// Get chain stats
	stats := bc.GetChainStats()

	if stats.TotalBlocks != 4 { // genesis + 3 added
		t.Errorf("Expected 4 total blocks, got %d", stats.TotalBlocks)
	}

	if stats.ValidatedBlocks != 4 { // all should be validated
		t.Errorf("Expected 4 validated blocks, got %d", stats.ValidatedBlocks)
	}

	if !stats.IsValid {
		t.Error("Chain should be valid")
	}

	if stats.NGOID != "NGO001" {
		t.Errorf("Expected NGO001, got %s", stats.NGOID)
	}

	if stats.ChainType != "donation" {
		t.Errorf("Expected donation chain type, got %s", stats.ChainType)
	}
}
