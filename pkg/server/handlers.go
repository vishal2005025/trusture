package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ngo-transparency-platform/pkg/auth"
	"ngo-transparency-platform/pkg/database"
	"ngo-transparency-platform/pkg/middleware"
)

// GetPlatformStatsHandler returns platform statistics
// @Summary Get platform statistics
// @Description Get overall platform statistics
// @Tags Public
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Router /api/v1/stats [get]
func (s *Server) GetPlatformStatsHandler(c *gin.Context) {
	stats := s.Platform.GetPlatformStats()
	middleware.StandardResponse(c, stats, "Platform statistics retrieved successfully")
}

// GetPublicNGOsHandler returns list of verified NGOs
// @Summary Get public NGO list
// @Description Get list of verified NGOs with basic information
// @Tags Public
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} middleware.PaginatedResponse
// @Router /api/v1/ngos [get]
func (s *Server) GetPublicNGOsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var ngos []database.NGOModel
	var total int64

	// Get verified NGOs
	query := database.DB.Where("kyc_verified = ?", true)
	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&ngos)

	middleware.PaginatedResponseData(c, ngos, page, limit, total)
}

// GetPublicNGOHandler returns public information about a specific NGO
// @Summary Get NGO public profile
// @Description Get public information about a specific NGO
// @Tags Public
// @Produce json
// @Param id path string true "NGO ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/v1/ngos/{id} [get]
func (s *Server) GetPublicNGOHandler(c *gin.Context) {
	ngoID := c.Param("id")
	
	var ngo database.NGOModel
	if err := database.DB.Where("ngo_id = ? AND kyc_verified = ?", ngoID, true).First(&ngo).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "ngo_not_found", "NGO not found", nil)
		return
	}

	// Get NGO rating and transparency data
	ratings := s.Platform.CalculateAllNGORatings(30)
	var ngoRating map[string]interface{}
	for _, rating := range ratings {
		if rating["ngo_id"] == ngoID {
			ngoRating = rating
			break
		}
	}

	response := gin.H{
		"ngo":    ngo,
		"rating": ngoRating,
	}

	middleware.StandardResponse(c, response, "NGO information retrieved successfully")
}

// GetNGORatingHandler returns NGO rating and transparency information
// @Summary Get NGO rating
// @Description Get NGO rating and transparency metrics
// @Tags Public
// @Produce json
// @Param id path string true "NGO ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/v1/ngos/{id}/rating [get]
func (s *Server) GetNGORatingHandler(c *gin.Context) {
	ngoID := c.Param("id")
	
	// Verify NGO exists
	var ngo database.NGOModel
	if err := database.DB.Where("ngo_id = ?", ngoID).First(&ngo).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "ngo_not_found", "NGO not found", nil)
		return
	}

	// Get rating data
	ratings := s.Platform.CalculateAllNGORatings(30)
	for _, rating := range ratings {
		if rating["ngo_id"] == ngoID {
			middleware.StandardResponse(c, rating, "NGO rating retrieved successfully")
			return
		}
	}

	middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "rating_not_found", "Rating data not found", nil)
}

// GetSystemStatusHandler returns system health status
// @Summary Get system status
// @Description Get system health and status information
// @Tags Public
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Router /api/v1/status [get]
func (s *Server) GetSystemStatusHandler(c *gin.Context) {
	status := gin.H{
		"service":     "Trusture API",
		"version":     "1.0.0",
		"status":      "healthy",
		"environment": s.Config.Platform.Environment,
		"database":    "connected",
		"blockchain":  "active",
	}

	// Check Polygon integration
	if s.Platform.PolygonIntegration != nil {
		stats := s.Platform.PolygonIntegration.GetNetworkStats()
		status["polygon"] = gin.H{
			"network":       stats.Network,
			"chain_id":      stats.ChainID,
			"wallet":        stats.WalletAddress,
			"contract":      stats.ContractAddress,
		}
	}

	middleware.StandardResponse(c, status, "System status retrieved successfully")
}

// VerifyBlockchainDataHandler verifies blockchain data by hash
// @Summary Verify blockchain data
// @Description Verify the authenticity of blockchain data by hash
// @Tags Public
// @Produce json
// @Param hash path string true "Block or transaction hash"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/v1/verify/{hash} [get]
func (s *Server) VerifyBlockchainDataHandler(c *gin.Context) {
	hash := c.Param("hash")
	
	// Check in local blockchain
	var block database.BlockchainBlockModel
	if err := database.DB.Where("hash = ?", hash).First(&block).Error; err == nil {
		verification := gin.H{
			"found":      true,
			"type":       "blockchain_block",
			"hash":       block.Hash,
			"index":      block.Index,
			"block_type": block.BlockType,
			"ngo_id":     block.NGOID,
			"validated":  block.Validated,
			"timestamp":  block.CreatedAt,
		}

		middleware.StandardResponse(c, verification, "Blockchain data verified successfully")
		return
	}

	// Check Polygon anchors if available
	if s.Platform.PolygonIntegration != nil {
		verification := s.Platform.PolygonIntegration.VerifyAnchoredHash(hash)
		if verification.Exists {
			response := gin.H{
				"found":         true,
				"type":          "polygon_anchor",
				"hash":          hash,
				"exists":        verification.Exists,
				"verified":      verification.Verified,
				"block_number":  verification.BlockNumber,
				"tx_hash":       verification.TxHash,
				"timestamp":     verification.Timestamp,
			}

			middleware.StandardResponse(c, response, "Polygon anchor verified successfully")
			return
		}
	}

	middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "hash_not_found", "Hash not found in blockchain records", nil)
}

// Protected endpoint handlers (require authentication)

// GetNGOProfileHandler returns NGO profile information
// @Summary Get NGO profile
// @Description Get NGO profile information (requires NGO authentication)
// @Tags NGO
// @Security Bearer
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/v1/ngos/profile [get]
func (s *Server) GetNGOProfileHandler(c *gin.Context) {
	userID, userType, entityID, err := auth.GetUserFromContext(c)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "unauthorized", "Unauthorized access", nil)
		return
	}

	if userType != "ngo" {
		middleware.ErrorResponseWithDetails(c, http.StatusForbidden, "forbidden", "Access denied", nil)
		return
	}

	var ngo database.NGOModel
	if err := database.DB.Where("ngo_id = ?", entityID).First(&ngo).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "ngo_not_found", "NGO profile not found", nil)
		return
	}

	// Get dashboard data
	dashboard, err := s.Platform.GetNGODashboard(entityID)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "dashboard_error", "Failed to retrieve dashboard data", nil)
		return
	}

	response := gin.H{
		"user_id":   userID,
		"profile":   ngo,
		"dashboard": dashboard,
	}

	middleware.StandardResponse(c, response, "NGO profile retrieved successfully")
}

// GetDonorProfileHandler returns Donor profile information
// @Summary Get Donor profile
// @Description Get Donor profile information (requires Donor authentication)
// @Tags Donor
// @Security Bearer
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/v1/donors/profile [get]
func (s *Server) GetDonorProfileHandler(c *gin.Context) {
	userID, userType, entityID, err := auth.GetUserFromContext(c)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "unauthorized", "Unauthorized access", nil)
		return
	}

	if userType != "donor" {
		middleware.ErrorResponseWithDetails(c, http.StatusForbidden, "forbidden", "Access denied", nil)
		return
	}

	var donor database.DonorModel
	if err := database.DB.Where("donor_id = ?", entityID).First(&donor).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "donor_not_found", "Donor profile not found", nil)
		return
	}

	// Get dashboard data
	dashboard, err := s.Platform.GetDonorDashboard(entityID)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "dashboard_error", "Failed to retrieve dashboard data", nil)
		return
	}

	response := gin.H{
		"user_id":   userID,
		"profile":   donor,
		"dashboard": dashboard,
	}

	middleware.StandardResponse(c, response, "Donor profile retrieved successfully")
}

// GetAuditorProfileHandler returns Auditor profile information
// @Summary Get Auditor profile
// @Description Get Auditor profile information (requires Auditor authentication)
// @Tags Auditor
// @Security Bearer
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/v1/auditors/profile [get]
func (s *Server) GetAuditorProfileHandler(c *gin.Context) {
	userID, userType, entityID, err := auth.GetUserFromContext(c)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "unauthorized", "Unauthorized access", nil)
		return
	}

	if userType != "auditor" {
		middleware.ErrorResponseWithDetails(c, http.StatusForbidden, "forbidden", "Access denied", nil)
		return
	}

	var auditor database.AuditorModel
	if err := database.DB.Where("auditor_id = ?", entityID).First(&auditor).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusNotFound, "auditor_not_found", "Auditor profile not found", nil)
		return
	}

	// Get dashboard data
	dashboard, err := s.Platform.GetAuditorDashboard(entityID)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "dashboard_error", "Failed to retrieve dashboard data", nil)
		return
	}

	response := gin.H{
		"user_id":   userID,
		"profile":   auditor,
		"dashboard": dashboard,
	}

	middleware.StandardResponse(c, response, "Auditor profile retrieved successfully")
}

// Placeholder handlers for remaining endpoints
// These would need to be implemented based on specific requirements

func (s *Server) GetNGODashboardHandler(c *gin.Context) {
	userID, userType, entityID, _ := auth.GetUserFromContext(c)
	dashboard, _ := s.Platform.GetNGODashboard(entityID)
	
	response := gin.H{
		"user_id":   userID,
		"user_type": userType,
		"entity_id": entityID,
		"dashboard": dashboard,
	}
	
	middleware.StandardResponse(c, response, "NGO dashboard retrieved successfully")
}

func (s *Server) GetDonorDashboardHandler(c *gin.Context) {
	userID, userType, entityID, _ := auth.GetUserFromContext(c)
	dashboard, _ := s.Platform.GetDonorDashboard(entityID)
	
	response := gin.H{
		"user_id":   userID,
		"user_type": userType,
		"entity_id": entityID,
		"dashboard": dashboard,
	}
	
	middleware.StandardResponse(c, response, "Donor dashboard retrieved successfully")
}

func (s *Server) GetAuditorDashboardHandler(c *gin.Context) {
	userID, userType, entityID, _ := auth.GetUserFromContext(c)
	dashboard, _ := s.Platform.GetAuditorDashboard(entityID)
	
	response := gin.H{
		"user_id":   userID,
		"user_type": userType,
		"entity_id": entityID,
		"dashboard": dashboard,
	}
	
	middleware.StandardResponse(c, response, "Auditor dashboard retrieved successfully")
}

// Placeholder handlers - to be implemented as needed
func (s *Server) UpdateNGOProfileHandler(c *gin.Context)        { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) UpdateDonorProfileHandler(c *gin.Context)      { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) UpdateAuditorProfileHandler(c *gin.Context)    { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetNGODonationsHandler(c *gin.Context)         { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetNGOExpendituresHandler(c *gin.Context)      { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) CreateExpenditureHandler(c *gin.Context)       { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetExpenditureHandler(c *gin.Context)          { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) UpdateExpenditureHandler(c *gin.Context)       { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetNGODonationBlocksHandler(c *gin.Context)    { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetNGOExpenditureBlocksHandler(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) SubmitNGOKYCHandler(c *gin.Context)            { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetNGOFinancialSummaryHandler(c *gin.Context)  { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetDonorDonationsHandler(c *gin.Context)       { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) CreateDonationHandler(c *gin.Context)          { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetDonationHandler(c *gin.Context)             { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetTaxBenefitsHandler(c *gin.Context)          { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetPreferredNGOsHandler(c *gin.Context)        { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) AddPreferredNGOHandler(c *gin.Context)         { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) RemovePreferredNGOHandler(c *gin.Context)      { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) SubmitDonorKYCHandler(c *gin.Context)          { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) CheckDonationLimitHandler(c *gin.Context)      { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetAuditorAuditsHandler(c *gin.Context)        { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetPendingExpendituresHandler(c *gin.Context)  { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) AuditExpenditureHandler(c *gin.Context)        { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetAuditHandler(c *gin.Context)                { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) SubmitAuditorKYCHandler(c *gin.Context)        { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetDonationTransactionHandler(c *gin.Context)  { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetExpenditureTransactionHandler(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetDonationReceiptHandler(c *gin.Context)      { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetComplianceReportHandler(c *gin.Context)     { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetBlockHandler(c *gin.Context)                { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) VerifyBlockHandler(c *gin.Context)             { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetPolygonAnchorsHandler(c *gin.Context)       { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) GetPolygonStatsHandler(c *gin.Context)         { c.JSON(501, gin.H{"error": "Not implemented yet"}) }
func (s *Server) AnchorBlockToPolygonHandler(c *gin.Context)    { c.JSON(501, gin.H{"error": "Not implemented yet"}) }