package xgin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get("user_id")

	if !exists {
		return "", false
	}

	return value.(string), true
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, &RfcValidationError{
		Type:   "https://example.com/errors/internal",
		Title:  "Internal server error",
		Status: http.StatusInternalServerError,
		Detail: "An unexpected error occurred",
		Errors: []RfcFieldError{},
	})
}

func AbortMissingUserID(c *gin.Context) {
	log.Printf("[ERROR] %s: user_id not found in context", c.FullPath())
	InternalError(c)
}

func ParseUUID(c *gin.Context) (string, bool) {
	ID := c.Param("id")
	_, err := uuid.Parse(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, RfcValidationError{
			Type:   "https://example.com/errors/validation",
			Title:  "validation error",
			Status: http.StatusBadRequest,
			Detail: "Invalid UUID format",
			Errors: []RfcFieldError{},
		})
		return "", false
	}
	return ID, true
}

func ErrorResponse(c *gin.Context, status int, title string, detail string) {
	c.JSON(status, &RfcValidationError{
		Type:   "https://example.com/errors/" + strings.ReplaceAll(strings.ToLower(title), " ", "-"),
		Title:  title,
		Status: status,
		Detail: detail,
		Errors: []RfcFieldError{},
	})
}
