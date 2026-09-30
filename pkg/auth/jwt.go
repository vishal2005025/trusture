package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"ngo-transparency-platform/pkg/config"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"` // "ngo", "donor", "auditor"
	EntityID string `json:"entity_id"` // NGO ID, Donor ID, or Auditor ID
	jwt.RegisteredClaims
}

// TokenResponse represents the response for token generation
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserType  string    `json:"user_type"`
	EntityID  string    `json:"entity_id"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserType string `json:"user_type" binding:"required,oneof=ngo donor auditor"`
}

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrTokenExpired   = errors.New("token expired")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidClaims  = errors.New("invalid claims")
)

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uint, email, userType, entityID string) (*TokenResponse, error) {
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpiryHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		UserType: userType,
		EntityID: entityID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "trusture-platform",
			Subject:   fmt.Sprintf("user_%d", userID),
			ID:        fmt.Sprintf("%d_%d", userID, time.Now().Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &TokenResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
		UserType:  userType,
		EntityID:  entityID,
	}, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if a password matches its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer token format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(bearerToken[1])
		if err != nil {
			var status int
			var message string
			
			switch err {
			case ErrTokenExpired:
				status = http.StatusUnauthorized
				message = "Token expired"
			case ErrInvalidToken, ErrInvalidClaims:
				status = http.StatusUnauthorized
				message = "Invalid token"
			default:
				status = http.StatusUnauthorized
				message = "Unauthorized"
			}

			c.JSON(status, gin.H{"error": message})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("entity_id", claims.EntityID)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireUserType middleware ensures the user has a specific user type
func RequireUserType(userTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User type not found in context"})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type format"})
			c.Abort()
			return
		}

		// Check if user type is in allowed types
		allowed := false
		for _, allowedType := range userTypes {
			if userTypeStr == allowedType {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Access denied. Required user type: %v, got: %s", userTypes, userTypeStr),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireEntityAccess middleware ensures the user can access specific entity data
func RequireEntityAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		entityID, exists := c.Get("entity_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Entity ID not found in context"})
			c.Abort()
			return
		}

		// Get entity ID from URL parameter
		urlEntityID := c.Param("id")
		if urlEntityID == "" {
			// For endpoints that don't have ID in URL, skip this check
			c.Next()
			return
		}

		if entityID != urlEntityID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this entity"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware validates JWT tokens if present, but doesn't require them
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check Bearer token format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := ValidateToken(bearerToken[1])
		if err != nil {
			// Token is invalid but we don't require auth, so continue
			c.Next()
			return
		}

		// Add claims to context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("entity_id", claims.EntityID)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserFromContext extracts user information from gin context
func GetUserFromContext(c *gin.Context) (uint, string, string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, "", "", ErrUnauthorized
	}

	userType, exists := c.Get("user_type")
	if !exists {
		return 0, "", "", ErrUnauthorized
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		return 0, "", "", ErrUnauthorized
	}

	return userID.(uint), userType.(string), entityID.(string), nil
}

// RefreshToken generates a new token with updated expiration
func RefreshToken(claims *Claims) (*TokenResponse, error) {
	// Create new token with same user info but updated expiration
	return GenerateToken(claims.UserID, claims.Email, claims.UserType, claims.EntityID)
}