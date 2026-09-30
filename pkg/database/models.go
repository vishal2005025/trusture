package database

import (
	"time"
	"encoding/json"
)

// User represents the base user model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // Never return password in JSON
	UserType  string    `json:"user_type" gorm:"not null"` // "ngo", "donor", "auditor"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NGOModel represents the database model for NGOs
type NGOModel struct {
	ID                     uint      `json:"id" gorm:"primaryKey"`
	UserID                 uint      `json:"user_id" gorm:"not null"`
	User                   User      `json:"user" gorm:"foreignKey:UserID"`
	NGOID                  string    `json:"ngo_id" gorm:"unique;not null"`
	Name                   string    `json:"name" gorm:"not null"`
	RegistrationNumber     string    `json:"registration_number" gorm:"unique;not null"`
	Category               string    `json:"category" gorm:"not null"`
	Rating                 float64   `json:"rating" gorm:"default:5.0"`
	KYCVerified            bool      `json:"kyc_verified" gorm:"default:false"`
	KYCData                string    `json:"kyc_data" gorm:"type:text"` // JSON string
	TotalDonationsReceived float64   `json:"total_donations_received" gorm:"default:0"`
	TotalExpenditureReported float64 `json:"total_expenditure_reported" gorm:"default:0"`
	TransparencyScore      int       `json:"transparency_score" gorm:"default:100"`
	PublicKey              string    `json:"public_key" gorm:"not null"`
	Certificates           string    `json:"certificates" gorm:"type:text"` // JSON string
	MultiSigSigners        string    `json:"multisig_signers" gorm:"type:text"` // JSON string
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// DonorModel represents the database model for Donors
type DonorModel struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	UserID              uint      `json:"user_id" gorm:"not null"`
	User                User      `json:"user" gorm:"foreignKey:UserID"`
	DonorID             string    `json:"donor_id" gorm:"unique;not null"`
	KYCVerified         bool      `json:"kyc_verified" gorm:"default:false"`
	KYCData             string    `json:"kyc_data" gorm:"type:text"` // JSON string
	TotalDonated        float64   `json:"total_donated" gorm:"default:0"`
	PreferredNGOs       string    `json:"preferred_ngos" gorm:"type:text"` // JSON string
	AnnualDonationLimit float64   `json:"annual_donation_limit" gorm:"default:1000000"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// AuditorModel represents the database model for Auditors
type AuditorModel struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	AuditorID       string    `json:"auditor_id" gorm:"unique;not null"`
	Name            string    `json:"name" gorm:"not null"`
	Credentials     string    `json:"credentials" gorm:"type:text"` // JSON string
	Specializations string    `json:"specializations" gorm:"type:text"` // JSON string
	Verified        bool      `json:"verified" gorm:"default:false"`
	Rating          float64   `json:"rating" gorm:"default:5.0"`
	PublicKey       string    `json:"public_key" gorm:"not null"`
	VerificationAuthority string `json:"verification_authority"`
	VerificationDate      *time.Time `json:"verification_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DonationModel represents the database model for Donations
type DonationModel struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TransactionID string    `json:"transaction_id" gorm:"unique;not null"`
	DonorID       string    `json:"donor_id" gorm:"not null"`
	NGOID         string    `json:"ngo_id" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	PlatformFee   float64   `json:"platform_fee" gorm:"default:0"`
	NetAmount     float64   `json:"net_amount" gorm:"not null"`
	PaymentMethod string    `json:"payment_method" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null;default:pending"`
	BlockHash     string    `json:"block_hash"`
	PolygonTxHash string    `json:"polygon_tx_hash"`
	ZKProofData   string    `json:"zk_proof_data" gorm:"type:text"` // JSON string
	EBillData     string    `json:"e_bill_data" gorm:"type:text"` // JSON string
	TaxBenefit    string    `json:"tax_benefit" gorm:"type:text"` // JSON string
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// ExpenditureModel represents the database model for Expenditures
type ExpenditureModel struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TransactionID     string    `json:"transaction_id" gorm:"unique;not null"`
	NGOID             string    `json:"ngo_id" gorm:"not null"`
	Amount            float64   `json:"amount" gorm:"not null"`
	Category          string    `json:"category" gorm:"not null"`
	Description       string    `json:"description" gorm:"type:text"`
	Status            string    `json:"status" gorm:"not null;default:pending_validation"`
	InvoiceDetails    string    `json:"invoice_details" gorm:"type:text"` // JSON string
	Attachments       string    `json:"attachments" gorm:"type:text"` // JSON string
	AuditorValidation string    `json:"auditor_validation" gorm:"type:text"` // JSON string
	ComplianceScore   float64   `json:"compliance_score" gorm:"default:0"`
	BlockHash         string    `json:"block_hash"`
	PolygonTxHash     string    `json:"polygon_tx_hash"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AuditModel represents the database model for Audits
type AuditModel struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	AuditID         string    `json:"audit_id" gorm:"unique;not null"`
	ExpenditureID   string    `json:"expenditure_id" gorm:"not null"`
	AuditorID       string    `json:"auditor_id" gorm:"not null"`
	ComplianceScore float64   `json:"compliance_score" gorm:"not null"`
	Findings        string    `json:"findings" gorm:"type:text"` // JSON string
	Recommendation  string    `json:"recommendation" gorm:"not null"`
	AuditNotes      string    `json:"audit_notes" gorm:"type:text"`
	Signature       string    `json:"signature" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BlockchainBlockModel represents blockchain blocks in database
type BlockchainBlockModel struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Index        int       `json:"index" gorm:"not null"`
	Hash         string    `json:"hash" gorm:"unique;not null"`
	PreviousHash string    `json:"previous_hash" gorm:"not null"`
	BlockType    string    `json:"block_type" gorm:"not null"` // donation, expenditure
	NGOID        string    `json:"ngo_id" gorm:"not null"`
	Data         string    `json:"data" gorm:"type:text"` // JSON string
	MerkleRoot   string    `json:"merkle_root" gorm:"not null"`
	Nonce        int       `json:"nonce" gorm:"default:0"`
	Validated    bool      `json:"validated" gorm:"default:false"`
	Validators   string    `json:"validators" gorm:"type:text"` // JSON string
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName methods to customize table names
func (NGOModel) TableName() string {
	return "ngos"
}

func (DonorModel) TableName() string {
	return "donors"
}

func (AuditorModel) TableName() string {
	return "auditors"
}

func (DonationModel) TableName() string {
	return "donations"
}

func (ExpenditureModel) TableName() string {
	return "expenditures"
}

func (AuditModel) TableName() string {
	return "audits"
}

func (BlockchainBlockModel) TableName() string {
	return "blockchain_blocks"
}

// Helper methods for JSON marshaling/unmarshaling
func (n *NGOModel) SetKYCData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	n.KYCData = string(jsonData)
	return nil
}

func (n *NGOModel) GetKYCData() (map[string]interface{}, error) {
	var data map[string]interface{}
	if n.KYCData == "" {
		return data, nil
	}
	err := json.Unmarshal([]byte(n.KYCData), &data)
	return data, err
}

func (d *DonorModel) SetKYCData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.KYCData = string(jsonData)
	return nil
}

func (d *DonorModel) GetKYCData() (map[string]interface{}, error) {
	var data map[string]interface{}
	if d.KYCData == "" {
		return data, nil
	}
	err := json.Unmarshal([]byte(d.KYCData), &data)
	return data, err
}