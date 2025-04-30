package snippets

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleErrorJSONAnswer(c *gin.Context, logger *zap.Logger,
	statusCode int, errorMessage, errorLogMessage, logSource string) {
	log.Printf("[ERROR] %s %s\n", logSource, errorLogMessage)
	logger.Error(errorLogMessage)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errorMessage,
	})
}

func HandleError(logger *zap.Logger, errorLogMessage, logSource string) {
	log.Printf("[ERROR] %s %s\n", logSource, errorLogMessage)
	logger.Error(errorLogMessage)
}

func HandleInfoLogs(logger *zap.Logger, logMessage, logSource string) {
	log.Printf("[INFO] %s %s\n", logSource, logMessage)
	logger.Info(logMessage)
}

func HandleDebugLogs(logger *zap.Logger, logMessage string) {
	log.Printf("[DEBUG] %s\n", logMessage)
	logger.Debug(logMessage)
}
