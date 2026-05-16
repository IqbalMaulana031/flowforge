package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

func Success(c *gin.Context, data any, meta *Meta) {
	body := gin.H{"success": true, "data": data}
	if meta != nil {
		body["meta"] = meta
	}
	c.JSON(http.StatusOK, body)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, statusCode int, code, message string, details ...ErrorDetail) {
	apiErr := APIError{Code: code, Message: message}
	if len(details) > 0 {
		apiErr.Details = details
	}

	c.JSON(statusCode, gin.H{
		"success":   false,
		"error":     apiErr,
		"requestId": c.GetString("requestId"),
	})
}
