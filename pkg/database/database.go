package database

import (
	"fmt"
	"log"

	"ngo-transparency-platform/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Set log level based on environment
	if config.IsProduction() {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(300) // 5 minutes

	log.Println("Database connected successfully")
	return nil
}

// MigrateDatabase runs all database migrations
func MigrateDatabase() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Auto-migrate all models
	err := DB.AutoMigrate(
		&User{},
		&NGOModel{},
		&DonorModel{},
		&AuditorModel{},
		&DonationModel{},
		&ExpenditureModel{},
		&AuditModel{},
		&BlockchainBlockModel{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Repository interface for common database operations
type Repository interface {
	Create(entity interface{}) error
	GetByID(id uint, entity interface{}) error
	Update(entity interface{}) error
	Delete(id uint, entity interface{}) error
	List(entities interface{}, limit, offset int) error
}

// BaseRepository implements common database operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{db: DB}
}

func (r *BaseRepository) Create(entity interface{}) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository) GetByID(id uint, entity interface{}) error {
	return r.db.First(entity, id).Error
}

func (r *BaseRepository) Update(entity interface{}) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository) Delete(id uint, entity interface{}) error {
	return r.db.Delete(entity, id).Error
}

func (r *BaseRepository) List(entities interface{}, limit, offset int) error {
	query := r.db
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(entities).Error
}

// Specialized repositories

// NGORepository handles NGO-specific database operations
type NGORepository struct {
	*BaseRepository
}

func NewNGORepository() *NGORepository {
	return &NGORepository{NewBaseRepository()}
}

func (r *NGORepository) GetByNGOID(ngoID string) (*NGOModel, error) {
	var ngo NGOModel
	err := r.db.Where("ngo_id = ?", ngoID).First(&ngo).Error
	return &ngo, err
}

func (r *NGORepository) GetByRegistrationNumber(regNumber string) (*NGOModel, error) {
	var ngo NGOModel
	err := r.db.Where("registration_number = ?", regNumber).First(&ngo).Error
	return &ngo, err
}

func (r *NGORepository) GetVerifiedNGOs() ([]NGOModel, error) {
	var ngos []NGOModel
	err := r.db.Where("kyc_verified = ?", true).Find(&ngos).Error
	return ngos, err
}

func (r *NGORepository) UpdateRating(ngoID string, rating float64, transparencyScore int) error {
	return r.db.Model(&NGOModel{}).Where("ngo_id = ?", ngoID).
		Updates(map[string]interface{}{
			"rating": rating,
			"transparency_score": transparencyScore,
		}).Error
}

// DonorRepository handles Donor-specific database operations
type DonorRepository struct {
	*BaseRepository
}

func NewDonorRepository() *DonorRepository {
	return &DonorRepository{NewBaseRepository()}
}

func (r *DonorRepository) GetByDonorID(donorID string) (*DonorModel, error) {
	var donor DonorModel
	err := r.db.Where("donor_id = ?", donorID).First(&donor).Error
	return &donor, err
}

func (r *DonorRepository) UpdateTotalDonated(donorID string, totalDonated float64) error {
	return r.db.Model(&DonorModel{}).Where("donor_id = ?", donorID).
		Update("total_donated", totalDonated).Error
}

// AuditorRepository handles Auditor-specific database operations
type AuditorRepository struct {
	*BaseRepository
}

func NewAuditorRepository() *AuditorRepository {
	return &AuditorRepository{NewBaseRepository()}
}

func (r *AuditorRepository) GetByAuditorID(auditorID string) (*AuditorModel, error) {
	var auditor AuditorModel
	err := r.db.Where("auditor_id = ?", auditorID).First(&auditor).Error
	return &auditor, err
}

func (r *AuditorRepository) GetVerifiedAuditors() ([]AuditorModel, error) {
	var auditors []AuditorModel
	err := r.db.Where("verified = ?", true).Find(&auditors).Error
	return auditors, err
}

// DonationRepository handles Donation-specific database operations
type DonationRepository struct {
	*BaseRepository
}

func NewDonationRepository() *DonationRepository {
	return &DonationRepository{NewBaseRepository()}
}

func (r *DonationRepository) GetByTransactionID(transactionID string) (*DonationModel, error) {
	var donation DonationModel
	err := r.db.Where("transaction_id = ?", transactionID).First(&donation).Error
	return &donation, err
}

func (r *DonationRepository) GetByDonorID(donorID string, limit, offset int) ([]DonationModel, error) {
	var donations []DonationModel
	query := r.db.Where("donor_id = ?", donorID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&donations).Error
	return donations, err
}

func (r *DonationRepository) GetByNGOID(ngoID string, limit, offset int) ([]DonationModel, error) {
	var donations []DonationModel
	query := r.db.Where("ngo_id = ?", ngoID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&donations).Error
	return donations, err
}

// ExpenditureRepository handles Expenditure-specific database operations
type ExpenditureRepository struct {
	*BaseRepository
}

func NewExpenditureRepository() *ExpenditureRepository {
	return &ExpenditureRepository{NewBaseRepository()}
}

func (r *ExpenditureRepository) GetByTransactionID(transactionID string) (*ExpenditureModel, error) {
	var expenditure ExpenditureModel
	err := r.db.Where("transaction_id = ?", transactionID).First(&expenditure).Error
	return &expenditure, err
}

func (r *ExpenditureRepository) GetByNGOID(ngoID string, limit, offset int) ([]ExpenditureModel, error) {
	var expenditures []ExpenditureModel
	query := r.db.Where("ngo_id = ?", ngoID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&expenditures).Error
	return expenditures, err
}

func (r *ExpenditureRepository) GetPendingValidation() ([]ExpenditureModel, error) {
	var expenditures []ExpenditureModel
	err := r.db.Where("status = ?", "pending_validation").Find(&expenditures).Error
	return expenditures, err
}