package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Validator represents a validator for a block
type Validator struct {
	ValidatorID    string    `json:"validator_id"`
	Signature      string    `json:"signature"`
	ValidationType string    `json:"validation_type"`
	Timestamp      time.Time `json:"timestamp"`
}

// Block represents a single block in the blockchain
type Block struct {
	Index        int         `json:"index"`
	Timestamp    time.Time   `json:"timestamp"`
	Data         interface{} `json:"data"`
	PreviousHash string      `json:"previous_hash"`
	BlockType    string      `json:"block_type"` // 'donation' or 'expenditure'
	Hash         string      `json:"hash"`
	Nonce        int         `json:"nonce"`
	Validated    bool        `json:"validated"`
	Validators   []Validator `json:"validators"`
	MerkleRoot   string      `json:"merkle_root"`
}

// NewBlock creates a new block
func NewBlock(index int, timestamp time.Time, data interface{}, previousHash, blockType string) *Block {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	if blockType == "" {
		blockType = "donation"
	}

	block := &Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
		BlockType:    blockType,
		Nonce:        0,
		Validated:    true, // Initial blocks are valid until tampered
		Validators:   make([]Validator, 0),
	}

	block.MerkleRoot = block.calculateMerkleRoot()
	block.Hash = block.calculateHash()

	return block
}

// calculateHash computes the hash of the block
func (b *Block) calculateHash() string {
	dataBytes, err := json.Marshal(b.Data)
	if err != nil {
		dataBytes = []byte(fmt.Sprintf("%v", b.Data))
	}

	record := fmt.Sprintf("%d%s%d%s%d%s%s",
		b.Index,
		b.PreviousHash,
		b.Timestamp.UnixNano(),
		string(dataBytes),
		b.Nonce,
		b.BlockType,
		b.MerkleRoot,
	)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// calculateMerkleRoot computes a simplified Merkle root
func (b *Block) calculateMerkleRoot() string {
	dataBytes, err := json.Marshal(b.Data)
	if err != nil {
		dataBytes = []byte(fmt.Sprintf("%v", b.Data))
	}

	hash := sha256.Sum256(dataBytes)
	return hex.EncodeToString(hash[:])
}

// MineBlock performs proof-of-work mining on the block
func (b *Block) MineBlock(difficulty int) {
	target := strings.Repeat("0", difficulty)

	for !strings.HasPrefix(b.Hash, target) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}

	fmt.Printf("Block mined: %s\n", b.Hash)
}

// AddValidator adds a validator to the block
func (b *Block) AddValidator(validatorID, signature, validationType string) {
	if validationType == "" {
		validationType = "general"
	}

	validator := Validator{
		ValidatorID:    validatorID,
		Signature:      signature,
		ValidationType: validationType,
		Timestamp:      time.Now(),
	}

	b.Validators = append(b.Validators, validator)
}

// IsValid checks if the block is valid
func (b *Block) IsValid() bool {
	// Regenerate hash and check if it matches
	currentHash := b.calculateHash()

	// Check timestamp is not in the future (allow 1 minute tolerance)
	if b.Timestamp.After(time.Now().Add(time.Minute)) {
		return false
	}

	// Verify hash and validation status
	return b.Hash == currentHash && b.Validated
}

// UpdateHash recalculates and updates the block's hash after data changes
func (b *Block) UpdateHash() {
	b.Hash = b.calculateHash()
}

// GetValidatorsByType returns validators of a specific type
func (b *Block) GetValidatorsByType(validationType string) []Validator {
	var validators []Validator
	for _, validator := range b.Validators {
		if validator.ValidationType == validationType {
			validators = append(validators, validator)
		}
	}
	return validators
}

// GetValidatorsCount returns the number of validators
func (b *Block) GetValidatorsCount() int {
	return len(b.Validators)
}

// Validate marks the block as validated
func (b *Block) Validate() {
	b.Validated = true
}

// GetBlockInfo returns basic information about the block
func (b *Block) GetBlockInfo() map[string]interface{} {
	return map[string]interface{}{
		"index":         b.Index,
		"timestamp":     b.Timestamp,
		"hash":          b.Hash,
		"previous_hash": b.PreviousHash,
		"block_type":    b.BlockType,
		"validated":     b.Validated,
		"validators":    len(b.Validators),
		"nonce":         b.Nonce,
		"merkle_root":   b.MerkleRoot,
	}
}
