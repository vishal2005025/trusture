package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Signature represents a signature from a signer
type Signature struct {
	Signer    string    `json:"signer"`
	Signature string    `json:"signature"`
	Timestamp time.Time `json:"timestamp"`
}

// Transaction represents a pending multi-signature transaction
type Transaction struct {
	Data       map[string]interface{} `json:"data"`
	Signatures []Signature            `json:"signatures"`
	Executed   bool                   `json:"executed"`
	Timestamp  time.Time              `json:"timestamp"`
	Creator    string                 `json:"creator"`
}

// TransactionResult represents the result of a transaction operation
type TransactionResult struct {
	Success          bool   `json:"success"`
	Executed         bool   `json:"executed"`
	Transaction      *Transaction `json:"transaction,omitempty"`
	SignaturesCount  int    `json:"signatures_count,omitempty"`
	Message          string `json:"message"`
}

// TransactionStatus represents the status of a transaction
type TransactionStatus struct {
	TxID               string      `json:"tx_id"`
	Executed           bool        `json:"executed"`
	SignaturesCount    int         `json:"signatures_count"`
	RequiredSignatures int         `json:"required_signatures"`
	Signatures         []Signature `json:"signatures"`
}

// MultiSigWallet manages multi-signature transactions
type MultiSigWallet struct {
	RequiredSignatures   int                    `json:"required_signatures"`
	Signers              []string               `json:"signers"`
	PendingTransactions  map[string]*Transaction `json:"pending_transactions"`
	ExecutedTransactions map[string]bool        `json:"executed_transactions"`
	mutex                sync.RWMutex
}

// NewMultiSigWallet creates a new multi-signature wallet
func NewMultiSigWallet(requiredSignatures int) *MultiSigWallet {
	if requiredSignatures < 1 {
		requiredSignatures = 2
	}

	return &MultiSigWallet{
		RequiredSignatures:   requiredSignatures,
		Signers:              make([]string, 0),
		PendingTransactions:  make(map[string]*Transaction),
		ExecutedTransactions: make(map[string]bool),
	}
}

// AddSigner adds a new signer to the wallet
func (w *MultiSigWallet) AddSigner(address string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Check if signer already exists
	for _, signer := range w.Signers {
		if signer == address {
			return
		}
	}
	w.Signers = append(w.Signers, address)
}

// RemoveSigner removes a signer from the wallet
func (w *MultiSigWallet) RemoveSigner(address string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	newSigners := make([]string, 0)
	for _, signer := range w.Signers {
		if signer != address {
			newSigners = append(newSigners, signer)
		}
	}
	w.Signers = newSigners
}

// CreateTransaction creates a new pending transaction
func (w *MultiSigWallet) CreateTransaction(txData map[string]interface{}) (string, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Generate transaction ID
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate transaction ID: %v", err)
	}
	txID := hex.EncodeToString(randomBytes)

	creator := ""
	if creatorVal, ok := txData["creator"]; ok {
		if creatorStr, ok := creatorVal.(string); ok {
			creator = creatorStr
		}
	}

	transaction := &Transaction{
		Data:       txData,
		Signatures: make([]Signature, 0),
		Executed:   false,
		Timestamp:  time.Now(),
		Creator:    creator,
	}

	w.PendingTransactions[txID] = transaction
	return txID, nil
}

// SignTransaction adds a signature to a pending transaction
func (w *MultiSigWallet) SignTransaction(txID, signerAddress, signature string) *TransactionResult {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	tx, exists := w.PendingTransactions[txID]
	if !exists || tx.Executed {
		return &TransactionResult{
			Success: false,
			Message: "Transaction not found or already executed",
		}
	}

	// Check if signer is authorized
	authorized := false
	for _, signer := range w.Signers {
		if signer == signerAddress {
			authorized = true
			break
		}
	}
	if !authorized {
		return &TransactionResult{
			Success: false,
			Message: "Unauthorized signer",
		}
	}

	// Check if signer already signed
	for _, sig := range tx.Signatures {
		if sig.Signer == signerAddress {
			return &TransactionResult{
				Success: false,
				Message: "Already signed by this signer",
			}
		}
	}

	// Add signature
	tx.Signatures = append(tx.Signatures, Signature{
		Signer:    signerAddress,
		Signature: signature,
		Timestamp: time.Now(),
	})

	// Check if we have enough signatures to execute
	if len(tx.Signatures) >= w.RequiredSignatures {
		tx.Executed = true
		w.ExecutedTransactions[txID] = true
		return &TransactionResult{
			Success:     true,
			Executed:    true,
			Transaction: tx,
			Message:     "Transaction executed successfully",
		}
	}

	return &TransactionResult{
		Success:         true,
		Executed:        false,
		SignaturesCount: len(tx.Signatures),
		Message:         fmt.Sprintf("%d/%d signatures collected", len(tx.Signatures), w.RequiredSignatures),
	}
}

// GetTransactionStatus returns the status of a transaction
func (w *MultiSigWallet) GetTransactionStatus(txID string) *TransactionStatus {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	tx, exists := w.PendingTransactions[txID]
	if !exists {
		return nil
	}

	return &TransactionStatus{
		TxID:               txID,
		Executed:           tx.Executed,
		SignaturesCount:    len(tx.Signatures),
		RequiredSignatures: w.RequiredSignatures,
		Signatures:         tx.Signatures,
	}
}

// GetSigners returns the list of authorized signers
func (w *MultiSigWallet) GetSigners() []string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	signers := make([]string, len(w.Signers))
	copy(signers, w.Signers)
	return signers
}

// GetPendingTransactionCount returns the number of pending transactions
func (w *MultiSigWallet) GetPendingTransactionCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	count := 0
	for _, tx := range w.PendingTransactions {
		if !tx.Executed {
			count++
		}
	}
	return count
}

// GetExecutedTransactionCount returns the number of executed transactions
func (w *MultiSigWallet) GetExecutedTransactionCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return len(w.ExecutedTransactions)
}
