package models

import (
	"fmt"
)

// ValidationError represents a validation error for nutritional data
// This struct provides detailed information about validation failures
type ValidationError struct {
	Field   string   `json:"field"`             // Name of the field that failed validation
	Value   float64  `json:"value"`             // The invalid value
	Message string   `json:"message"`           // Human-readable error message
	Min     *float64 `json:"min,omitempty"`     // Minimum allowed value (if applicable)
	Max     *float64 `json:"max,omitempty"`     // Maximum allowed value (if applicable)
}

// Error implements the error interface for ValidationError
func (ve ValidationError) Error() string {
	return ve.Message
}

// ErrorType represents the category of error that occurred
// This enum helps classify errors for appropriate handling and user messaging
type ErrorType string

const (
	ValidationErrorType  ErrorType = "validation"   // Input validation failures
	StorageErrorType     ErrorType = "storage"      // File system and data persistence errors
	DatabaseErrorType    ErrorType = "database"     // Food database access errors
	CalculationErrorType ErrorType = "calculation"  // Nutritional score calculation errors
	NetworkErrorType     ErrorType = "network"      // Network-related errors (future use)
	ConfigErrorType      ErrorType = "config"       // Configuration and settings errors
	UserInputErrorType   ErrorType = "user_input"   // User input processing errors
	ExportErrorType      ErrorType = "export"       // Data export operation errors
	ImportErrorType      ErrorType = "import"       // Data import operation errors (future use)
	SystemErrorType      ErrorType = "system"       // System-level errors
)

// NutritionalError represents a structured error with context and type information
// This struct provides detailed error information for better debugging and user experience
type NutritionalError struct {
	Type        ErrorType `json:"type"`                   // Category of the error
	Message     string    `json:"message"`                // Human-readable error message
	Field       string    `json:"field,omitempty"`        // Specific field that caused the error (if applicable)
	Code        string    `json:"code,omitempty"`         // Error code for programmatic handling
	Details     string    `json:"details,omitempty"`      // Additional technical details
	Suggestions []string  `json:"suggestions,omitempty"`  // Suggested actions to resolve the error
	Timestamp   string    `json:"timestamp,omitempty"`    // When the error occurred
}

// Error implements the error interface for NutritionalError
func (ne NutritionalError) Error() string {
	if ne.Field != "" {
		return fmt.Sprintf("%s error in field '%s': %s", ne.Type, ne.Field, ne.Message)
	}
	return fmt.Sprintf("%s error: %s", ne.Type, ne.Message)
}

// NewValidationError creates a new validation error with helpful context
func NewValidationError(field, message string, suggestions ...string) NutritionalError {
	return NutritionalError{
		Type:        ValidationErrorType,
		Message:     message,
		Field:       field,
		Code:        "VALIDATION_FAILED",
		Suggestions: suggestions,
	}
}

// NewStorageError creates a new storage-related error
func NewStorageError(message, details string) NutritionalError {
	return NutritionalError{
		Type:    StorageErrorType,
		Message: message,
		Code:    "STORAGE_FAILED",
		Details: details,
		Suggestions: []string{
			"Check file permissions",
			"Ensure sufficient disk space",
			"Verify data directory exists",
		},
	}
}

// NewDatabaseError creates a new database-related error
func NewDatabaseError(message, details string) NutritionalError {
	return NutritionalError{
		Type:    DatabaseErrorType,
		Message: message,
		Code:    "DATABASE_ERROR",
		Details: details,
		Suggestions: []string{
			"Check if food database is properly loaded",
			"Verify database file integrity",
			"Try restarting the application",
		},
	}
}

// NewCalculationError creates a new calculation-related error
func NewCalculationError(message, details string) NutritionalError {
	return NutritionalError{
		Type:    CalculationErrorType,
		Message: message,
		Code:    "CALCULATION_ERROR",
		Details: details,
		Suggestions: []string{
			"Verify all nutritional values are valid numbers",
			"Check that score type is appropriate for the food",
			"Ensure nutritional data is within acceptable ranges",
		},
	}
}

// NewUserInputError creates a new user input error
func NewUserInputError(message string, suggestions ...string) NutritionalError {
	return NutritionalError{
		Type:        UserInputErrorType,
		Message:     message,
		Code:        "INPUT_ERROR",
		Suggestions: suggestions,
	}
}

// NewExportError creates a new export-related error
func NewExportError(message, details string) NutritionalError {
	return NutritionalError{
		Type:    ExportErrorType,
		Message: message,
		Code:    "EXPORT_ERROR",
		Details: details,
		Suggestions: []string{
			"Check export directory permissions",
			"Ensure sufficient disk space",
			"Verify export format is supported",
		},
	}
}

// NewConfigError creates a new configuration-related error
func NewConfigError(message, details string) NutritionalError {
	return NutritionalError{
		Type:    ConfigErrorType,
		Message: message,
		Code:    "CONFIG_ERROR",
		Details: details,
		Suggestions: []string{
			"Check configuration file format",
			"Verify configuration file permissions",
			"Reset to default configuration if needed",
		},
	}
}

// ErrorCollection represents multiple errors that occurred during an operation
// This struct is useful for operations that can have multiple validation or processing errors
type ErrorCollection struct {
	Errors      []NutritionalError `json:"errors"`       // List of individual errors
	Operation   string             `json:"operation"`    // The operation that failed
	Summary     string             `json:"summary"`      // Brief summary of the errors
	ErrorCount  int                `json:"error_count"`  // Total number of errors
	WarningCount int               `json:"warning_count"` // Number of non-critical warnings
}

// Error implements the error interface for ErrorCollection
func (ec ErrorCollection) Error() string {
	if ec.ErrorCount == 1 {
		return fmt.Sprintf("1 error in %s: %s", ec.Operation, ec.Summary)
	}
	return fmt.Sprintf("%d errors in %s: %s", ec.ErrorCount, ec.Operation, ec.Summary)
}

// HasErrors returns true if the collection contains any errors
func (ec ErrorCollection) HasErrors() bool {
	return ec.ErrorCount > 0
}

// HasWarnings returns true if the collection contains any warnings
func (ec ErrorCollection) HasWarnings() bool {
	return ec.WarningCount > 0
}

// AddError adds a new error to the collection
func (ec *ErrorCollection) AddError(err NutritionalError) {
	ec.Errors = append(ec.Errors, err)
	ec.ErrorCount++
}

// GetErrorsByType returns all errors of a specific type
func (ec ErrorCollection) GetErrorsByType(errorType ErrorType) []NutritionalError {
	var filtered []NutritionalError
	for _, err := range ec.Errors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// GetErrorsByField returns all errors related to a specific field
func (ec ErrorCollection) GetErrorsByField(field string) []NutritionalError {
	var filtered []NutritionalError
	for _, err := range ec.Errors {
		if err.Field == field {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// Common error messages and codes for consistency across the application
var (
	// Validation error messages
	ErrInvalidEnergyValue = "Energy value must be between 0 and 4000 kJ per 100g"
	ErrInvalidSugarValue  = "Sugar value must be between 0 and 100g per 100g"
	ErrInvalidFatValue    = "Saturated fat value must be between 0 and 100g per 100g"
	ErrInvalidSodiumValue = "Sodium value must be between 0 and 10000mg per 100g"
	ErrInvalidFruitValue  = "Fruit/vegetable percentage must be between 0 and 100"
	ErrInvalidFibreValue  = "Fiber value must be between 0 and 50g per 100g"
	ErrInvalidProteinValue = "Protein value must be between 0 and 100g per 100g"
	
	// Food-related error messages
	ErrFoodNotFound      = "Food not found in database"
	ErrInvalidFoodID     = "Invalid food ID format"
	ErrFoodNameRequired  = "Food name is required"
	ErrFoodNameTooLong   = "Food name must be less than 200 characters"
	ErrInvalidCategory   = "Invalid food category"
	
	// Storage error messages
	ErrStorageNotInitialized = "Storage system not initialized"
	ErrFileNotFound         = "Required data file not found"
	ErrPermissionDenied     = "Permission denied accessing data file"
	ErrDiskSpaceFull        = "Insufficient disk space for operation"
	ErrCorruptedData        = "Data file appears to be corrupted"
	
	// Calculation error messages
	ErrInvalidScoreType     = "Invalid score type for calculation"
	ErrCalculationFailed    = "Failed to calculate nutritional score"
	ErrIncompleteData       = "Incomplete nutritional data for calculation"
	
	// Export error messages
	ErrUnsupportedFormat    = "Unsupported export format"
	ErrExportFailed         = "Failed to export data"
	ErrInvalidExportPath    = "Invalid export file path"
	
	// User input error messages
	ErrInvalidMenuChoice    = "Invalid menu choice selected"
	ErrEmptySearchQuery     = "Search query cannot be empty"
	ErrSearchQueryTooShort  = "Search query must be at least 2 characters"
	ErrInvalidNumericInput  = "Invalid numeric input"
)

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if ne, ok := err.(NutritionalError); ok {
		return ne.Type == ValidationErrorType
	}
	return false
}

// IsStorageError checks if an error is a storage error
func IsStorageError(err error) bool {
	if ne, ok := err.(NutritionalError); ok {
		return ne.Type == StorageErrorType
	}
	return false
}

// IsDatabaseError checks if an error is a database error
func IsDatabaseError(err error) bool {
	if ne, ok := err.(NutritionalError); ok {
		return ne.Type == DatabaseErrorType
	}
	return false
}

// WrapError wraps a standard error into a NutritionalError with additional context
func WrapError(err error, errorType ErrorType, message string) NutritionalError {
	return NutritionalError{
		Type:    errorType,
		Message: message,
		Details: err.Error(),
	}
}