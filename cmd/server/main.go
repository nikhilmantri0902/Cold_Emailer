package main

import (
	"cold_emailer/api"
	"cold_emailer/constants"
	"cold_emailer/db"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() func(context.Context) error {

	var secureOption otlptracegrpc.Option

	if strings.ToLower(constants.INSECURE_MODE) == "false" || constants.INSECURE_MODE == "0" || strings.ToLower(constants.INSECURE_MODE) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(constants.COLLECTOR_URL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", constants.SERVICE_NAME),
			attribute.String("library.language", "go"),
			attribute.String("deployment.environment", constants.ENV),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {

	cleanup := initTracer()
	defer cleanup(context.Background())

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
	log.Println("SERVICE_NAME: ", constants.SERVICE_NAME)
	log.Println("OTEL_EXPORTER_OTLP_ENDPOINT: ", constants.COLLECTOR_URL)
	log.Println("INSECURE_MODE: ", constants.INSECURE_MODE)

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
	r.Use(otelgin.Middleware(constants.SERVICE_NAME))

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
	apiRoutes.POST("/generate-late-spans", api.GenerateLateSpansHandler)

	// Frontend API endpoints
	frontendRoutes := r.Group("/frontend")
	frontendRoutes.GET("/email-logs", api.GetEmailLogsHandler)
	frontendRoutes.GET("/companies", api.GetCompaniesHandler)
	frontendRoutes.GET("/companies/:company_id/contacts", api.GetCompanyContactsHandler)
	frontendRoutes.GET("/config", api.GetSystemConfigHandler)
	frontendRoutes.GET("/job-status", api.GetJobStatusHandler)
	frontendRoutes.POST("/enrich-database", api.EnrichDatabaseHandler)

	// Gmail OAuth2 endpoints
	r.GET("/gmail-auth-initiate", api.GmailAuthInitiateHandler)
	r.GET("/gmail-oauth2callback", api.GmailOAuth2CallbackHandler)

	log.Printf("Server starting on port %s", constants.PORT)
	r.Run(":" + constants.PORT)
}
