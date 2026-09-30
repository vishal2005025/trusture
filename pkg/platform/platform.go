package platform

import (
	"fmt"
	"math/big"
	"ngo-transparency-platform/pkg/entities"
	"ngo-transparency-platform/pkg/polygon"
	"ngo-transparency-platform/pkg/transactions"
	"sort"
	"sync"
	"time"
)

// SystemStats represents platform-wide statistics
type SystemStats struct {
	TotalTransactions int       `json:"total_transactions"`
	TotalDonations    float64   `json:"total_donations"`
	TotalExpenditures float64   `json:"total_expenditures"`
	PlatformFee       float64   `json:"platform_fee"`
	CreatedAt         time.Time `json:"created_at"`
}

// PlatformStats represents comprehensive platform statistics
type PlatformStats struct {
	TotalNGOs            int      `json:"total_ngos"`
	TotalDonors          int      `json:"total_donors"`
	TotalAuditors        int      `json:"total_auditors"`
	TotalTransactions    int      `json:"total_transactions"`
	TotalDonations       float64  `json:"total_donations"`
	TotalExpenditures    float64  `json:"total_expenditures"`
	PlatformFeeCollected float64  `json:"platform_fee_collected"`
	VerifiedNGOs         int      `json:"verified_ngos"`
	VerifiedDonors       int      `json:"verified_donors"`
	VerifiedAuditors     int      `json:"verified_auditors"`
	KYCAuthorities       int      `json:"kyc_authorities"`
	DaysActive           int      `json:"days_active"`
	AverageNGORating     float64  `json:"average_ngo_rating"`
	Categories           []string `json:"categories"`
}

// NGOTransparencyPlatform is the main platform orchestrator
type NGOTransparencyPlatform struct {
	NGOs               map[string]*entities.NGO     `json:"ngos"`
	Donors             map[string]*entities.Donor   `json:"donors"`
	Auditors           map[string]*entities.Auditor `json:"auditors"`
	PolygonIntegration *polygon.PolygonIntegration  `json:"polygon_integration"`
	SystemStats        SystemStats                  `json:"system_stats"`
	KYCAuthorities     map[string]bool              `json:"kyc_authorities"`
	mutex              sync.RWMutex
}

// NewNGOTransparencyPlatform creates a new platform instance
func NewNGOTransparencyPlatform() *NGOTransparencyPlatform {
	return &NGOTransparencyPlatform{
		NGOs:           make(map[string]*entities.NGO),
		Donors:         make(map[string]*entities.Donor),
		Auditors:       make(map[string]*entities.Auditor),
		KYCAuthorities: make(map[string]bool),
		SystemStats: SystemStats{
			TotalTransactions: 0,
			TotalDonations:    0,
			TotalExpenditures: 0,
			PlatformFee:       0.01, // 1% platform fee
			CreatedAt:         time.Now(),
		},
	}
}

// InitializePolygon initializes Polygon blockchain integration
func (p *NGOTransparencyPlatform) InitializePolygon(providerURL, privateKey string, gasLimit int64, gasPrice *big.Int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.PolygonIntegration = polygon.NewPolygonIntegration(providerURL, privateKey, gasLimit, gasPrice)
}

// RegisterNGO registers a new NGO on the platform
func (p *NGOTransparencyPlatform) RegisterNGO(ngoID, name, registrationNumber, category string, kycData map[string]interface{}, signers []string) (*entities.NGO, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.NGOs[ngoID]; exists {
		return nil, fmt.Errorf("NGO already registered")
	}

	ngo := entities.NewNGO(ngoID, name, registrationNumber, category, kycData, signers)
	p.NGOs[ngoID] = ngo

	return ngo, nil
}

// VerifyNGOKYC verifies NGO's KYC
func (p *NGOTransparencyPlatform) VerifyNGOKYC(ngoID, authorityID string, certificates []entities.Certificate) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ngo, exists := p.NGOs[ngoID]
	if !exists {
		return fmt.Errorf("NGO not found")
	}

	p.KYCAuthorities[authorityID] = true
	ngo.VerifyKYC(authorityID, certificates)
	return nil
}

// RegisterDonor registers a new donor on the platform
func (p *NGOTransparencyPlatform) RegisterDonor(donorID string, kycData map[string]interface{}) (*entities.Donor, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.Donors[donorID]; exists {
		return nil, fmt.Errorf("donor already registered")
	}

	donor := entities.NewDonor(donorID, kycData)
	p.Donors[donorID] = donor

	return donor, nil
}

// VerifyDonorKYC verifies donor's KYC
func (p *NGOTransparencyPlatform) VerifyDonorKYC(donorID, authorityID, verificationLevel string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	donor, exists := p.Donors[donorID]
	if !exists {
		return fmt.Errorf("donor not found")
	}

	p.KYCAuthorities[authorityID] = true
	donor.VerifyKYC(authorityID, verificationLevel)
	return nil
}

// RegisterAuditor registers a new auditor on the platform
func (p *NGOTransparencyPlatform) RegisterAuditor(auditorID, name string, credentials interface{}, specializations []string) (*entities.Auditor, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.Auditors[auditorID]; exists {
		return nil, fmt.Errorf("auditor already registered")
	}

	auditor := entities.NewAuditor(auditorID, name, credentials, specializations)
	p.Auditors[auditorID] = auditor

	return auditor, nil
}

// VerifyAuditorCredentials verifies auditor's credentials
func (p *NGOTransparencyPlatform) VerifyAuditorCredentials(auditorID, verificationAuthority string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	auditor, exists := p.Auditors[auditorID]
	if !exists {
		return fmt.Errorf("auditor not found")
	}

	auditor.VerifyCredentials(verificationAuthority)
	return nil
}

// ProcessDonation processes a donation transaction
func (p *NGOTransparencyPlatform) ProcessDonation(donorID, ngoID string, amount float64, paymentMethod string) (map[string]interface{}, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	donor, donorExists := p.Donors[donorID]
	ngo, ngoExists := p.NGOs[ngoID]

	if !donorExists {
		return nil, fmt.Errorf("donor not found")
	}
	if !ngoExists {
		return nil, fmt.Errorf("NGO not found")
	}
	if !donor.KYCVerified {
		return nil, fmt.Errorf("donor KYC not verified")
	}
	if !ngo.KYCData.Verified {
		return nil, fmt.Errorf("NGO KYC not verified")
	}

	// Check donation limit
	limitCheck := donor.CheckDonationLimit(amount)
	if !limitCheck.CanDonate {
		return nil, fmt.Errorf("donation exceeds annual limit. Remaining: ₹%.2f", limitCheck.RemainingLimit)
	}

	// Calculate platform fee
	platformFee := amount * p.SystemStats.PlatformFee
	netAmount := amount - platformFee

	donation := transactions.NewDonationTransaction(donorID, ngoID, netAmount, paymentMethod, donor.KYCData.DocumentHash)

	result, err := ngo.ProcessDonation(donation)
	if err != nil {
		return nil, err
	}

	donor.AddDonation(donation)

	// Update system stats
	p.SystemStats.TotalTransactions++
	p.SystemStats.TotalDonations += netAmount

	// Anchor to Polygon if available
	if p.PolygonIntegration != nil {
		additionalData := map[string]interface{}{
			"amount":       netAmount,
			"platform_fee": platformFee,
		}

		anchor, err := p.PolygonIntegration.AnchorBlockHash(result.BlockHash, ngoID, "donation", additionalData)
		if err == nil {
			result.EBill = map[string]interface{}{
				"polygon_anchor": anchor,
				"original_ebill": donation.EBill,
			}
		}
	}

	return map[string]interface{}{
		"success":        result.Success,
		"block_hash":     result.BlockHash,
		"transaction_id": result.TransactionID,
		"block_index":    result.BlockIndex,
		"e_bill":         result.EBill,
		"platform_fee":   platformFee,
		"net_amount":     netAmount,
		"gross_amount":   amount,
	}, nil
}

// ProcessExpenditure processes an expenditure transaction
func (p *NGOTransparencyPlatform) ProcessExpenditure(ngoID string, expenditureData map[string]interface{}, auditorID string) (map[string]interface{}, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ngo, ngoExists := p.NGOs[ngoID]
	auditor, auditorExists := p.Auditors[auditorID]

	if !ngoExists {
		return nil, fmt.Errorf("NGO not found")
	}
	if !auditorExists {
		return nil, fmt.Errorf("auditor not found")
	}
	if !auditor.Verified {
		return nil, fmt.Errorf("auditor not verified")
	}

	// Extract expenditure data
	amount, ok := expenditureData["amount"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid amount")
	}

	category, _ := expenditureData["category"].(string)
	description, _ := expenditureData["description"].(string)

	// Create invoice details (simplified)
	invoiceDetails := transactions.InvoiceDetails{
		InvoiceNumber: fmt.Sprintf("INV-%d", time.Now().Unix()),
		GSTIN:         "27ABCDE1234F1Z5", // Example GSTIN
		VendorName:    "Test Vendor",
		VendorGSTIN:   "29ABCDE1234F1Z3", // Example Vendor GSTIN
		InvoiceDate:   time.Now(),
		Documents:     []string{"invoice.pdf"},
	}

	// Create expenditure transaction
	expenditure := transactions.NewExpenditureTransaction(ngoID, amount, category, description, invoiceDetails, []transactions.Attachment{
		{
			Filename:   "invoice.pdf",
			Hash:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			Type:       "application/pdf",
			UploadedAt: time.Now(),
		},
	})

	// Auditor performs audit
	auditResult := auditor.AuditExpenditure(expenditure, "")

	// Auto-approve based on audit result
	shouldApprove := auditResult.ComplianceScore >= 70.0 // Require at least 70% compliance
	expenditure.ValidateByAuditor(auditorID, shouldApprove, auditResult.Recommendation, &auditResult.ComplianceScore)

	if !shouldApprove {
		return nil, fmt.Errorf("expenditure rejected by auditor (score: %.2f%%): %s",
			auditResult.ComplianceScore, auditResult.Recommendation)
	}

	result, err := ngo.ProcessExpenditure(expenditure)
	if err != nil {
		return nil, err
	}

	// Update system stats
	p.SystemStats.TotalTransactions++
	p.SystemStats.TotalExpenditures += amount

	// Anchor to Polygon if available
	if p.PolygonIntegration != nil {
		additionalData := map[string]interface{}{
			"amount":   amount,
			"category": category,
		}

		anchor, err := p.PolygonIntegration.AnchorBlockHash(result.BlockHash, ngoID, "expenditure", additionalData)
		if err == nil {
			result.EBill = map[string]interface{}{
				"polygon_anchor": anchor,
			}
		}
	}

	return map[string]interface{}{
		"success":        result.Success,
		"block_hash":     result.BlockHash,
		"transaction_id": result.TransactionID,
		"block_index":    result.BlockIndex,
		"audit_result":   auditResult,
	}, nil
}

// CalculateAllNGORatings calculates ratings for all NGOs
func (p *NGOTransparencyPlatform) CalculateAllNGORatings(periodDays int) []map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var ratings []map[string]interface{}

	for ngoID, ngo := range p.NGOs {
		rating := ngo.CalculateRating(periodDays)
		ratingInfo := map[string]interface{}{
			"ngo_id":                ngoID,
			"name":                  ngo.Name,
			"category":              ngo.Category,
			"kyc_verified":          ngo.KYCData.Verified,
			"rating":                rating.Rating,
			"transparency_score":    rating.TransparencyScore,
			"utilization_rate":      rating.UtilizationRate,
			"gap_percentage":        rating.GapPercentage,
			"total_donations":       rating.TotalDonations,
			"total_expenditures":    rating.TotalExpenditures,
			"documentation_quality": rating.DocumentationQuality,
		}
		ratings = append(ratings, ratingInfo)
	}

	// Sort by rating (descending)
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i]["rating"].(float64) > ratings[j]["rating"].(float64)
	})

	return ratings
}

// GetPlatformStats returns comprehensive platform statistics
func (p *NGOTransparencyPlatform) GetPlatformStats() PlatformStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	now := time.Now()
	platformAge := now.Sub(p.SystemStats.CreatedAt)
	daysActive := int(platformAge.Hours() / 24)

	verifiedNGOs := 0
	verifiedDonors := 0
	verifiedAuditors := 0
	totalRating := 0.0
	categories := make(map[string]bool)

	for _, ngo := range p.NGOs {
		if ngo.KYCData.Verified {
			verifiedNGOs++
		}
		totalRating += ngo.Rating
		categories[ngo.Category] = true
	}

	for _, donor := range p.Donors {
		if donor.KYCVerified {
			verifiedDonors++
		}
	}

	for _, auditor := range p.Auditors {
		if auditor.Verified {
			verifiedAuditors++
		}
	}

	averageRating := 0.0
	if len(p.NGOs) > 0 {
		averageRating = totalRating / float64(len(p.NGOs))
	}

	categoryList := make([]string, 0, len(categories))
	for category := range categories {
		categoryList = append(categoryList, category)
	}

	return PlatformStats{
		TotalNGOs:            len(p.NGOs),
		TotalDonors:          len(p.Donors),
		TotalAuditors:        len(p.Auditors),
		TotalTransactions:    p.SystemStats.TotalTransactions,
		TotalDonations:       p.SystemStats.TotalDonations,
		TotalExpenditures:    p.SystemStats.TotalExpenditures,
		PlatformFeeCollected: p.SystemStats.TotalDonations * p.SystemStats.PlatformFee,
		VerifiedNGOs:         verifiedNGOs,
		VerifiedDonors:       verifiedDonors,
		VerifiedAuditors:     verifiedAuditors,
		KYCAuthorities:       len(p.KYCAuthorities),
		DaysActive:           daysActive,
		AverageNGORating:     averageRating,
		Categories:           categoryList,
	}
}

// GetNGODashboard returns NGO dashboard information
func (p *NGOTransparencyPlatform) GetNGODashboard(ngoID string) (map[string]interface{}, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	ngo, exists := p.NGOs[ngoID]
	if !exists {
		return nil, fmt.Errorf("NGO not found")
	}

	stats := ngo.GetBlockchainStats()
	financialSummary := ngo.GetFinancialSummary(12)
	ratingDetails := ngo.CalculateRating(30)

	return map[string]interface{}{
		"stats":             stats,
		"financial_summary": financialSummary,
		"rating_details":    ratingDetails,
		"multisig_status": map[string]interface{}{
			"signers":              ngo.MultiSigWallet.GetSigners(),
			"required_signatures":  2,
			"pending_transactions": ngo.MultiSigWallet.GetPendingTransactionCount(),
		},
	}, nil
}

// GetDonorDashboard returns donor dashboard information
func (p *NGOTransparencyPlatform) GetDonorDashboard(donorID string) (map[string]interface{}, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	donor, exists := p.Donors[donorID]
	if !exists {
		return nil, fmt.Errorf("donor not found")
	}

	stats := donor.GetDonorStats()
	recentDonations := donor.GetDonationHistory(10)
	currentYearTaxBenefits := donor.GetAnnualTaxBenefits(0) // Current year

	preferredNGOs := make([]map[string]interface{}, 0)
	for _, ngoID := range donor.PreferredNGOs {
		if ngo, exists := p.NGOs[ngoID]; exists {
			preferredNGOs = append(preferredNGOs, map[string]interface{}{
				"ngo_id":   ngoID,
				"name":     ngo.Name,
				"rating":   ngo.Rating,
				"category": ngo.Category,
			})
		}
	}

	return map[string]interface{}{
		"stats":                     stats,
		"recent_donations":          recentDonations,
		"current_year_tax_benefits": currentYearTaxBenefits,
		"preferred_ngos":            preferredNGOs,
	}, nil
}

// GetAuditorDashboard returns auditor dashboard information
func (p *NGOTransparencyPlatform) GetAuditorDashboard(auditorID string) (map[string]interface{}, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	auditor, exists := p.Auditors[auditorID]
	if !exists {
		return nil, fmt.Errorf("auditor not found")
	}

	stats := auditor.GetAuditorStats()
	recentAudits := auditor.GetRecentAudits(5)

	return map[string]interface{}{
		"stats":         stats,
		"recent_audits": recentAudits,
	}, nil
}
