package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"ngo-transparency-platform/pkg/crypto"
	"ngo-transparency-platform/pkg/transactions"
	"time"
)

// DonorKYCData represents KYC information for donors
type DonorKYCData struct {
	DocumentHash          string    `json:"document_hash"`
	VerificationDate      *time.Time `json:"verification_date"`
	VerificationAuthority string    `json:"verification_authority"`
	DocumentsSubmitted    []string  `json:"documents_submitted"`
	VerificationLevel     string    `json:"verification_level"`
}

// DonationRecord represents a single donation record
type DonationRecord struct {
	TransactionID string                            `json:"transaction_id"`
	NGOID         string                            `json:"ngo_id"`
	Amount        float64                           `json:"amount"`
	Timestamp     time.Time                         `json:"timestamp"`
	EBill         *transactions.EBill               `json:"e_bill"`
	ZKProof       *crypto.ZKProof            `json:"zk_proof"`
	TaxBenefit    transactions.TaxBenefit          `json:"tax_benefit"`
}

// TaxBenefitSummary represents annual tax benefit summary
type TaxBenefitSummary struct {
	Year              int                   `json:"year"`
	TotalDonated      float64              `json:"total_donated"`
	TotalDeductible   float64              `json:"total_deductible"`
	EstimatedTaxSaving float64             `json:"estimated_tax_saving"`
	Donations         []DonationRecord     `json:"donations"`
}

// DonationLimit represents donation limit checking result
type DonationLimit struct {
	CanDonate         bool    `json:"can_donate"`
	CurrentYearTotal  float64 `json:"current_year_total"`
	Limit             float64 `json:"limit"`
	RemainingLimit    float64 `json:"remaining_limit"`
}

// DonorStats represents donor statistics
type DonorStats struct {
	DonorID              string    `json:"donor_id"`
	KYCVerified          bool      `json:"kyc_verified"`
	VerificationLevel    string    `json:"verification_level"`
	TotalDonated         float64   `json:"total_donated"`
	DonationCount        int       `json:"donation_count"`
	CurrentYearDonations float64   `json:"current_year_donations"`
	CurrentYearCount     int       `json:"current_year_count"`
	PreferredNGOsCount   int       `json:"preferred_ngos_count"`
	AverageDonation      float64   `json:"average_donation"`
	MemberSince          time.Time `json:"member_since"`
	AnnualLimit          float64   `json:"annual_limit"`
}

// Donor represents a donor in the system
type Donor struct {
	DonorID         string            `json:"donor_id"`
	KYCVerified     bool              `json:"kyc_verified"`
	KYCData         DonorKYCData      `json:"kyc_data"`
	DonationHistory []DonationRecord  `json:"donation_history"`
	TotalDonated    float64           `json:"total_donated"`
	PreferredNGOs   []string          `json:"preferred_ngos"`
	TaxBenefits     []TaxBenefitSummary `json:"tax_benefits"`
	CreatedAt       time.Time         `json:"created_at"`
	AnnualDonationLimit float64       `json:"annual_donation_limit"`
}

// NewDonor creates a new donor instance
func NewDonor(donorID string, kycData map[string]interface{}) *Donor {
	// Create KYC data hash
	kycDocHash := ""
	documentsSubmitted := []string{}
	annualLimit := 1000000.0 // 10 lakh default

	if docs, ok := kycData["documents"]; ok {
		if docsList, ok := docs.([]string); ok {
			documentsSubmitted = docsList
			kycBytes := fmt.Sprintf("%v", docs)
			hash := sha256.Sum256([]byte(kycBytes))
			kycDocHash = hex.EncodeToString(hash[:])
		}
	}

	if limit, ok := kycData["annual_limit"].(float64); ok {
		annualLimit = limit
	}

	return &Donor{
		DonorID:     donorID,
		KYCVerified: false,
		KYCData: DonorKYCData{
			DocumentHash:       kycDocHash,
			DocumentsSubmitted: documentsSubmitted,
		},
		DonationHistory:     make([]DonationRecord, 0),
		TotalDonated:        0,
		PreferredNGOs:       make([]string, 0),
		TaxBenefits:         make([]TaxBenefitSummary, 0),
		CreatedAt:           time.Now(),
		AnnualDonationLimit: annualLimit,
	}
}

// VerifyKYC verifies the donor's KYC
func (d *Donor) VerifyKYC(authorityID, verificationLevel string) bool {
	if verificationLevel == "" {
		verificationLevel = "basic"
	}

	d.KYCVerified = true
	now := time.Now()
	d.KYCData.VerificationDate = &now
	d.KYCData.VerificationAuthority = authorityID
	d.KYCData.VerificationLevel = verificationLevel

	// Update annual limit based on verification level
	if verificationLevel == "premium" {
		d.AnnualDonationLimit = 5000000 // 50 lakh for premium KYC
	}

	return d.KYCVerified
}

// AddDonation adds a donation record to the donor's history
func (d *Donor) AddDonation(donation *transactions.DonationTransaction) {
	donationRecord := DonationRecord{
		TransactionID: donation.TransactionID,
		NGOID:         donation.NGOID,
		Amount:        donation.Amount,
		Timestamp:     donation.Timestamp,
		EBill:         donation.EBill,
		ZKProof:       donation.ZKProof,
		TaxBenefit:    donation.EBill.TaxBenefit,
	}

	d.DonationHistory = append(d.DonationHistory, donationRecord)
	d.TotalDonated += donation.Amount

	// Add to tax benefits for the current year
	d.updateTaxBenefits(donationRecord)
}

// updateTaxBenefits updates the tax benefits summary
func (d *Donor) updateTaxBenefits(donation DonationRecord) {
	year := donation.Timestamp.Year()

	// Find or create tax benefit entry for this year
	var yearBenefit *TaxBenefitSummary
	for i := range d.TaxBenefits {
		if d.TaxBenefits[i].Year == year {
			yearBenefit = &d.TaxBenefits[i]
			break
		}
	}

	if yearBenefit == nil {
		// Create new year entry
		newBenefit := TaxBenefitSummary{
			Year:      year,
			Donations: make([]DonationRecord, 0),
		}
		d.TaxBenefits = append(d.TaxBenefits, newBenefit)
		yearBenefit = &d.TaxBenefits[len(d.TaxBenefits)-1]
	}

	// Add donation to this year's summary
	yearBenefit.Donations = append(yearBenefit.Donations, donation)
	yearBenefit.TotalDonated += donation.Amount
	yearBenefit.TotalDeductible += donation.TaxBenefit.DeductibleAmount
	yearBenefit.EstimatedTaxSaving += donation.TaxBenefit.TaxSaving
}

// AddPreferredNGO adds an NGO to the preferred list
func (d *Donor) AddPreferredNGO(ngoID string) {
	for _, existingNGO := range d.PreferredNGOs {
		if existingNGO == ngoID {
			return // Already in preferred list
		}
	}
	d.PreferredNGOs = append(d.PreferredNGOs, ngoID)
}

// RemovePreferredNGO removes an NGO from the preferred list
func (d *Donor) RemovePreferredNGO(ngoID string) {
	var newPreferred []string
	for _, existingNGO := range d.PreferredNGOs {
		if existingNGO != ngoID {
			newPreferred = append(newPreferred, existingNGO)
		}
	}
	d.PreferredNGOs = newPreferred
}

// GetDonationHistory returns donation history (with optional limit)
func (d *Donor) GetDonationHistory(limit int) []DonationRecord {
	if limit <= 0 || limit >= len(d.DonationHistory) {
		return d.DonationHistory
	}

	// Return the most recent donations
	start := len(d.DonationHistory) - limit
	return d.DonationHistory[start:]
}

// GetAnnualTaxBenefits returns tax benefits for a specific year
func (d *Donor) GetAnnualTaxBenefits(year int) *TaxBenefitSummary {
	if year == 0 {
		year = time.Now().Year()
	}

	for _, benefit := range d.TaxBenefits {
		if benefit.Year == year {
			return &benefit
		}
	}

	// Return empty summary if no donations for this year
	return &TaxBenefitSummary{
		Year:      year,
		Donations: make([]DonationRecord, 0),
	}
}

// CheckDonationLimit checks if the donor can make a donation of the specified amount
func (d *Donor) CheckDonationLimit(amount float64) DonationLimit {
	currentYear := time.Now().Year()
	yearlyTotal := 0.0

	for _, donation := range d.DonationHistory {
		if donation.Timestamp.Year() == currentYear {
			yearlyTotal += donation.Amount
		}
	}

	canDonate := (yearlyTotal + amount) <= d.AnnualDonationLimit
	remainingLimit := d.AnnualDonationLimit - yearlyTotal

	return DonationLimit{
		CanDonate:        canDonate,
		CurrentYearTotal: yearlyTotal,
		Limit:            d.AnnualDonationLimit,
		RemainingLimit:   remainingLimit,
	}
}

// GetDonorStats returns comprehensive donor statistics
func (d *Donor) GetDonorStats() DonorStats {
	currentYear := time.Now().Year()
	currentYearTotal := 0.0
	currentYearCount := 0

	for _, donation := range d.DonationHistory {
		if donation.Timestamp.Year() == currentYear {
			currentYearTotal += donation.Amount
			currentYearCount++
		}
	}

	averageDonation := 0.0
	if len(d.DonationHistory) > 0 {
		averageDonation = d.TotalDonated / float64(len(d.DonationHistory))
	}

	verificationLevel := "none"
	if d.KYCVerified {
		verificationLevel = d.KYCData.VerificationLevel
	}

	return DonorStats{
		DonorID:              d.DonorID,
		KYCVerified:          d.KYCVerified,
		VerificationLevel:    verificationLevel,
		TotalDonated:         d.TotalDonated,
		DonationCount:        len(d.DonationHistory),
		CurrentYearDonations: currentYearTotal,
		CurrentYearCount:     currentYearCount,
		PreferredNGOsCount:   len(d.PreferredNGOs),
		AverageDonation:      averageDonation,
		MemberSince:          d.CreatedAt,
		AnnualLimit:          d.AnnualDonationLimit,
	}
}

// GetDonationsByNGO returns donations made to a specific NGO
func (d *Donor) GetDonationsByNGO(ngoID string) []DonationRecord {
	var donations []DonationRecord
	for _, donation := range d.DonationHistory {
		if donation.NGOID == ngoID {
			donations = append(donations, donation)
		}
	}
	return donations
}

// GetDonationsByDateRange returns donations within a date range
func (d *Donor) GetDonationsByDateRange(startDate, endDate time.Time) []DonationRecord {
	var donations []DonationRecord
	for _, donation := range d.DonationHistory {
		if (donation.Timestamp.After(startDate) || donation.Timestamp.Equal(startDate)) &&
			(donation.Timestamp.Before(endDate) || donation.Timestamp.Equal(endDate)) {
			donations = append(donations, donation)
		}
	}
	return donations
}

// GetMonthlyDonationSummary returns monthly donation summary for the current year
func (d *Donor) GetMonthlyDonationSummary() map[string]float64 {
	currentYear := time.Now().Year()
	monthlySummary := make(map[string]float64)

	// Initialize all months
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	for _, month := range months {
		monthlySummary[month] = 0.0
	}

	// Calculate monthly totals
	for _, donation := range d.DonationHistory {
		if donation.Timestamp.Year() == currentYear {
			monthName := donation.Timestamp.Month().String()
			monthlySummary[monthName] += donation.Amount
		}
	}

	return monthlySummary
}
