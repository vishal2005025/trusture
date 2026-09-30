package transactions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"ngo-transparency-platform/pkg/crypto"
	"strings"
	"time"
)

// TaxBenefit represents tax benefit information
type TaxBenefit struct {
	Section          string  `json:"section"`
	DeductibleAmount float64 `json:"deductible_amount"`
	TaxSaving        float64 `json:"tax_saving"`
	Note             string  `json:"note"`
}

// EBill represents an electronic bill/receipt
type EBill struct {
	BillID           string     `json:"bill_id"`
	TransactionID    string     `json:"transaction_id"`
	Amount           float64    `json:"amount"`
	Currency         string     `json:"currency"`
	Timestamp        time.Time  `json:"timestamp"`
	NGOID            string     `json:"ngo_id"`
	DonorHash        string     `json:"donor_hash"`
	PaymentMethod    string     `json:"payment_method"`
	TaxBenefit       TaxBenefit `json:"tax_benefit"`
	ReceiptNumber    string     `json:"receipt_number"`
	Signature        string     `json:"signature"`
	QRCode           string     `json:"qr_code"`
	DownloadURL      string     `json:"download_url"`
	ValidityPeriod   string     `json:"validity_period"`
}

// DonationTransaction represents a donation transaction
type DonationTransaction struct {
	TransactionID   string            `json:"transaction_id"`
	DonorID         string            `json:"donor_id"`
	NGOID           string            `json:"ngo_id"`
	Amount          float64           `json:"amount"`
	PaymentMethod   string            `json:"payment_method"`
	Timestamp       time.Time         `json:"timestamp"`
	Status          string            `json:"status"`
	DonorKYCHash    string            `json:"donor_kyc_hash"`
	ZKProof       *crypto.ZKProof            `json:"zk_proof"`
	EBill           *EBill            `json:"e_bill"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	FailedAt        *time.Time        `json:"failed_at,omitempty"`
	FailureReason   string            `json:"failure_reason,omitempty"`
}

// NewDonationTransaction creates a new donation transaction
func NewDonationTransaction(donorID, ngoID string, amount float64, paymentMethod, donorKYCHash string) *DonationTransaction {
	// Generate transaction ID
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	transactionID := hex.EncodeToString(randomBytes)

	timestamp := time.Now()

	transaction := &DonationTransaction{
		TransactionID: transactionID,
		DonorID:       donorID,
		NGOID:         ngoID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Timestamp:     timestamp,
		Status:        "pending",
		DonorKYCHash:  donorKYCHash,
	}

	// Generate ZK proof
	transaction.ZKProof = crypto.GenerateProof(donorID, amount, timestamp)

	// Generate e-bill
	transaction.EBill = transaction.generateEBill()

	return transaction
}

// generateEBill creates an electronic bill for the donation
func (dt *DonationTransaction) generateEBill() *EBill {
	// Generate bill ID
	randomBytes := make([]byte, 12)
	rand.Read(randomBytes)
	billID := hex.EncodeToString(randomBytes)

	// Generate donor hash (anonymized)
	donorHashBytes := sha256.Sum256([]byte(dt.DonorID))
	donorHash := hex.EncodeToString(donorHashBytes[:])

	// Calculate tax benefits
	taxBenefit := dt.calculateTaxBenefit()

	// Generate receipt number
	receiptNumber := fmt.Sprintf("RCP-%d-%s",
		dt.Timestamp.Unix(),
		strings.ToUpper(hex.EncodeToString(randomBytes[:5])))

	billData := &EBill{
		BillID:         billID,
		TransactionID:  dt.TransactionID,
		Amount:         dt.Amount,
		Currency:       "INR",
		Timestamp:      dt.Timestamp,
		NGOID:          dt.NGOID,
		DonorHash:      donorHash,
		PaymentMethod:  dt.PaymentMethod,
		TaxBenefit:     taxBenefit,
		ReceiptNumber:  receiptNumber,
		DownloadURL:    fmt.Sprintf("https://receipts.ngo/%s", billID),
		ValidityPeriod: "7 years", // For tax purposes
	}

	// Generate signature
	billData.Signature = dt.generateBillSignature(billData)

	// Generate QR code
	billData.QRCode = dt.generateQRCode(billData)

	return billData
}

// calculateTaxBenefit calculates the tax benefit for the donation
func (dt *DonationTransaction) calculateTaxBenefit() TaxBenefit {
	// 80G deduction calculation (simplified)
	maxDeduction := math.Min(dt.Amount, 10000) // Simplified calculation

	return TaxBenefit{
		Section:          "80G",
		DeductibleAmount: maxDeduction,
		TaxSaving:        maxDeduction * 0.3, // Assuming 30% tax bracket
		Note:             "Consult tax advisor for accurate calculations",
	}
}

// generateQRCode generates a QR code for the bill
func (dt *DonationTransaction) generateQRCode(billData *EBill) string {
	qrData := map[string]interface{}{
		"bill_id":   billData.BillID,
		"amount":    billData.Amount,
		"ngo_id":    billData.NGOID,
		"timestamp": billData.Timestamp.Unix(),
	}

	qrBytes, _ := json.Marshal(qrData)
	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

	return fmt.Sprintf("QR:%s", qrBase64)
}

// generateBillSignature generates a signature for the e-bill
func (dt *DonationTransaction) generateBillSignature(billData *EBill) string {
	// Create a copy of bill data without signature, QR code, and download URL
	signatureData := map[string]interface{}{
		"bill_id":         billData.BillID,
		"transaction_id":  billData.TransactionID,
		"amount":          billData.Amount,
		"currency":        billData.Currency,
		"timestamp":       billData.Timestamp.Unix(),
		"ngo_id":          billData.NGOID,
		"donor_hash":      billData.DonorHash,
		"payment_method":  billData.PaymentMethod,
		"tax_benefit":     billData.TaxBenefit,
		"receipt_number":  billData.ReceiptNumber,
		"validity_period": billData.ValidityPeriod,
	}

	signatureBytes, _ := json.Marshal(signatureData)
	hash := sha256.Sum256(signatureBytes)
	return hex.EncodeToString(hash[:])
}

// ValidateEBill validates the e-bill integrity
func (dt *DonationTransaction) ValidateEBill() bool {
	if dt.EBill == nil {
		return false
	}

	// Create a copy and recalculate signature
	originalSignature := dt.EBill.Signature
	recalculatedSignature := dt.generateBillSignature(dt.EBill)

	return originalSignature == recalculatedSignature
}

// MarkComplete marks the transaction as completed
func (dt *DonationTransaction) MarkComplete() {
	dt.Status = "completed"
	now := time.Now()
	dt.CompletedAt = &now
}

// MarkFailed marks the transaction as failed
func (dt *DonationTransaction) MarkFailed(reason string) {
	dt.Status = "failed"
	now := time.Now()
	dt.FailedAt = &now
	dt.FailureReason = reason
}

// IsCompleted checks if the transaction is completed
func (dt *DonationTransaction) IsCompleted() bool {
	return dt.Status == "completed"
}

// IsFailed checks if the transaction is failed
func (dt *DonationTransaction) IsFailed() bool {
	return dt.Status == "failed"
}

// IsPending checks if the transaction is pending
func (dt *DonationTransaction) IsPending() bool {
	return dt.Status == "pending"
}

// GetTransactionSummary returns a summary of the transaction
func (dt *DonationTransaction) GetTransactionSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"transaction_id":  dt.TransactionID,
		"ngo_id":          dt.NGOID,
		"amount":          dt.Amount,
		"currency":        "INR",
		"payment_method":  dt.PaymentMethod,
		"status":          dt.Status,
		"timestamp":       dt.Timestamp,
		"has_zk_proof":    dt.ZKProof != nil,
		"has_e_bill":      dt.EBill != nil,
	}

	if dt.CompletedAt != nil {
		summary["completed_at"] = *dt.CompletedAt
	}

	if dt.FailedAt != nil {
		summary["failed_at"] = *dt.FailedAt
		summary["failure_reason"] = dt.FailureReason
	}

	if dt.EBill != nil {
		summary["receipt_number"] = dt.EBill.ReceiptNumber
		summary["tax_saving"] = dt.EBill.TaxBenefit.TaxSaving
	}

	return summary
}

// GetEBillInfo returns e-bill information
func (dt *DonationTransaction) GetEBillInfo() map[string]interface{} {
	if dt.EBill == nil {
		return nil
	}

	return map[string]interface{}{
		"bill_id":         dt.EBill.BillID,
		"receipt_number":  dt.EBill.ReceiptNumber,
		"amount":          dt.EBill.Amount,
		"currency":        dt.EBill.Currency,
		"tax_benefit":     dt.EBill.TaxBenefit,
		"download_url":    dt.EBill.DownloadURL,
		"validity_period": dt.EBill.ValidityPeriod,
		"qr_code":         dt.EBill.QRCode,
	}
}
