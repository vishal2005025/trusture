package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "ngo-transparency-platform/docs" // Import generated docs
	"ngo-transparency-platform/pkg/config"
	"ngo-transparency-platform/pkg/server"
)

// @title Trusture API
// @version 1.0
// @description Blockchain-based NGO donation auditing framework API
// @termsOfService https://trusture.io/terms
// @contact.name API Support
// @contact.url https://trusture.io/support
// @contact.email support@trusture.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create server
	srv := server.NewServer(cfg)

	// Initialize server
	if err := srv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		log.Println("Received shutdown signal, shutting down gracefully...")
		
		if err := srv.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		
		os.Exit(0)
	}()

	// Start server
	log.Printf("Starting Trusture API server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}