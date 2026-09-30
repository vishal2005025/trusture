package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"ngo-transparency-platform/pkg/transactions"
	"strings"
	"time"
)

// AuditResult represents the result of an audit
type AuditResult struct {
	AuditID         string    `json:"audit_id"`
	ExpenditureID   string    `json:"expenditure_id"`
	AuditorID       string    `json:"auditor_id"`
	Timestamp       time.Time `json:"timestamp"`
	ComplianceScore float64   `json:"compliance_score"`
	Findings        []string  `json:"findings"`
	Recommendation  string    `json:"recommendation"`
	AuditNotes      string    `json:"audit_notes"`
	Signature       string    `json:"signature"`
}

// AuditorStats represents auditor statistics
type AuditorStats struct {
	AuditorID              string    `json:"auditor_id"`
	Name                   string    `json:"name"`
	Verified               bool      `json:"verified"`
	Specializations        []string  `json:"specializations"`
	Rating                 float64   `json:"rating"`
	TotalAudits            int       `json:"total_audits"`
	ApprovedAudits         int       `json:"approved_audits"`
	ApprovalRate           string    `json:"approval_rate"`
	AverageComplianceScore float64   `json:"average_compliance_score"`
	MemberSince            time.Time `json:"member_since"`
}

// Auditor represents an auditor in the system
type Auditor struct {
	AuditorID             string        `json:"auditor_id"`
	Name                  string        `json:"name"`
	Credentials           interface{}   `json:"credentials"`
	Specializations       []string      `json:"specializations"` // ['financial', 'compliance', 'technical']
	Verified              bool          `json:"verified"`
	AuditHistory          []AuditResult `json:"audit_history"`
	Rating                float64       `json:"rating"`
	CreatedAt             time.Time     `json:"created_at"`
	PublicKey             string        `json:"public_key"`
	VerificationAuthority string        `json:"verification_authority,omitempty"`
	VerificationDate      *time.Time    `json:"verification_date,omitempty"`
}

// NewAuditor creates a new auditor instance
func NewAuditor(auditorID, name string, credentials interface{}, specializations []string) *Auditor {
	if specializations == nil {
		specializations = make([]string, 0)
	}

	// Generate public key
	publicKeyData := auditorID + name
	hash := sha256.Sum256([]byte(publicKeyData))
	publicKey := hex.EncodeToString(hash[:])

	return &Auditor{
		AuditorID:       auditorID,
		Name:            name,
		Credentials:     credentials,
		Specializations: specializations,
		Verified:        false,
		AuditHistory:    make([]AuditResult, 0),
		Rating:          5.0,
		CreatedAt:       time.Now(),
		PublicKey:       publicKey,
	}
}

// VerifyCredentials verifies the auditor's credentials
func (a *Auditor) VerifyCredentials(verificationAuthority string) bool {
	a.Verified = true
	a.VerificationAuthority = verificationAuthority
	now := time.Now()
	a.VerificationDate = &now
	return a.Verified
}

// AuditExpenditure performs an audit on an expenditure transaction
func (a *Auditor) AuditExpenditure(expenditure *transactions.ExpenditureTransaction, auditNotes string) *AuditResult {
	// Generate audit ID - handle variable length auditorIDs
	auditorIDPart := a.AuditorID
	if len(a.AuditorID) > 8 {
		auditorIDPart = a.AuditorID[:8]
	}
	auditID := fmt.Sprintf("AUD-%d-%s", time.Now().Unix(), auditorIDPart)

	findings := a.generateFindings(expenditure)
	recommendation := a.generateRecommendation(expenditure)

	// Generate signature
	signatureData := fmt.Sprintf("%s%s%s", auditID, expenditure.TransactionID, a.AuditorID)
	hash := sha256.Sum256([]byte(signatureData))
	signature := hex.EncodeToString(hash[:])

	auditResult := &AuditResult{
		AuditID:         auditID,
		ExpenditureID:   expenditure.TransactionID,
		AuditorID:       a.AuditorID,
		Timestamp:       time.Now(),
		ComplianceScore: expenditure.ComplianceScore,
		Findings:        findings,
		Recommendation:  recommendation,
		AuditNotes:      auditNotes,
		Signature:       signature,
	}

	a.AuditHistory = append(a.AuditHistory, *auditResult)
	return auditResult
}

// generateFindings generates audit findings based on expenditure analysis
func (a *Auditor) generateFindings(expenditure *transactions.ExpenditureTransaction) []string {
	var findings []string

	if !expenditure.VerifyGSTIN(expenditure.InvoiceDetails.GSTIN) {
		findings = append(findings, "Invalid GSTIN format")
	}

	if expenditure.InvoiceDetails.BankTransactionID == "" && expenditure.InvoiceDetails.ChequeNumber == "" {
		findings = append(findings, "Missing payment proof")
	}

	if len(expenditure.InvoiceDetails.Documents) == 0 {
		findings = append(findings, "No supporting documents provided")
	}

	if expenditure.ComplianceScore < 80 {
		findings = append(findings, fmt.Sprintf("Low compliance score: %.1f%%", expenditure.ComplianceScore))
	}

	// Check for vendor GSTIN if available
	if expenditure.InvoiceDetails.VendorGSTIN != "" && !expenditure.VerifyGSTIN(expenditure.InvoiceDetails.VendorGSTIN) {
		findings = append(findings, "Invalid vendor GSTIN format")
	}

	// Check invoice age
	daysDiff := time.Since(expenditure.InvoiceDetails.InvoiceDate).Hours() / 24
	if daysDiff > 90 {
		findings = append(findings, fmt.Sprintf("Invoice is older than 90 days (%.0f days)", daysDiff))
	}

	// Check for missing critical information
	if expenditure.InvoiceDetails.InvoiceNumber == "" {
		findings = append(findings, "Missing invoice number")
	}

	if expenditure.InvoiceDetails.VendorName == "" {
		findings = append(findings, "Missing vendor name")
	}

	return findings
}

// generateRecommendation generates a recommendation based on compliance score and findings
func (a *Auditor) generateRecommendation(expenditure *transactions.ExpenditureTransaction) string {
	score := expenditure.ComplianceScore

	switch {
	case score >= 90:
		return "Approve - Excellent compliance"
	case score >= 80:
		return "Approve - Good compliance with minor observations"
	case score >= 70:
		return "Approve with conditions - Address noted observations"
	case score >= 60:
		return "Conditional approval - Requires additional documentation"
	case score >= 50:
		return "Review required - Significant compliance gaps"
	default:
		return "Reject - Insufficient compliance and documentation"
	}
}

// GetAuditorStats returns comprehensive statistics about the auditor
func (a *Auditor) GetAuditorStats() AuditorStats {
	approvedCount := 0
	totalComplianceScore := 0.0

	for _, audit := range a.AuditHistory {
		if strings.Contains(strings.ToLower(audit.Recommendation), "approve") {
			approvedCount++
		}
		totalComplianceScore += audit.ComplianceScore
	}

	approvalRate := "0%"
	averageScore := 0.0

	if len(a.AuditHistory) > 0 {
		approvalRate = fmt.Sprintf("%.2f%%", (float64(approvedCount)/float64(len(a.AuditHistory)))*100)
		averageScore = totalComplianceScore / float64(len(a.AuditHistory))
	}

	return AuditorStats{
		AuditorID:              a.AuditorID,
		Name:                   a.Name,
		Verified:               a.Verified,
		Specializations:        a.Specializations,
		Rating:                 a.Rating,
		TotalAudits:            len(a.AuditHistory),
		ApprovedAudits:         approvedCount,
		ApprovalRate:           approvalRate,
		AverageComplianceScore: averageScore,
		MemberSince:            a.CreatedAt,
	}
}

// GetRecentAudits returns the most recent audits performed by this auditor
func (a *Auditor) GetRecentAudits(limit int) []AuditResult {
	if limit <= 0 {
		return a.AuditHistory
	}

	if len(a.AuditHistory) <= limit {
		return a.AuditHistory
	}

	// Return the last 'limit' audits
	start := len(a.AuditHistory) - limit
	return a.AuditHistory[start:]
}

// GetAuditsByDateRange returns audits within a specific date range
func (a *Auditor) GetAuditsByDateRange(startDate, endDate time.Time) []AuditResult {
	var audits []AuditResult

	for _, audit := range a.AuditHistory {
		if (audit.Timestamp.After(startDate) || audit.Timestamp.Equal(startDate)) &&
			(audit.Timestamp.Before(endDate) || audit.Timestamp.Equal(endDate)) {
			audits = append(audits, audit)
		}
	}

	return audits
}

// GetAuditByID finds an audit by its ID
func (a *Auditor) GetAuditByID(auditID string) *AuditResult {
	for _, audit := range a.AuditHistory {
		if audit.AuditID == auditID {
			return &audit
		}
	}
	return nil
}

// UpdateRating updates the auditor's rating based on performance metrics
func (a *Auditor) UpdateRating() float64 {
	if len(a.AuditHistory) == 0 {
		a.Rating = 5.0
		return a.Rating
	}

	// Base rating
	rating := 5.0

	// Calculate approval rate
	approvedCount := 0
	totalScore := 0.0

	for _, audit := range a.AuditHistory {
		if strings.Contains(strings.ToLower(audit.Recommendation), "approve") {
			approvedCount++
		}
		totalScore += audit.ComplianceScore
	}

	approvalRate := float64(approvedCount) / float64(len(a.AuditHistory))
	averageScore := totalScore / float64(len(a.AuditHistory))

	// Adjust rating based on approval rate
	if approvalRate > 0.8 {
		rating += 0.2
	} else if approvalRate < 0.3 {
		rating -= 0.5
	}

	// Adjust rating based on average compliance score
	if averageScore > 80 {
		rating += 0.3
	} else if averageScore < 60 {
		rating -= 0.3
	}

	// Experience bonus (more audits = slight bonus)
	if len(a.AuditHistory) > 100 {
		rating += 0.2
	} else if len(a.AuditHistory) > 50 {
		rating += 0.1
	}

	// Verification bonus
	if a.Verified {
		rating += 0.3
	}

	// Ensure rating stays within bounds
	if rating > 5.0 {
		rating = 5.0
	}
	if rating < 1.0 {
		rating = 1.0
	}

	a.Rating = rating
	return rating
}

// HasSpecialization checks if the auditor has a specific specialization
func (a *Auditor) HasSpecialization(specialization string) bool {
	for _, spec := range a.Specializations {
		if strings.EqualFold(spec, specialization) {
			return true
		}
	}
	return false
}

// AddSpecialization adds a new specialization to the auditor
func (a *Auditor) AddSpecialization(specialization string) {
	if !a.HasSpecialization(specialization) {
		a.Specializations = append(a.Specializations, specialization)
	}
}

// RemoveSpecialization removes a specialization from the auditor
func (a *Auditor) RemoveSpecialization(specialization string) {
	var newSpecs []string
	for _, spec := range a.Specializations {
		if !strings.EqualFold(spec, specialization) {
			newSpecs = append(newSpecs, spec)
		}
	}
	a.Specializations = newSpecs
}

// GetAuditSummary returns a summary of audit activities
func (a *Auditor) GetAuditSummary() map[string]interface{} {
	stats := a.GetAuditorStats()

	// Calculate monthly average
	monthsActive := int(time.Since(a.CreatedAt).Hours() / (24 * 30))
	if monthsActive < 1 {
		monthsActive = 1
	}
	monthlyAverage := float64(stats.TotalAudits) / float64(monthsActive)

	return map[string]interface{}{
		"auditor_id":             a.AuditorID,
		"name":                   a.Name,
		"verified":               a.Verified,
		"rating":                 a.Rating,
		"total_audits":           stats.TotalAudits,
		"approved_audits":        stats.ApprovedAudits,
		"approval_rate":          stats.ApprovalRate,
		"average_compliance":     fmt.Sprintf("%.1f%%", stats.AverageComplianceScore),
		"specializations":        a.Specializations,
		"monthly_average":        fmt.Sprintf("%.1f", monthlyAverage),
		"member_since":           a.CreatedAt.Format("2006-01-02"),
		"verification_authority": a.VerificationAuthority,
	}
}
