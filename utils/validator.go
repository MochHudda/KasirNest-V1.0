package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// ValidatePrice validates that price is a positive number
func ValidatePrice(price float64) bool {
	return price > 0
}

// ValidateStock validates that stock is not negative
func ValidateStock(stock int) bool {
	return stock >= 0
}

// ValidateQuantity validates that quantity is positive
func ValidateQuantity(quantity int) bool {
	return quantity > 0
}

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// ValidateMinLength checks if a string has minimum length
func ValidateMinLength(value string, minLength int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) >= minLength
}

// ValidateMaxLength checks if a string doesn't exceed maximum length
func ValidateMaxLength(value string, maxLength int) bool {
	return utf8.RuneCountInString(value) <= maxLength
}

// ValidatePasswordStrength checks password strength
func ValidatePasswordStrength(password string) bool {
	if len(password) < 6 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// ValidateBarcode validates barcode format (basic check)
func ValidateBarcode(barcode string) bool {
	if barcode == "" {
		return true // Barcode is optional
	}

	// Check if barcode contains only numbers and has reasonable length
	matched, _ := regexp.MatchString(`^\d{8,13}$`, barcode)
	return matched
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationResult holds validation results
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (vr *ValidationResult) AddError(field, message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// GetErrorMessage returns formatted error message
func (vr *ValidationResult) GetErrorMessage() string {
	if vr.IsValid {
		return ""
	}

	var messages []string
	for _, err := range vr.Errors {
		messages = append(messages, err.Field+": "+err.Message)
	}

	return strings.Join(messages, "\n")
}
