package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"ngo-transparency-platform/pkg/blockchain"
	"ngo-transparency-platform/pkg/crypto"
	"ngo-transparency-platform/pkg/transactions"
	"time"
)

// KYCData represents KYC information for NGO
type KYCData struct {
	Verified             bool      `json:"verified"`
	DocumentsHash        string    `json:"documents_hash"`
	VerificationDate     *time.Time `json:"verification_date"`
	VerificationAuthority string    `json:"verification_authority"`
}

// Certificate represents a certificate/certification
type Certificate struct {
	Type      string `json:"type"`
	Number    string `json:"number"`
	ValidUntil string `json:"valid_until"`
}

// ProcessResult represents the result of processing a transaction
type ProcessResult struct {
	Success       bool        `json:"success"`
	BlockHash     string      `json:"block_hash"`
	TransactionID string      `json:"transaction_id"`
	BlockIndex    int         `json:"block_index"`
	EBill         interface{} `json:"e_bill,omitempty"`
}

// RatingDetails represents detailed rating information
type RatingDetails struct {
	Rating              float64 `json:"rating"`
	TransparencyScore   int     `json:"transparency_score"`
	UtilizationRate     string  `json:"utilization_rate"`
	GapPercentage       string  `json:"gap_percentage"`
	TotalDonations      float64 `json:"total_donations"`
	TotalExpenditures   float64 `json:"total_expenditures"`
	PeriodDays          int     `json:"period_days"`
	DocumentationQuality string  `json:"documentation_quality"`
}

// FinancialSummary represents financial summary information
type FinancialSummary struct {
	Period             string             `json:"period"`
	TotalDonations     float64            `json:"total_donations"`
	TotalExpenditures  float64            `json:"total_expenditures"`
	DonationCount      int                `json:"donation_count"`
	ExpenditureCount   int                `json:"expenditure_count"`
	CategoryBreakdown  map[string]float64 `json:"category_breakdown"`
	AverageDonation    float64            `json:"average_donation"`
	MonthlyAverage     map[string]float64 `json:"monthly_average"`
}

// NGO represents a non-governmental organization
type NGO struct {
	NGOID                    string                       `json:"ngo_id"`
	Name                     string                       `json:"name"`
	RegistrationNumber       string                       `json:"registration_number"`
	Category                 string                       `json:"category"`
	Rating                   float64                      `json:"rating"`
	KYCData                  KYCData                      `json:"kyc_data"`
	DonationBlockchain       *blockchain.Blockchain       `json:"donation_blockchain"`
	ExpenditureBlockchain    *blockchain.Blockchain       `json:"expenditure_blockchain"`
	MultiSigWallet           *crypto.MultiSigWallet       `json:"multi_sig_wallet"`
	TotalDonationsReceived   float64                      `json:"total_donations_received"`
	TotalExpenditureReported float64                      `json:"total_expenditure_reported"`
	TransparencyScore        int                          `json:"transparency_score"`
	CreatedAt                time.Time                    `json:"created_at"`
	LastAuditDate            *time.Time                   `json:"last_audit_date"`
	Certificates             []Certificate                `json:"certificates"`
	PublicKey                string                       `json:"public_key"`
}

// NewNGO creates a new NGO instance
func NewNGO(ngoID, name, registrationNumber, category string, kycData map[string]interface{}, signers []string) *NGO {
	// Generate public key
	publicKeyData := ngoID + registrationNumber
	hash := sha256.Sum256([]byte(publicKeyData))
	publicKey := hex.EncodeToString(hash[:])

	// Create KYC data hash
	kycDocHash := ""
	if docs, ok := kycData["documents"]; ok {
		kycBytes := fmt.Sprintf("%v", docs)
		hash := sha256.Sum256([]byte(kycBytes))
		kycDocHash = hex.EncodeToString(hash[:])
	}

	ngo := &NGO{
		NGOID:              ngoID,
		Name:               name,
		RegistrationNumber: registrationNumber,
		Category:           category,
		Rating:             5.0,
		KYCData: KYCData{
			Verified:      false,
			DocumentsHash: kycDocHash,
		},
		DonationBlockchain:       blockchain.NewBlockchain(ngoID, "donation", 2),
		ExpenditureBlockchain:    blockchain.NewBlockchain(ngoID, "expenditure", 2),
		MultiSigWallet:           crypto.NewMultiSigWallet(2),
		TotalDonationsReceived:   0,
		TotalExpenditureReported: 0,
		TransparencyScore:        100,
		CreatedAt:                time.Now(),
		Certificates:             make([]Certificate, 0),
		PublicKey:                publicKey,
	}

	// Add multi-sig signers
	for _, signer := range signers {
		ngo.MultiSigWallet.AddSigner(signer)
	}

	return ngo
}

// VerifyKYC verifies the NGO's KYC data
func (ngo *NGO) VerifyKYC(authorityID string, certificates []Certificate) bool {
	ngo.KYCData.Verified = true
	now := time.Now()
	ngo.KYCData.VerificationDate = &now
	ngo.KYCData.VerificationAuthority = authorityID
	ngo.Certificates = certificates
	
	return ngo.KYCData.Verified
}

// ProcessDonation processes a donation transaction
func (ngo *NGO) ProcessDonation(donation *transactions.DonationTransaction) (*ProcessResult, error) {
	if !donation.ValidateEBill() {
		return nil, fmt.Errorf("invalid e-bill")
	}

	// Verify ZK proof
	if !crypto.VerifyProof(donation.ZKProof, donation.Amount, donation.Timestamp) {
		return nil, fmt.Errorf("invalid zero-knowledge proof")
	}

	// Create block data
	blockData := map[string]interface{}{
		"type":           "donation",
		"transaction_id": donation.TransactionID,
		"donor_hash":     ngo.generateDonorHash(donation.DonorID),
		"amount":         donation.Amount,
		"currency":       "INR",
		"zk_proof":       donation.ZKProof,
		"e_bill":         donation.EBill,
		"timestamp":      donation.Timestamp,
		"payment_method": donation.PaymentMethod,
	}

	block := blockchain.NewBlock(
		ngo.DonationBlockchain.GetChainLength(),
		time.Now(),
		blockData,
		ngo.DonationBlockchain.GetLatestBlock().Hash,
		"donation",
	)

	// Validate block with e-bill
	block.Validate()
	block.AddValidator("ebill_system", donation.EBill.Signature, "ebill")
	block.AddValidator("zk_system", donation.ZKProof.Proof, "zkproof")

	if ngo.DonationBlockchain.AddBlock(block) {
		ngo.TotalDonationsReceived += donation.Amount
		donation.MarkComplete()

		return &ProcessResult{
			Success:       true,
			BlockHash:     block.Hash,
			TransactionID: donation.TransactionID,
			BlockIndex:    block.Index,
			EBill:         donation.EBill,
		}, nil
	} else {
		donation.MarkFailed("Block validation failed")
		return nil, fmt.Errorf("failed to add block to blockchain")
	}
}

// ProcessExpenditure processes an expenditure transaction
func (ngo *NGO) ProcessExpenditure(expenditure *transactions.ExpenditureTransaction) (*ProcessResult, error) {
	if expenditure.AuditorValidation == nil || !expenditure.AuditorValidation.IsValid {
		return nil, fmt.Errorf("expenditure not validated by auditor")
	}

	if !expenditure.VerifyGSTIN(expenditure.InvoiceDetails.GSTIN) {
		return nil, fmt.Errorf("invalid GSTIN format")
	}

	// Check compliance score
	if expenditure.ComplianceScore < 60 {
		return nil, fmt.Errorf("low compliance score: %.1f%%. Minimum required: 60%%", expenditure.ComplianceScore)
	}

	// Create block data
	blockData := map[string]interface{}{
		"type":               "expenditure",
		"transaction_id":     expenditure.TransactionID,
		"amount":             expenditure.Amount,
		"currency":           "INR",
		"category":           expenditure.Category,
		"description":        expenditure.Description,
		"invoice_details":    expenditure.InvoiceDetails,
		"auditor_validation": expenditure.AuditorValidation,
		"compliance_score":   expenditure.ComplianceScore,
		"timestamp":          expenditure.Timestamp,
		"attachments":        ngo.extractAttachmentHashes(expenditure.Attachments),
	}

	block := blockchain.NewBlock(
		ngo.ExpenditureBlockchain.GetChainLength(),
		time.Now(),
		blockData,
		ngo.ExpenditureBlockchain.GetLatestBlock().Hash,
		"expenditure",
	)

	// Validate block with auditor signature
	block.Validate()
	block.AddValidator(
		expenditure.AuditorValidation.AuditorID,
		expenditure.AuditorValidation.Signature,
		"auditor",
	)

	if ngo.ExpenditureBlockchain.AddBlock(block) {
		ngo.TotalExpenditureReported += expenditure.Amount

		return &ProcessResult{
			Success:       true,
			BlockHash:     block.Hash,
			TransactionID: expenditure.TransactionID,
			BlockIndex:    block.Index,
		}, nil
	} else {
		return nil, fmt.Errorf("failed to add expenditure block to blockchain")
	}
}

// CalculateRating calculates the NGO's rating based on recent activity
func (ngo *NGO) CalculateRating(periodDays int) RatingDetails {
	periodMs := time.Duration(periodDays) * 24 * time.Hour
	startDate := time.Now().Add(-periodMs)

	donations := ngo.DonationBlockchain.GetBlocksByDateRange(startDate, time.Now())
	expenditures := ngo.ExpenditureBlockchain.GetBlocksByDateRange(startDate, time.Now())

	totalDonations := ngo.sumBlockAmounts(donations)
	totalExpenditures := ngo.sumBlockAmounts(expenditures)

	utilizationRate := 0.0
	if totalDonations > 0 {
		utilizationRate = totalExpenditures / totalDonations
	}

	gap := math.Abs(totalDonations - totalExpenditures)
	gapPercentage := 0.0
	if totalDonations > 0 {
		gapPercentage = (gap / totalDonations) * 100
	}

	// Rating calculation (1.0 to 5.0)
	rating := 5.0

	// Penalize large gaps
	if gapPercentage > 50 {
		rating -= 2.0
	} else if gapPercentage > 30 {
		rating -= 1.0
	} else if gapPercentage > 15 {
		rating -= 0.5
	}

	// Reward optimal utilization (60-85% is ideal)
	if utilizationRate >= 0.6 && utilizationRate <= 0.85 {
		rating += 0.5
	} else if utilizationRate < 0.3 || utilizationRate > 0.95 {
		rating -= 0.5
	}

	// Transparency and documentation bonus
	documentationQuality := ngo.calculateDocumentationQuality()
	rating += documentationQuality * 0.5

	// KYC verification bonus
	if ngo.KYCData.Verified {
		rating += 0.2
	}

	// Certificate bonus
	if len(ngo.Certificates) > 0 {
		rating += 0.3
	}

	rating = math.Max(1.0, math.Min(5.0, rating))
	ngo.Rating = rating
	ngo.TransparencyScore = int(math.Round((rating / 5.0) * 100))

	return RatingDetails{
		Rating:               rating,
		TransparencyScore:    ngo.TransparencyScore,
		UtilizationRate:      fmt.Sprintf("%.2f%%", utilizationRate*100),
		GapPercentage:        fmt.Sprintf("%.2f%%", gapPercentage),
		TotalDonations:       totalDonations,
		TotalExpenditures:    totalExpenditures,
		PeriodDays:          periodDays,
		DocumentationQuality: fmt.Sprintf("%.1f%%", documentationQuality*100),
	}
}

// GetBlockchainStats returns blockchain statistics
func (ngo *NGO) GetBlockchainStats() map[string]interface{} {
	return map[string]interface{}{
		"ngo_id":                       ngo.NGOID,
		"name":                         ngo.Name,
		"category":                     ngo.Category,
		"rating":                       ngo.Rating,
		"transparency_score":           ngo.TransparencyScore,
		"kyc_verified":                 ngo.KYCData.Verified,
		"total_donations_received":     ngo.TotalDonationsReceived,
		"total_expenditure_reported":   ngo.TotalExpenditureReported,
		"donation_blockchain_length":   ngo.DonationBlockchain.GetChainLength(),
		"expenditure_blockchain_length": ngo.ExpenditureBlockchain.GetChainLength(),
		"donation_chain_valid":         ngo.DonationBlockchain.IsChainValid(),
		"expenditure_chain_valid":      ngo.ExpenditureBlockchain.IsChainValid(),
		"last_audit_date":              ngo.LastAuditDate,
		"certificates_count":           len(ngo.Certificates),
		"created_at":                   ngo.CreatedAt,
	}
}

// GetFinancialSummary returns financial summary for the specified months
func (ngo *NGO) GetFinancialSummary(months int) FinancialSummary {
	monthMs := time.Duration(months) * 30 * 24 * time.Hour
	startDate := time.Now().Add(-monthMs)

	donations := ngo.DonationBlockchain.GetBlocksByDateRange(startDate, time.Now())
	expenditures := ngo.ExpenditureBlockchain.GetBlocksByDateRange(startDate, time.Now())

	// Category-wise expenditure breakdown
	categoryBreakdown := make(map[string]float64)
	for _, block := range expenditures {
		if blockData, ok := block.Data.(map[string]interface{}); ok {
			if category, ok := blockData["category"].(string); ok {
				if amount, ok := blockData["amount"].(float64); ok {
					categoryBreakdown[category] += amount
				}
			}
		}
	}

	totalDonations := ngo.sumBlockAmounts(donations)
	totalExpenditures := ngo.sumBlockAmounts(expenditures)
	
	averageDonation := 0.0
	if len(donations) > 0 {
		averageDonation = totalDonations / float64(len(donations))
	}

	monthlyAverage := map[string]float64{
		"donations":    totalDonations / float64(months),
		"expenditures": totalExpenditures / float64(months),
	}

	return FinancialSummary{
		Period:            fmt.Sprintf("%d months", months),
		TotalDonations:    totalDonations,
		TotalExpenditures: totalExpenditures,
		DonationCount:     len(donations),
		ExpenditureCount:  len(expenditures),
		CategoryBreakdown: categoryBreakdown,
		AverageDonation:   averageDonation,
		MonthlyAverage:    monthlyAverage,
	}
}

// Helper methods

func (ngo *NGO) generateDonorHash(donorID string) string {
	hash := sha256.Sum256([]byte(donorID))
	return hex.EncodeToString(hash[:])
}

func (ngo *NGO) extractAttachmentHashes(attachments []transactions.Attachment) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attachments))
	for i, att := range attachments {
		result[i] = map[string]interface{}{
			"filename": att.Filename,
			"hash":     att.Hash,
			"type":     att.Type,
		}
	}
	return result
}

func (ngo *NGO) sumBlockAmounts(blocks []*blockchain.Block) float64 {
	total := 0.0
	for _, block := range blocks {
		if blockData, ok := block.Data.(map[string]interface{}); ok {
			if amount, ok := blockData["amount"].(float64); ok {
				total += amount
			}
		}
	}
	return total
}

func (ngo *NGO) calculateDocumentationQuality() float64 {
	recentBlocks := ngo.ExpenditureBlockchain.GetRecentBlocks(10)
	if len(recentBlocks) == 0 {
		return 1.0 // Perfect score if no expenditures yet
	}

	score := 0.0
	validBlocks := 0

	for _, block := range recentBlocks {
		if blockData, ok := block.Data.(map[string]interface{}); ok {
			if blockType, ok := blockData["type"].(string); ok && blockType == "expenditure" {
				validBlocks++
				blockScore := 0.0

				// Basic invoice details (40%)
				if invoiceDetails, ok := blockData["invoice_details"].(map[string]interface{}); ok {
					if invoiceNum, ok := invoiceDetails["invoice_number"].(string); ok && invoiceNum != "" {
						blockScore += 0.1
					}
					if gstin, ok := invoiceDetails["gstin"].(string); ok && gstin != "" {
						blockScore += 0.1
					}
					if vendorName, ok := invoiceDetails["vendor_name"].(string); ok && vendorName != "" {
						blockScore += 0.1
					}
					if vendorGSTIN, ok := invoiceDetails["vendor_gstin"].(string); ok && vendorGSTIN != "" {
						blockScore += 0.1
					}
				}

				// Supporting documents (30%)
				if invoiceDetails, ok := blockData["invoice_details"].(map[string]interface{}); ok {
					if docs, ok := invoiceDetails["documents"].([]interface{}); ok && len(docs) > 0 {
						blockScore += 0.2
					}
				}
				if attachments, ok := blockData["attachments"].([]interface{}); ok && len(attachments) > 0 {
					blockScore += 0.1
				}

				// Validation and compliance (30%)
				if block.Validated {
					blockScore += 0.1
				}
				if len(block.Validators) > 0 {
					blockScore += 0.1
				}
				if complianceScore, ok := blockData["compliance_score"].(float64); ok && complianceScore >= 80 {
					blockScore += 0.1
				}

				score += blockScore
			}
		}
	}

	if validBlocks == 0 {
		return 1.0
	}
	return math.Min(1.0, score/float64(validBlocks))
}
