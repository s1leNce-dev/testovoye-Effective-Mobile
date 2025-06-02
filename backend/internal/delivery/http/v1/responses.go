package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

// createPersonResponse описывает успешный ответ при создании Person
// swagger:model createPersonResponse
type createPersonResponse struct {
	// UUID нового Person
	// example: 3fa85f64-5717-4562-b3fc-2c963f66afa6
	ID string `json:"id"`
}

// simpleMessageResponse описывает ответ вида {"message": "success"}
// swagger:model simpleMessageResponse
type simpleMessageResponse struct {
	// Сообщение
	// example: success
	Message string `json:"message"`
}

// ErrorResponse описывает формат ответа с ошибкой
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Описание ошибки
	// example: invalid input body
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message interface{}, logMessage interface{}) {
	log.Printf("[ERROR] %s\n", logMessage)
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error": message,
	})
}
