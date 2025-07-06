package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload user profile and resume
func UploadProfileHandler(c *gin.Context) {
	// TODO: Implement file upload logic
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}

// Upload targets (CSV or JSON list)
func UploadTargetsHandler(c *gin.Context) {
	// TODO: Implement targets upload logic
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}

// Generate personalized email for a target
func GenerateEmailHandler(c *gin.Context) {
	// TODO: Call OpenAI, generate and return email
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}

// Send email to a target (attach resume)
func SendEmailHandler(c *gin.Context) {
	// TODO: Implement Gmail send logic
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}

// Get status of sent emails
func StatusHandler(c *gin.Context) {
	// TODO: Return email send status from DB
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}

// Get logs (basic)
func LogsHandler(c *gin.Context) {
	// TODO: Return recent operation logs
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "Not implemented"})
}
