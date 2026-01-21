package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse is the response structure for paginated data
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 response for resource creation
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response with the given status code
func Error(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends a 401 error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, "unauthorized")
}

// Forbidden sends a 403 error response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, "forbidden")
}

// NotFound sends a 404 error response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "not_found")
}

// InternalError sends a 500 error response
func InternalError(c *gin.Context, message string, err string) {
	Error(c, http.StatusInternalServerError, message, err)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, items interface{}, total int64, page, limit int) {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Data retrieved successfully",
		Data: PaginatedResponse{
			Items:      items,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}
