package xgin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RfcValidationError struct {
	Type   string          `json:"type"`
	Title  string          `json:"title"`
	Status int             `json:"status"`
	Detail string          `json:"detail"`
	Errors []RfcFieldError `json:"errors"`
}

type RfcFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validationMessages = map[string]string{
	"required":             "field is required",
	"email":                "invalid email address",
	"gt":                   "must be greater than 0",
	"required_without_all": "at least one field is required",
}

func valMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", fe.Param())
	default:
		if msg, ok := validationMessages[fe.Tag()]; ok {
			return msg
		}
		return fe.Tag()
	}
}

func BindError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError

	switch {
	case errors.As(err, &ve):
		jsonErrors := make([]RfcFieldError, len(ve))

		for idx := range ve {
			jsonErrors[idx] = RfcFieldError{
				Field:   ve[idx].Field(),
				Message: valMessage(ve[idx]),
			}
		}
		jsonErr := &RfcValidationError{
			Type:   "https://example.com/errors/validation",
			Title:  "Validation error",
			Status: http.StatusUnprocessableEntity,
			Detail: "One or more fields are invalid",
			Errors: jsonErrors,
		}
		c.JSON(http.StatusUnprocessableEntity, jsonErr)

	case errors.As(err, &syntaxErr):
		jsonErr := &RfcValidationError{
			Type:   "https://example.com/errors/bad-request",
			Title:  "Syntax error",
			Status: http.StatusBadRequest,
			Detail: "Malformed JSON",
			Errors: []RfcFieldError{},
		}
		c.JSON(http.StatusBadRequest, jsonErr)

	case errors.As(err, &unmarshalErr):
		jsonErr := &RfcValidationError{
			Type:   "https://example.com/errors/bad-request",
			Title:  "Wrong field type",
			Status: http.StatusBadRequest,
			Detail: "Invalid request body",
			Errors: []RfcFieldError{
				{Field: unmarshalErr.Field, Message: "invalid type"}}}
		c.JSON(http.StatusBadRequest, jsonErr)

	default:
		c.JSON(http.StatusBadRequest, &RfcValidationError{
			Type:   "https://example.com/errors/bad-request",
			Title:  "Bad request",
			Status: http.StatusBadRequest,
			Detail: "Invalid request body",
			Errors: []RfcFieldError{},
		})
	}
}
