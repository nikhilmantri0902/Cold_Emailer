package main

import (
	"cold_emailer/api"
	"cold_emailer/constants"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	constants.PrintENV()

	// Initialize Gin
	r := gin.Default()

	// Example health route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiRoutes := r.Group("/api")
	apiRoutes.POST("/profile", api.UploadProfileHandler)
	apiRoutes.POST("/targets", api.UploadTargetsHandler)
	apiRoutes.POST("/generate-email", api.GenerateEmailHandler)
	apiRoutes.POST("/send-email", api.SendEmailHandler)
	apiRoutes.GET("/status", api.StatusHandler)
	apiRoutes.GET("/logs", api.LogsHandler)

	log.Printf("Server starting on port %s", constants.PORT)
	r.Run(":" + constants.PORT)
}
