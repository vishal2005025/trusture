package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ngo-transparency-platform/pkg/config"
)

var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger(cfg *config.Config) {
	Logger = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Set log format
	if cfg.Logging.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Success   bool                   `json:"success"`
	Data      interface{}            `json:"data"`
	Message   string                 `json:"message,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination `json:"pagination"`
	RequestID  string      `json:"request_id,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// RequestLogger middleware logs all HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Generate request ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Get user info if available
		userID, _ := c.Get("user_id")
		userType, _ := c.Get("user_type")

		logFields := logrus.Fields{
			"request_id":   requestID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"status_code":  statusCode,
			"duration_ms":  duration.Milliseconds(),
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"query_params": c.Request.URL.RawQuery,
		}

		// Add user info if authenticated
		if userID != nil {
			logFields["user_id"] = userID
			logFields["user_type"] = userType
		}

		// Log with appropriate level based on status code
		switch {
		case statusCode >= 500:
			Logger.WithFields(logFields).Error("Server error")
		case statusCode >= 400:
			Logger.WithFields(logFields).Warn("Client error")
		case statusCode >= 300:
			Logger.WithFields(logFields).Info("Redirect")
		default:
			Logger.WithFields(logFields).Info("Request processed")
		}
	}
}

// ErrorHandler middleware handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				requestID, _ := c.Get("request_id")
				
				Logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      fmt.Sprintf("%v", err),
					"stack":      string(debug.Stack()),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				// Return error response
				errorResponse := ErrorResponse{
					Error:     "internal_server_error",
					Message:   "An unexpected error occurred",
					RequestID: fmt.Sprintf("%v", requestID),
					Timestamp: time.Now(),
				}

				// Don't expose panic details in production
				if config.IsDevelopment() {
					errorResponse.Details = map[string]interface{}{
						"panic": fmt.Sprintf("%v", err),
					}
				}

				c.JSON(http.StatusInternalServerError, errorResponse)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// CORS middleware configuration
func CORSMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	
	if config.IsDevelopment() {
		corsConfig.AllowOrigins = []string{"*"}
		corsConfig.AllowHeaders = []string{"*"}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	} else {
		// Production CORS settings
		corsConfig.AllowOrigins = []string{
			"https://yourdomain.com", 
			"https://app.yourdomain.com",
		}
		corsConfig.AllowHeaders = []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
		}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	return cors.New(corsConfig)
}

// RateLimit middleware (basic implementation)
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old requests
		if requests, exists := rl.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) <= rl.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[clientIP] = validRequests
		}

		// Check rate limit
		if len(rl.requests[clientIP]) >= rl.limit {
			Logger.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"limit":     rl.limit,
				"window":    rl.window.String(),
			}).Warn("Rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:     "rate_limit_exceeded",
				Message:   "Too many requests. Please try again later.",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Add current request
		rl.requests[clientIP] = append(rl.requests[clientIP], now)
		c.Next()
	}
}

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		if !config.IsDevelopment() {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		c.Next()
	}
}

// ContentType middleware ensures JSON content type for API endpoints
func JSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for GET requests and specific paths
		if c.Request.Method == "GET" || 
		   strings.HasPrefix(c.Request.URL.Path, "/docs") ||
		   strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:     "invalid_content_type",
				Message:   "Content-Type must be application/json",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Health check middleware
func HealthCheck() gin.HandlerFunc {
	startTime := time.Now()
	
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			uptime := time.Since(startTime)
			
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"uptime":    uptime.String(),
				"version":   "1.0.0",
			})
			return
		}
		
		c.Next()
	}
}

// Helper functions

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// StandardResponse creates a standardized success response
func StandardResponse(c *gin.Context, data interface{}, message string) {
	requestID, _ := c.Get("request_id")
	
	response := SuccessResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: fmt.Sprintf("%v", requestID),
		Timestamp: time.Now(),
	}
	
	c.JSON(http.StatusOK, response)
}

// ErrorResponseWithDetails creates a standardized error response with details
func ErrorResponseWithDetails(c *gin.Context, statusCode int, errorCode, message string, details map[string]interface{}) {
	requestID, _ := c.Get("request_id")
	
	response := ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Details:   details,
		RequestID: fmt.Sprintf("%v", requestID),
		Timestamp: time.Now(),
	}
	
	c.JSON(statusCode, response)
}

// PaginatedResponseData creates a paginated response
func PaginatedResponseData(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	requestID, _ := c.Get("request_id")
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	pagination := Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
	
	response := PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		RequestID:  fmt.Sprintf("%v", requestID),
		Timestamp:  time.Now(),
	}
	
	c.JSON(http.StatusOK, response)
}