package transactions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"
)

// InvoiceDetails represents invoice information
type InvoiceDetails struct {
	InvoiceNumber     string    `json:"invoice_number"`
	GSTIN             string    `json:"gstin"`
	VendorName        string    `json:"vendor_name"`
	VendorGSTIN       string    `json:"vendor_gstin"`
	InvoiceDate       time.Time `json:"invoice_date"`
	Documents         []string  `json:"documents"`
	BankTransactionID string    `json:"bank_transaction_id"`
	ChequeNumber      string    `json:"cheque_number"`
}

// Attachment represents a file attachment
type Attachment struct {
	Filename   string    `json:"filename"`
	Hash       string    `json:"hash"`
	Type       string    `json:"type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// AuditorValidation represents auditor validation information
type AuditorValidation struct {
	AuditorID  string    `json:"auditor_id"`
	IsValid    bool      `json:"is_valid"`
	Remarks    string    `json:"remarks"`
	AuditScore float64   `json:"audit_score"`
	Timestamp  time.Time `json:"timestamp"`
	Signature  string    `json:"signature"`
}

// ExpenditureTransaction represents an expenditure transaction
type ExpenditureTransaction struct {
	TransactionID     string             `json:"transaction_id"`
	NGOID             string             `json:"ngo_id"`
	Amount            float64            `json:"amount"`
	Category          string             `json:"category"`
	Description       string             `json:"description"`
	Timestamp         time.Time          `json:"timestamp"`
	InvoiceDetails    InvoiceDetails     `json:"invoice_details"`
	Attachments       []Attachment       `json:"attachments"`
	Status            string             `json:"status"`
	AuditorValidation *AuditorValidation `json:"auditor_validation,omitempty"`
	ComplianceScore   float64            `json:"compliance_score"`
}

// NewExpenditureTransaction creates a new expenditure transaction
func NewExpenditureTransaction(ngoID string, amount float64, category, description string, invoiceDetails InvoiceDetails, attachments []Attachment) *ExpenditureTransaction {
	// Generate transaction ID
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	transactionID := hex.EncodeToString(randomBytes)

	if attachments == nil {
		attachments = make([]Attachment, 0)
	}

	transaction := &ExpenditureTransaction{
		TransactionID:  transactionID,
		NGOID:          ngoID,
		Amount:         amount,
		Category:       category,
		Description:    description,
		Timestamp:      time.Now(),
		InvoiceDetails: invoiceDetails,
		Attachments:    attachments,
		Status:         "pending_validation",
	}

	// Calculate compliance score
	transaction.ComplianceScore = transaction.calculateComplianceScore()

	return transaction
}

// calculateComplianceScore calculates the compliance score for the expenditure
func (et *ExpenditureTransaction) calculateComplianceScore() float64 {
	var score float64 = 0
	maxScore := 100.0

	// Invoice number present (20%)
	if et.InvoiceDetails.InvoiceNumber != "" {
		score += 20
	}

	// GSTIN validation (20%)
	if et.VerifyGSTIN(et.InvoiceDetails.GSTIN) {
		score += 20
	}

	// Vendor details complete (15%)
	if et.InvoiceDetails.VendorName != "" && et.InvoiceDetails.VendorGSTIN != "" {
		score += 15
	}

	// Payment proof (15%)
	if et.InvoiceDetails.BankTransactionID != "" || et.InvoiceDetails.ChequeNumber != "" {
		score += 15
	}

	// Supporting documents (15%)
	if len(et.InvoiceDetails.Documents) > 0 {
		score += 15
	}

	// Attachments (10%)
	if len(et.Attachments) > 0 {
		score += 10
	}

	// Recent invoice (within 90 days) (5%)
	daysDiff := time.Since(et.InvoiceDetails.InvoiceDate).Hours() / 24
	if daysDiff <= 90 {
		score += 5
	}

	return min(score, maxScore)
}

// ValidateByAuditor validates the expenditure by an auditor
func (et *ExpenditureTransaction) ValidateByAuditor(auditorID string, isValid bool, remarks string, score *float64) {
	if score == nil {
		defaultScore := et.calculateComplianceScore()
		score = &defaultScore
	}

	// Create validation signature
	hashData := fmt.Sprintf("%s|%s|%t|%s|%f|%d",
		et.TransactionID,
		auditorID,
		isValid,
		remarks,
		*score,
		time.Now().Unix(),
	)
	hash := sha256.Sum256([]byte(hashData))
	signature := hex.EncodeToString(hash[:])

	et.AuditorValidation = &AuditorValidation{
		AuditorID:  auditorID,
		IsValid:    isValid,
		Remarks:    remarks,
		AuditScore: *score,
		Timestamp:  time.Now(),
		Signature:  signature,
	}

	if isValid {
		et.Status = "validated"
		et.ComplianceScore = *score
	} else {
		et.Status = "rejected"
		et.ComplianceScore = 0
	}
}

// VerifyGSTIN validates GSTIN format
func (et *ExpenditureTransaction) VerifyGSTIN(gstin string) bool {
	if gstin == "" {
		return false
	}

	// GSTIN format: 2-digit state code + 10-character PAN + 1-character entity code + 1-character check digit + 1-character default 'Z' + 1-character check digit
	gstinRegex := regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`)
	return gstinRegex.MatchString(gstin)
}

// VerifyInvoiceUniqueness checks if the invoice number is unique among existing invoices
func (et *ExpenditureTransaction) VerifyInvoiceUniqueness(existingInvoices []string) bool {
	for _, existing := range existingInvoices {
		if existing == et.InvoiceDetails.InvoiceNumber {
			return false
		}
	}
	return true
}

// AddAttachment adds a new attachment to the transaction
func (et *ExpenditureTransaction) AddAttachment(filename, hash, attachmentType string) {
	attachment := Attachment{
		Filename:   filename,
		Hash:       hash,
		Type:       attachmentType,
		UploadedAt: time.Now(),
	}

	et.Attachments = append(et.Attachments, attachment)

	// Recalculate compliance score
	et.ComplianceScore = et.calculateComplianceScore()
}

// IsValidated checks if the transaction is validated
func (et *ExpenditureTransaction) IsValidated() bool {
	return et.Status == "validated"
}

// IsRejected checks if the transaction is rejected
func (et *ExpenditureTransaction) IsRejected() bool {
	return et.Status == "rejected"
}

// IsPendingValidation checks if the transaction is pending validation
func (et *ExpenditureTransaction) IsPendingValidation() bool {
	return et.Status == "pending_validation"
}

// GetTransactionSummary returns a summary of the expenditure transaction
func (et *ExpenditureTransaction) GetTransactionSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"transaction_id":   et.TransactionID,
		"ngo_id":           et.NGOID,
		"amount":           et.Amount,
		"currency":         "INR",
		"category":         et.Category,
		"description":      et.Description,
		"status":           et.Status,
		"timestamp":        et.Timestamp,
		"compliance_score": et.ComplianceScore,
		"attachments":      len(et.Attachments),
		"gstin_valid":      et.VerifyGSTIN(et.InvoiceDetails.GSTIN),
	}

	if et.AuditorValidation != nil {
		summary["auditor_validation"] = map[string]interface{}{
			"auditor_id":   et.AuditorValidation.AuditorID,
			"is_valid":     et.AuditorValidation.IsValid,
			"remarks":      et.AuditorValidation.Remarks,
			"audit_score":  et.AuditorValidation.AuditScore,
			"validated_at": et.AuditorValidation.Timestamp,
		}
	}

	return summary
}

// GetInvoiceInfo returns detailed invoice information
func (et *ExpenditureTransaction) GetInvoiceInfo() map[string]interface{} {
	return map[string]interface{}{
		"invoice_number":      et.InvoiceDetails.InvoiceNumber,
		"gstin":               et.InvoiceDetails.GSTIN,
		"vendor_name":         et.InvoiceDetails.VendorName,
		"vendor_gstin":        et.InvoiceDetails.VendorGSTIN,
		"invoice_date":        et.InvoiceDetails.InvoiceDate,
		"documents":           et.InvoiceDetails.Documents,
		"bank_transaction_id": et.InvoiceDetails.BankTransactionID,
		"cheque_number":       et.InvoiceDetails.ChequeNumber,
		"gstin_valid":         et.VerifyGSTIN(et.InvoiceDetails.GSTIN),
		"vendor_gstin_valid":  et.VerifyGSTIN(et.InvoiceDetails.VendorGSTIN),
	}
}

// GetComplianceBreakdown returns detailed compliance score breakdown
func (et *ExpenditureTransaction) GetComplianceBreakdown() map[string]interface{} {
	breakdown := map[string]interface{}{
		"total_score": et.ComplianceScore,
		"max_score":   100.0,
	}

	components := make(map[string]interface{})

	// Invoice number check
	components["invoice_number"] = map[string]interface{}{
		"score":     getScore(et.InvoiceDetails.InvoiceNumber != "", 20),
		"max_score": 20,
		"status":    et.InvoiceDetails.InvoiceNumber != "",
	}

	// GSTIN validation
	components["gstin_validation"] = map[string]interface{}{
		"score":     getScore(et.VerifyGSTIN(et.InvoiceDetails.GSTIN), 20),
		"max_score": 20,
		"status":    et.VerifyGSTIN(et.InvoiceDetails.GSTIN),
	}

	// Vendor details
	vendorComplete := et.InvoiceDetails.VendorName != "" && et.InvoiceDetails.VendorGSTIN != ""
	components["vendor_details"] = map[string]interface{}{
		"score":     getScore(vendorComplete, 15),
		"max_score": 15,
		"status":    vendorComplete,
	}

	// Payment proof
	hasPaymentProof := et.InvoiceDetails.BankTransactionID != "" || et.InvoiceDetails.ChequeNumber != ""
	components["payment_proof"] = map[string]interface{}{
		"score":     getScore(hasPaymentProof, 15),
		"max_score": 15,
		"status":    hasPaymentProof,
	}

	// Supporting documents
	hasDocuments := len(et.InvoiceDetails.Documents) > 0
	components["supporting_documents"] = map[string]interface{}{
		"score":     getScore(hasDocuments, 15),
		"max_score": 15,
		"status":    hasDocuments,
		"count":     len(et.InvoiceDetails.Documents),
	}

	// Attachments
	hasAttachments := len(et.Attachments) > 0
	components["attachments"] = map[string]interface{}{
		"score":     getScore(hasAttachments, 10),
		"max_score": 10,
		"status":    hasAttachments,
		"count":     len(et.Attachments),
	}

	// Invoice recency
	daysDiff := time.Since(et.InvoiceDetails.InvoiceDate).Hours() / 24
	isRecent := daysDiff <= 90
	components["invoice_recency"] = map[string]interface{}{
		"score":     getScore(isRecent, 5),
		"max_score": 5,
		"status":    isRecent,
		"days_old":  int(daysDiff),
		"threshold": 90,
	}

	breakdown["components"] = components
	return breakdown
}

// getScore returns the score based on condition and max value
func getScore(condition bool, maxScore float64) float64 {
	if condition {
		return maxScore
	}
	return 0
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// GetValidationRecommendation returns a recommendation based on compliance score
func (et *ExpenditureTransaction) GetValidationRecommendation() string {
	score := et.ComplianceScore

	switch {
	case score >= 90:
		return "Approve - Excellent compliance"
	case score >= 70:
		return "Approve with minor observations"
	case score >= 50:
		return "Conditional approval - requires additional documentation"
	default:
		return "Reject - insufficient compliance"
	}
}

// GetComplianceIssues returns a list of compliance issues
func (et *ExpenditureTransaction) GetComplianceIssues() []string {
	var issues []string

	if et.InvoiceDetails.InvoiceNumber == "" {
		issues = append(issues, "Missing invoice number")
	}

	if !et.VerifyGSTIN(et.InvoiceDetails.GSTIN) {
		issues = append(issues, "Invalid GSTIN format")
	}

	if et.InvoiceDetails.VendorName == "" {
		issues = append(issues, "Missing vendor name")
	}

	if et.InvoiceDetails.VendorGSTIN == "" {
		issues = append(issues, "Missing vendor GSTIN")
	}

	if et.InvoiceDetails.BankTransactionID == "" && et.InvoiceDetails.ChequeNumber == "" {
		issues = append(issues, "Missing payment proof (bank transaction ID or cheque number)")
	}

	if len(et.InvoiceDetails.Documents) == 0 {
		issues = append(issues, "No supporting documents provided")
	}

	if len(et.Attachments) == 0 {
		issues = append(issues, "No file attachments provided")
	}

	daysDiff := time.Since(et.InvoiceDetails.InvoiceDate).Hours() / 24
	if daysDiff > 90 {
		issues = append(issues, fmt.Sprintf("Invoice is too old (%d days, threshold: 90 days)", int(daysDiff)))
	}

	return issues
}
