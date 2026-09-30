package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ngo-transparency-platform/pkg/auth"
	"ngo-transparency-platform/pkg/database"
	"ngo-transparency-platform/pkg/middleware"
)

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=6"`
	UserType             string `json:"user_type" binding:"required,oneof=ngo donor auditor"`
	Name                 string `json:"name" binding:"required"`
	RegistrationNumber   string `json:"registration_number,omitempty"` // For NGOs
	Category             string `json:"category,omitempty"`            // For NGOs
	Specializations      []string `json:"specializations,omitempty"`   // For Auditors
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Register a new user (NGO, Donor, or Auditor)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Router /api/v1/auth/register [post]
func (s *Server) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusBadRequest, "validation_error", "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser database.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		middleware.ErrorResponseWithDetails(c, http.StatusConflict, "user_exists", "User with this email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "password_hash_error", "Failed to hash password", nil)
		return
	}

	// Create user
	user := database.User{
		Email:    req.Email,
		Password: hashedPassword,
		UserType: req.UserType,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "database_error", "Failed to create user", nil)
		return
	}

	// Create entity based on user type
	var entityID string
	switch req.UserType {
	case "ngo":
		entityID, err = s.createNGOEntity(user.ID, req)
	case "donor":
		entityID, err = s.createDonorEntity(user.ID, req)
	case "auditor":
		entityID, err = s.createAuditorEntity(user.ID, req)
	}

	if err != nil {
		// Rollback user creation if entity creation fails
		database.DB.Delete(&user)
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "entity_creation_error", "Failed to create entity", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Generate JWT token
	tokenResponse, err := auth.GenerateToken(user.ID, user.Email, user.UserType, entityID)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "token_generation_error", "Failed to generate token", nil)
		return
	}

	middleware.StandardResponse(c, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"user_type": user.UserType,
			"entity_id": entityID,
		},
		"token": tokenResponse,
	}, "User registered successfully")
}

// LoginHandler handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/v1/auth/login [post]
func (s *Server) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusBadRequest, "validation_error", "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Find user
	var user database.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
		return
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
		return
	}

	// Get entity ID based on user type
	entityID, err := s.getEntityID(user.ID, user.UserType)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "entity_lookup_error", "Failed to lookup entity", nil)
		return
	}

	// Generate JWT token
	tokenResponse, err := auth.GenerateToken(user.ID, user.Email, user.UserType, entityID)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "token_generation_error", "Failed to generate token", nil)
		return
	}

	middleware.StandardResponse(c, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"user_type": user.UserType,
			"entity_id": entityID,
		},
		"token": tokenResponse,
	}, "Login successful")
}

// RefreshTokenHandler handles JWT token refresh
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (s *Server) RefreshTokenHandler(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "no_claims", "No claims found in token", nil)
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		middleware.ErrorResponseWithDetails(c, http.StatusUnauthorized, "invalid_claims", "Invalid claims format", nil)
		return
	}

	// Generate new token
	tokenResponse, err := auth.RefreshToken(userClaims)
	if err != nil {
		middleware.ErrorResponseWithDetails(c, http.StatusInternalServerError, "token_generation_error", "Failed to refresh token", nil)
		return
	}

	middleware.StandardResponse(c, gin.H{
		"token": tokenResponse,
	}, "Token refreshed successfully")
}

// LogoutHandler handles user logout
// @Summary User logout
// @Description Logout user (client-side token invalidation)
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} middleware.SuccessResponse
// @Router /api/v1/auth/logout [post]
func (s *Server) LogoutHandler(c *gin.Context) {
	// In a JWT implementation, logout is typically handled client-side
	// by removing the token. For server-side token invalidation,
	// you would need to implement a token blacklist.
	
	middleware.StandardResponse(c, nil, "Logout successful")
}

// Helper functions

func (s *Server) createNGOEntity(userID uint, req RegisterRequest) (string, error) {
	// Generate NGO ID
	ngoID := generateEntityID("NGO")
	
	ngoModel := database.NGOModel{
		UserID:             userID,
		NGOID:              ngoID,
		Name:               req.Name,
		RegistrationNumber: req.RegistrationNumber,
		Category:           req.Category,
		Rating:             5.0,
		PublicKey:          generatePublicKey(ngoID),
	}

	if err := database.DB.Create(&ngoModel).Error; err != nil {
		return "", err
	}

	// Register in platform
	kycData := map[string]interface{}{
		"registration_number": req.RegistrationNumber,
		"category":           req.Category,
	}
	
	_, err := s.Platform.RegisterNGO(ngoID, req.Name, req.RegistrationNumber, req.Category, kycData, []string{})
	return ngoID, err
}

func (s *Server) createDonorEntity(userID uint, req RegisterRequest) (string, error) {
	// Generate Donor ID
	donorID := generateEntityID("DONOR")
	
	donorModel := database.DonorModel{
		UserID:              userID,
		DonorID:             donorID,
		AnnualDonationLimit: 1000000, // Default 10 lakh
	}

	if err := database.DB.Create(&donorModel).Error; err != nil {
		return "", err
	}

	// Register in platform
	kycData := map[string]interface{}{
		"annual_limit": 1000000,
	}
	
	_, err := s.Platform.RegisterDonor(donorID, kycData)
	return donorID, err
}

func (s *Server) createAuditorEntity(userID uint, req RegisterRequest) (string, error) {
	// Generate Auditor ID
	auditorID := generateEntityID("AUD")
	
	auditorModel := database.AuditorModel{
		UserID:          userID,
		AuditorID:       auditorID,
		Name:            req.Name,
		Rating:          5.0,
		PublicKey:       generatePublicKey(auditorID),
	}

	// Set specializations
	if len(req.Specializations) > 0 {
		specializationsJSON, _ := json.Marshal(req.Specializations)
		auditorModel.Specializations = string(specializationsJSON)
	}

	if err := database.DB.Create(&auditorModel).Error; err != nil {
		return "", err
	}

	// Register in platform
	credentials := map[string]interface{}{
		"name": req.Name,
	}
	
	_, err := s.Platform.RegisterAuditor(auditorID, req.Name, credentials, req.Specializations)
	return auditorID, err
}

func (s *Server) getEntityID(userID uint, userType string) (string, error) {
	switch userType {
	case "ngo":
		var ngo database.NGOModel
		if err := database.DB.Where("user_id = ?", userID).First(&ngo).Error; err != nil {
			return "", err
		}
		return ngo.NGOID, nil
	case "donor":
		var donor database.DonorModel
		if err := database.DB.Where("user_id = ?", userID).First(&donor).Error; err != nil {
			return "", err
		}
		return donor.DonorID, nil
	case "auditor":
		var auditor database.AuditorModel
		if err := database.DB.Where("user_id = ?", userID).First(&auditor).Error; err != nil {
			return "", err
		}
		return auditor.AuditorID, nil
	default:
		return "", nil
	}
}

// Helper utility functions
func generateEntityID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func generatePublicKey(entityID string) string {
	return fmt.Sprintf("pk_%s_%d", entityID, time.Now().UnixNano())
}