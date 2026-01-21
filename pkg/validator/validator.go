package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is the global validator instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(data interface{}) []ValidationError {
	var errors []ValidationError

	err := Validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			switch err.Tag() {
			case "required":
				message = err.Field() + " is required"
			case "email":
				message = err.Field() + " must be a valid email address"
			case "min":
				message = err.Field() + " must be at least " + err.Param() + " characters"
			case "max":
				message = err.Field() + " must be at most " + err.Param() + " characters"
			case "gte":
				message = err.Field() + " must be greater than or equal to " + err.Param()
			case "lte":
				message = err.Field() + " must be less than or equal to " + err.Param()
			case "oneof":
				message = err.Field() + " must be one of: " + err.Param()
			case "uuid":
				message = err.Field() + " must be a valid UUID"
			case "url":
				message = err.Field() + " must be a valid URL"
			default:
				message = err.Field() + " is invalid"
			}
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: message,
			})
		}
	}

	return errors
}

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if the password meets requirements
// Requirements: at least 8 characters
func IsValidPassword(password string) bool {
	return len(password) >= 8
}

// SanitizeString trims whitespace and removes dangerous characters
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")
	// Trim hyphens from start and end
	s = strings.Trim(s, "-")
	return s
}
