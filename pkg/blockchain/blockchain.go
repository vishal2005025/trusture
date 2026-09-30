package blockchain

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// ChainStats represents statistics about the blockchain
type ChainStats struct {
	TotalBlocks     int       `json:"total_blocks"`
	ValidatedBlocks int       `json:"validated_blocks"`
	ChainType       string    `json:"chain_type"`
	NGOID           string    `json:"ngo_id"`
	LastBlockTime   time.Time `json:"last_block_time"`
	IsValid         bool      `json:"is_valid"`
}

// Blockchain represents a blockchain for NGO transactions
type Blockchain struct {
	NGOID         string   `json:"ngo_id"`
	ChainType     string   `json:"chain_type"` // 'donation' or 'expenditure'
	Chain         []*Block `json:"chain"`
	Difficulty    int      `json:"difficulty"`
	PendingBlocks []*Block `json:"pending_blocks"`
	NetworkNodes  []string `json:"network_nodes"`
	mutex         sync.RWMutex
}

// NewBlockchain creates a new blockchain
func NewBlockchain(ngoID, chainType string, difficulty int) *Blockchain {
	if difficulty < 1 {
		difficulty = 2
	}

	if chainType != "donation" && chainType != "expenditure" {
		chainType = "donation"
	}

	blockchain := &Blockchain{
		NGOID:         ngoID,
		ChainType:     chainType,
		Difficulty:    difficulty,
		PendingBlocks: make([]*Block, 0),
		NetworkNodes:  make([]string, 0),
	}

	// Create and add genesis block
	genesisBlock := blockchain.createGenesisBlock()
	blockchain.Chain = []*Block{genesisBlock}

	return blockchain
}

// createGenesisBlock creates the first block in the chain
func (bc *Blockchain) createGenesisBlock() *Block {
	genesisData := map[string]interface{}{
		"type":       "genesis",
		"ngo_id":     bc.NGOID,
		"chain_type": bc.ChainType,
		"message":    fmt.Sprintf("Genesis block for %s chain of NGO %s", bc.ChainType, bc.NGOID),
	}

	genesisBlock := NewBlock(0, time.Now(), genesisData, "0", bc.ChainType)
	genesisBlock.Validated = true
	genesisBlock.AddValidator("system", "genesis_signature", "genesis")

	return genesisBlock
}

// GetLatestBlock returns the latest block in the chain
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if len(bc.Chain) == 0 {
		return nil
	}
	return bc.Chain[len(bc.Chain)-1]
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(newBlock *Block) bool {
	bc.mutex.RLock()
	if len(bc.Chain) == 0 {
		bc.mutex.RUnlock()
		return false
	}
	latestBlock := bc.Chain[len(bc.Chain)-1]
	bc.mutex.RUnlock()

	// If the block has an explicit previous hash, verify it matches
	if newBlock.PreviousHash != "" && newBlock.PreviousHash != latestBlock.Hash {
		log.Printf("Invalid previous hash: expected %s, got %s", latestBlock.Hash, newBlock.PreviousHash)
		return false
	}

	// Set index and previous hash if not already set
	if newBlock.PreviousHash == "" {
		newBlock.PreviousHash = latestBlock.Hash
	}
	if newBlock.Index == 0 {
		newBlock.Index = len(bc.Chain)
	}

	// Mine the block without holding the main lock
	newBlock.MineBlock(bc.Difficulty)

	// Acquire write lock to validate and append
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Re-check chain state before adding
	if len(bc.Chain) == 0 || bc.Chain[len(bc.Chain)-1].Hash != newBlock.PreviousHash {
		log.Println("Chain state changed while mining; rejecting block")
		return false
	}

	// Final validation before adding
	if bc.validateBlockInternal(newBlock) {
		newBlock.Validated = true
		newBlock.AddValidator("system", "auto", "mined")
		bc.Chain = append(bc.Chain, newBlock)
		return true
	}
	return false
}

// validateBlockInternal validates a block internally (assumes lock is held)
func (bc *Blockchain) validateBlockInternal(block *Block) bool {
	if len(bc.Chain) == 0 {
		return false
	}

	previousBlock := bc.Chain[len(bc.Chain)-1]

	// Check if previous hash is correct
	if block.PreviousHash != previousBlock.Hash {
		log.Printf("Invalid previous hash: expected %s, got %s", previousBlock.Hash, block.PreviousHash)
		return false
	}

	// Check if hash is correct
	if block.Hash != block.calculateHash() {
		log.Println("Invalid block hash")
		return false
	}

	// Check if block is properly mined
	target := strings.Repeat("0", bc.Difficulty)
	if !strings.HasPrefix(block.Hash, target) {
		log.Println("Block not properly mined")
		return false
	}

	return true
}

// ValidateBlock validates a block (thread-safe)
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.validateBlockInternal(block)
}

// IsChainValid checks if the entire blockchain is valid
func (bc *Blockchain) IsChainValid() bool {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		if !currentBlock.IsValid() {
			log.Printf("Block %d is invalid", i)
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash {
			log.Printf("Block %d has invalid previous hash", i)
			return false
		}
	}
	return true
}

// GetBlockByHash finds a block by its hash
func (bc *Blockchain) GetBlockByHash(hash string) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for _, block := range bc.Chain {
		if block.Hash == hash {
			return block
		}
	}
	return nil
}

// GetBlocksByDateRange returns blocks within a date range
func (bc *Blockchain) GetBlocksByDateRange(startDate, endDate time.Time) []*Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	var blocks []*Block
	for _, block := range bc.Chain {
		if (block.Timestamp.After(startDate) || block.Timestamp.Equal(startDate)) &&
			(block.Timestamp.Before(endDate) || block.Timestamp.Equal(endDate)) {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// GetChainStats returns statistics about the blockchain
func (bc *Blockchain) GetChainStats() ChainStats {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	validatedCount := 0
	var lastBlockTime time.Time

	for _, block := range bc.Chain {
		if block.Validated {
			validatedCount++
		}
		if block.Timestamp.After(lastBlockTime) {
			lastBlockTime = block.Timestamp
		}
	}

	return ChainStats{
		TotalBlocks:     len(bc.Chain),
		ValidatedBlocks: validatedCount,
		ChainType:       bc.ChainType,
		NGOID:           bc.NGOID,
		LastBlockTime:   lastBlockTime,
		IsValid:         bc.IsChainValid(),
	}
}

// GetBlockRange returns a range of blocks
func (bc *Blockchain) GetBlockRange(start, end int) []*Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if start < 0 {
		start = 0
	}
	if end >= len(bc.Chain) {
		end = len(bc.Chain) - 1
	}
	if start > end {
		return []*Block{}
	}

	blocks := make([]*Block, end-start+1)
	copy(blocks, bc.Chain[start:end+1])
	return blocks
}

// GetBlockByIndex returns a block by its index
func (bc *Blockchain) GetBlockByIndex(index int) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if index < 0 || index >= len(bc.Chain) {
		return nil
	}
	return bc.Chain[index]
}

// GetChainLength returns the length of the blockchain
func (bc *Blockchain) GetChainLength() int {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return len(bc.Chain)
}

// AddNetworkNode adds a network node
func (bc *Blockchain) AddNetworkNode(nodeAddress string) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	for _, node := range bc.NetworkNodes {
		if node == nodeAddress {
			return
		}
	}
	bc.NetworkNodes = append(bc.NetworkNodes, nodeAddress)
}

// RemoveNetworkNode removes a network node
func (bc *Blockchain) RemoveNetworkNode(nodeAddress string) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	newNodes := make([]string, 0)
	for _, node := range bc.NetworkNodes {
		if node != nodeAddress {
			newNodes = append(newNodes, node)
		}
	}
	bc.NetworkNodes = newNodes
}

// GetNetworkNodes returns the list of network nodes
func (bc *Blockchain) GetNetworkNodes() []string {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	nodes := make([]string, len(bc.NetworkNodes))
	copy(nodes, bc.NetworkNodes)
	return nodes
}

// AddPendingBlock adds a block to the pending list
func (bc *Blockchain) AddPendingBlock(block *Block) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	bc.PendingBlocks = append(bc.PendingBlocks, block)
}

// GetPendingBlocks returns all pending blocks
func (bc *Blockchain) GetPendingBlocks() []*Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	blocks := make([]*Block, len(bc.PendingBlocks))
	copy(blocks, bc.PendingBlocks)
	return blocks
}

// ClearPendingBlocks clears all pending blocks
func (bc *Blockchain) ClearPendingBlocks() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	bc.PendingBlocks = make([]*Block, 0)
}

// GetRecentBlocks returns the most recent blocks
func (bc *Blockchain) GetRecentBlocks(count int) []*Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if count <= 0 {
		return []*Block{}
	}

	chainLen := len(bc.Chain)
	if count > chainLen {
		count = chainLen
	}

	start := chainLen - count
	blocks := make([]*Block, count)
	copy(blocks, bc.Chain[start:])
	return blocks
}
