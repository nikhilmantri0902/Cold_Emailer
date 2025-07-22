package main

import (
	"cold_emailer/api"
	"cold_emailer/constants"
	"cold_emailer/db"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// to add migration file:
	// go run cmd/server/main.go --generate-migration profile_info_resume_path
	migrationSuffix := flag.String("generate-migration", "", "Generate a new migration file with the given suffix and exit")
	flag.Parse()

	if *migrationSuffix != "" {
		path, err := db.GenerateMigrationFile(*migrationSuffix)
		if err != nil {
			log.Fatalf("Failed to generate migration: %v", err)
		}
		fmt.Printf("Created migration: %s\n", path)
		os.Exit(0)
	}

	constants.PrintENV()

	// Initialize DB (singleton)
	if err := db.InitDB(constants.PG_DB_URL); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.GetDB().Close()

	// Run migrations
	if err := db.RunMigrations(db.GetDB()); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize storage service
	api.InitStorage()

	// Initialize Gin
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Example health route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiRoutes := r.Group("/api")
	apiRoutes.POST("/profile", api.UploadProfileHandler)
	apiRoutes.POST("/generate-email", api.GenerateEmailHandler)
	apiRoutes.POST("/send-single-email", api.SendSingleEmailHandler)
	apiRoutes.POST("/enrich-database", api.EnrichDatabaseHandler)
	apiRoutes.POST("/backfill-company-details", api.BackfillCompanyDetails)
	apiRoutes.POST("/send-few-initial-emails", api.SendFewInitialEmailsHandler)
	apiRoutes.POST("/send-few-follow-up-emails", api.SendFewFollowUpEmailsHandler)
	// New endpoints for frontend
	apiRoutes.GET("/email-logs", api.GetEmailLogsHandler)
	apiRoutes.GET("/companies", api.GetCompaniesHandler)
	apiRoutes.GET("/companies/:company_id/contacts", api.GetCompanyContactsHandler)
	apiRoutes.GET("/config", api.GetSystemConfigHandler)

	// Gmail OAuth2 endpoints
	r.GET("/gmail-auth-initiate", api.GmailAuthInitiateHandler)
	r.GET("/gmail-oauth2callback", api.GmailOAuth2CallbackHandler)

	log.Printf("Server starting on port %s", constants.PORT)
	r.Run(":" + constants.PORT)
}
