package server

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ngo-transparency-platform/pkg/auth"
	"ngo-transparency-platform/pkg/config"
	"ngo-transparency-platform/pkg/database"
	"ngo-transparency-platform/pkg/middleware"
	"ngo-transparency-platform/pkg/platform"
)

// Server represents the HTTP server
type Server struct {
	Config   *config.Config
	Router   *gin.Engine
	Platform *platform.NGOTransparencyPlatform
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	return &Server{
		Config: cfg,
		Router: gin.New(),
	}
}

// InitializePlatform initializes the core platform
func (s *Server) InitializePlatform() error {
	// Initialize the transparency platform
	s.Platform = platform.NewNGOTransparencyPlatform()

	// Initialize Polygon integration if configured
	if s.Config.Blockchain.PolygonRPC != "" {
		gasPrice := big.NewInt(s.Config.Blockchain.GasPriceGwei)
		gasPrice.Mul(gasPrice, big.NewInt(1e9)) // Convert Gwei to Wei

		s.Platform.InitializePolygon(
			s.Config.Blockchain.PolygonRPC,
			s.Config.Blockchain.PrivateKey,
			s.Config.Blockchain.GasLimit,
			gasPrice,
		)

		log.Println("Polygon integration initialized")
	}

	return nil
}

// SetupMiddleware configures all middleware
func (s *Server) SetupMiddleware() {
	// Initialize logger
	middleware.InitLogger(s.Config)

	// Recovery and error handling
	s.Router.Use(middleware.ErrorHandler())
	
	// Request logging
	s.Router.Use(middleware.RequestLogger())
	
	// CORS
	s.Router.Use(middleware.CORSMiddleware())
	
	// Security headers
	s.Router.Use(middleware.SecurityHeaders())
	
	// Rate limiting (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	s.Router.Use(rateLimiter.Middleware())
	
	// Health check
	s.Router.Use(middleware.HealthCheck())
	
	// Content-Type validation for non-GET requests
	s.Router.Use(middleware.JSONContentType())
}

// SetupRoutes configures all API routes
func (s *Server) SetupRoutes() {
	// Root route
	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Trusture API",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/swagger/index.html",
		})
	})

	// API v1 routes
	v1 := s.Router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("")
		{
			s.setupAuthRoutes(public)
			s.setupPublicRoutes(public)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			s.setupNGORoutes(protected)
			s.setupDonorRoutes(protected)
			s.setupAuditorRoutes(protected)
			s.setupTransactionRoutes(protected)
			s.setupBlockchainRoutes(protected)
		}
	}

	// Swagger documentation
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

// setupAuthRoutes sets up authentication routes
func (s *Server) setupAuthRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", s.RegisterHandler)
		authGroup.POST("/login", s.LoginHandler)
		authGroup.POST("/refresh", auth.AuthMiddleware(), s.RefreshTokenHandler)
		authGroup.POST("/logout", auth.AuthMiddleware(), s.LogoutHandler)
	}
}

// setupPublicRoutes sets up public routes that don't require authentication
func (s *Server) setupPublicRoutes(router *gin.RouterGroup) {
	// NGO public information
	router.GET("/ngos", s.GetPublicNGOsHandler)
	router.GET("/ngos/:id", s.GetPublicNGOHandler)
	router.GET("/ngos/:id/rating", s.GetNGORatingHandler)
	
	// Platform statistics
	router.GET("/stats", s.GetPlatformStatsHandler)
	
	// Blockchain verification
	router.GET("/verify/:hash", s.VerifyBlockchainDataHandler)
	
	// Health and status
	router.GET("/status", s.GetSystemStatusHandler)
}

// setupNGORoutes sets up NGO-specific routes
func (s *Server) setupNGORoutes(router *gin.RouterGroup) {
	ngoGroup := router.Group("/ngos")
	ngoGroup.Use(auth.RequireUserType("ngo"))
	{
		ngoGroup.GET("/profile", s.GetNGOProfileHandler)
		ngoGroup.PUT("/profile", s.UpdateNGOProfileHandler)
		ngoGroup.GET("/dashboard", s.GetNGODashboardHandler)
		ngoGroup.GET("/donations", s.GetNGODonationsHandler)
		ngoGroup.GET("/expenditures", s.GetNGOExpendituresHandler)
		ngoGroup.POST("/expenditures", s.CreateExpenditureHandler)
		ngoGroup.GET("/expenditures/:id", s.GetExpenditureHandler)
		ngoGroup.PUT("/expenditures/:id", s.UpdateExpenditureHandler)
		ngoGroup.GET("/blockchain/donations", s.GetNGODonationBlocksHandler)
		ngoGroup.GET("/blockchain/expenditures", s.GetNGOExpenditureBlocksHandler)
		ngoGroup.POST("/kyc/submit", s.SubmitNGOKYCHandler)
		ngoGroup.GET("/financial-summary", s.GetNGOFinancialSummaryHandler)
	}
}

// setupDonorRoutes sets up Donor-specific routes
func (s *Server) setupDonorRoutes(router *gin.RouterGroup) {
	donorGroup := router.Group("/donors")
	donorGroup.Use(auth.RequireUserType("donor"))
	{
		donorGroup.GET("/profile", s.GetDonorProfileHandler)
		donorGroup.PUT("/profile", s.UpdateDonorProfileHandler)
		donorGroup.GET("/dashboard", s.GetDonorDashboardHandler)
		donorGroup.GET("/donations", s.GetDonorDonationsHandler)
		donorGroup.POST("/donations", s.CreateDonationHandler)
		donorGroup.GET("/donations/:id", s.GetDonationHandler)
		donorGroup.GET("/tax-benefits", s.GetTaxBenefitsHandler)
		donorGroup.GET("/preferred-ngos", s.GetPreferredNGOsHandler)
		donorGroup.POST("/preferred-ngos/:ngo_id", s.AddPreferredNGOHandler)
		donorGroup.DELETE("/preferred-ngos/:ngo_id", s.RemovePreferredNGOHandler)
		donorGroup.POST("/kyc/submit", s.SubmitDonorKYCHandler)
		donorGroup.GET("/limit-check", s.CheckDonationLimitHandler)
	}
}

// setupAuditorRoutes sets up Auditor-specific routes
func (s *Server) setupAuditorRoutes(router *gin.RouterGroup) {
	auditorGroup := router.Group("/auditors")
	auditorGroup.Use(auth.RequireUserType("auditor"))
	{
		auditorGroup.GET("/profile", s.GetAuditorProfileHandler)
		auditorGroup.PUT("/profile", s.UpdateAuditorProfileHandler)
		auditorGroup.GET("/dashboard", s.GetAuditorDashboardHandler)
		auditorGroup.GET("/audits", s.GetAuditorAuditsHandler)
		auditorGroup.GET("/pending-expenditures", s.GetPendingExpendituresHandler)
		auditorGroup.POST("/audit/:expenditure_id", s.AuditExpenditureHandler)
		auditorGroup.GET("/audits/:id", s.GetAuditHandler)
		auditorGroup.POST("/kyc/submit", s.SubmitAuditorKYCHandler)
	}
}

// setupTransactionRoutes sets up transaction-related routes
func (s *Server) setupTransactionRoutes(router *gin.RouterGroup) {
	txGroup := router.Group("/transactions")
	{
		txGroup.GET("/donations/:id", s.GetDonationTransactionHandler)
		txGroup.GET("/expenditures/:id", s.GetExpenditureTransactionHandler)
		txGroup.GET("/donations/:id/receipt", s.GetDonationReceiptHandler)
		txGroup.GET("/expenditures/:id/compliance", s.GetComplianceReportHandler)
	}
}

// setupBlockchainRoutes sets up blockchain-related routes
func (s *Server) setupBlockchainRoutes(router *gin.RouterGroup) {
	blockchainGroup := router.Group("/blockchain")
	{
		blockchainGroup.GET("/blocks/:hash", s.GetBlockHandler)
		blockchainGroup.GET("/verify/:hash", auth.OptionalAuth(), s.VerifyBlockHandler)
		blockchainGroup.GET("/polygon/anchors", s.GetPolygonAnchorsHandler)
		blockchainGroup.GET("/polygon/stats", s.GetPolygonStatsHandler)
		blockchainGroup.POST("/anchor/:block_hash", auth.RequireUserType("ngo"), s.AnchorBlockToPolygonHandler)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port)
	
	middleware.Logger.WithFields(map[string]interface{}{
		"address":     address,
		"environment": s.Config.Platform.Environment,
		"version":     "1.0.0",
	}).Info("Starting Trusture API server")

	return s.Router.Run(address)
}

// Graceful shutdown
func (s *Server) Shutdown() error {
	middleware.Logger.Info("Shutting down server...")
	
	// Close database connections
	if err := database.CloseDatabase(); err != nil {
		middleware.Logger.WithError(err).Error("Error closing database")
		return err
	}
	
	middleware.Logger.Info("Server shutdown complete")
	return nil
}

// Initialize initializes the server with all dependencies
func (s *Server) Initialize() error {
	// Initialize database
	if err := database.InitDatabase(s.Config); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run database migrations
	if err := database.MigrateDatabase(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize platform
	if err := s.InitializePlatform(); err != nil {
		return fmt.Errorf("failed to initialize platform: %w", err)
	}

	// Setup middleware
	s.SetupMiddleware()

	// Setup routes
	s.SetupRoutes()

	middleware.Logger.Info("Server initialized successfully")
	return nil
}