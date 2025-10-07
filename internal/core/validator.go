package core

import (
	"fmt"
	"nutritional-score/pkg/models"
	"strings"
)

// InputValidator implements validation logic for nutritional data and user inputs
// This struct ensures data integrity and provides helpful error messages
type InputValidator struct {
	validationRules models.NutritionalDataValidation
}

// NewInputValidator creates a new input validator with default validation rules
func NewInputValidator() *InputValidator {
	return &InputValidator{
		validationRules: models.DefaultValidationRules(),
	}
}

// NewInputValidatorWithRules creates a validator with custom validation rules
func NewInputValidatorWithRules(rules models.NutritionalDataValidation) *InputValidator {
	return &InputValidator{
		validationRules: rules,
	}
}

// ValidateNutritionalData validates all nutritional data fields against defined rules
// Returns a slice of validation errors for any invalid values
func (iv *InputValidator) ValidateNutritionalData(data models.NutritionalData) []models.ValidationError {
	var errors []models.ValidationError

	// Validate Energy (kJ per 100g)
	energy := float64(data.Energy)
	if energy < iv.validationRules.EnergyMin {
		errors = append(errors, models.ValidationError{
			Field:   "energy",
			Value:   energy,
			Message: fmt.Sprintf("Energy cannot be less than %.1f kJ per 100g", iv.validationRules.EnergyMin),
			Min:     &iv.validationRules.EnergyMin,
			Max:     &iv.validationRules.EnergyMax,
		})
	}
	if energy > iv.validationRules.EnergyMax {
		errors = append(errors, models.ValidationError{
			Field:   "energy",
			Value:   energy,
			Message: fmt.Sprintf("Energy cannot exceed %.1f kJ per 100g", iv.validationRules.EnergyMax),
			Min:     &iv.validationRules.EnergyMin,
			Max:     &iv.validationRules.EnergyMax,
		})
	}

	// Validate Sugars (g per 100g)
	sugars := float64(data.Sugars)
	if sugars < iv.validationRules.SugarsMin {
		errors = append(errors, models.ValidationError{
			Field:   "sugars",
			Value:   sugars,
			Message: fmt.Sprintf("Sugar content cannot be less than %.1f g per 100g", iv.validationRules.SugarsMin),
			Min:     &iv.validationRules.SugarsMin,
			Max:     &iv.validationRules.SugarsMax,
		})
	}
	if sugars > iv.validationRules.SugarsMax {
		errors = append(errors, models.ValidationError{
			Field:   "sugars",
			Value:   sugars,
			Message: fmt.Sprintf("Sugar content cannot exceed %.1f g per 100g", iv.validationRules.SugarsMax),
			Min:     &iv.validationRules.SugarsMin,
			Max:     &iv.validationRules.SugarsMax,
		})
	}

	// Validate Saturated Fatty Acids (g per 100g)
	satFat := float64(data.SaturatedFattyAcids)
	if satFat < iv.validationRules.SaturatedFatMin {
		errors = append(errors, models.ValidationError{
			Field:   "saturated_fatty_acids",
			Value:   satFat,
			Message: fmt.Sprintf("Saturated fat content cannot be less than %.1f g per 100g", iv.validationRules.SaturatedFatMin),
			Min:     &iv.validationRules.SaturatedFatMin,
			Max:     &iv.validationRules.SaturatedFatMax,
		})
	}
	if satFat > iv.validationRules.SaturatedFatMax {
		errors = append(errors, models.ValidationError{
			Field:   "saturated_fatty_acids",
			Value:   satFat,
			Message: fmt.Sprintf("Saturated fat content cannot exceed %.1f g per 100g", iv.validationRules.SaturatedFatMax),
			Min:     &iv.validationRules.SaturatedFatMin,
			Max:     &iv.validationRules.SaturatedFatMax,
		})
	}

	// Validate Sodium (mg per 100g)
	sodium := float64(data.Sodium)
	if sodium < iv.validationRules.SodiumMin {
		errors = append(errors, models.ValidationError{
			Field:   "sodium",
			Value:   sodium,
			Message: fmt.Sprintf("Sodium content cannot be less than %.1f mg per 100g", iv.validationRules.SodiumMin),
			Min:     &iv.validationRules.SodiumMin,
			Max:     &iv.validationRules.SodiumMax,
		})
	}
	if sodium > iv.validationRules.SodiumMax {
		errors = append(errors, models.ValidationError{
			Field:   "sodium",
			Value:   sodium,
			Message: fmt.Sprintf("Sodium content cannot exceed %.1f mg per 100g", iv.validationRules.SodiumMax),
			Min:     &iv.validationRules.SodiumMin,
			Max:     &iv.validationRules.SodiumMax,
		})
	}

	// Validate Fruits/Vegetables/Nuts percentage
	fruits := float64(data.Fruits)
	if fruits < iv.validationRules.FruitsMin {
		errors = append(errors, models.ValidationError{
			Field:   "fruits",
			Value:   fruits,
			Message: fmt.Sprintf("Fruits/vegetables/nuts percentage cannot be less than %.1f%%", iv.validationRules.FruitsMin),
			Min:     &iv.validationRules.FruitsMin,
			Max:     &iv.validationRules.FruitsMax,
		})
	}
	if fruits > iv.validationRules.FruitsMax {
		errors = append(errors, models.ValidationError{
			Field:   "fruits",
			Value:   fruits,
			Message: fmt.Sprintf("Fruits/vegetables/nuts percentage cannot exceed %.1f%%", iv.validationRules.FruitsMax),
			Min:     &iv.validationRules.FruitsMin,
			Max:     &iv.validationRules.FruitsMax,
		})
	}

	// Validate Fiber (g per 100g)
	fiber := float64(data.Fibre)
	if fiber < iv.validationRules.FibreMin {
		errors = append(errors, models.ValidationError{
			Field:   "fibre",
			Value:   fiber,
			Message: fmt.Sprintf("Fiber content cannot be less than %.1f g per 100g", iv.validationRules.FibreMin),
			Min:     &iv.validationRules.FibreMin,
			Max:     &iv.validationRules.FibreMax,
		})
	}
	if fiber > iv.validationRules.FibreMax {
		errors = append(errors, models.ValidationError{
			Field:   "fibre",
			Value:   fiber,
			Message: fmt.Sprintf("Fiber content cannot exceed %.1f g per 100g", iv.validationRules.FibreMax),
			Min:     &iv.validationRules.FibreMin,
			Max:     &iv.validationRules.FibreMax,
		})
	}

	// Validate Protein (g per 100g)
	protein := float64(data.Protein)
	if protein < iv.validationRules.ProteinMin {
		errors = append(errors, models.ValidationError{
			Field:   "protein",
			Value:   protein,
			Message: fmt.Sprintf("Protein content cannot be less than %.1f g per 100g", iv.validationRules.ProteinMin),
			Min:     &iv.validationRules.ProteinMin,
			Max:     &iv.validationRules.ProteinMax,
		})
	}
	if protein > iv.validationRules.ProteinMax {
		errors = append(errors, models.ValidationError{
			Field:   "protein",
			Value:   protein,
			Message: fmt.Sprintf("Protein content cannot exceed %.1f g per 100g", iv.validationRules.ProteinMax),
			Min:     &iv.validationRules.ProteinMin,
			Max:     &iv.validationRules.ProteinMax,
		})
	}

	return errors
}

// ValidateFood validates a complete food item including name, category, and nutritional data
func (iv *InputValidator) ValidateFood(food models.Food) []models.ValidationError {
	var errors []models.ValidationError

	// Validate food name
	if strings.TrimSpace(food.Name) == "" {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Value:   0, // Not applicable for string fields
			Message: "Food name is required and cannot be empty",
		})
	}

	if len(food.Name) > 200 {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Value:   float64(len(food.Name)),
			Message: "Food name must be less than 200 characters",
			Max:     func() *float64 { v := 200.0; return &v }(),
		})
	}

	// Validate food category
	if strings.TrimSpace(food.Category) == "" {
		errors = append(errors, models.ValidationError{
			Field:   "category",
			Value:   0,
			Message: "Food category is required",
		})
	}

	// Validate food ID format (if provided)
	if food.ID != "" && !iv.isValidFoodID(food.ID) {
		errors = append(errors, models.ValidationError{
			Field:   "id",
			Value:   0,
			Message: "Food ID must contain only alphanumeric characters, hyphens, and underscores",
		})
	}

	// Validate nutritional data
	nutritionalErrors := iv.ValidateNutritionalData(food.NutritionalData)
	errors = append(errors, nutritionalErrors...)

	return errors
}

// ValidateScoreType checks if the provided score type is valid
func (iv *InputValidator) ValidateScoreType(scoreType models.ScoreType) error {
	switch scoreType {
	case models.FoodType, models.BeverageType, models.WaterType, models.CheeseType:
		return nil
	default:
		return models.NewValidationError("score_type", 
			fmt.Sprintf("Invalid score type: %d. Must be 0 (Food), 1 (Beverage), 2 (Water), or 3 (Cheese)", int(scoreType)),
			"Use 0 for Food, 1 for Beverage, 2 for Water, or 3 for Cheese")
	}
}

// ValidateSearchQuery checks if a search query is valid and useful
func (iv *InputValidator) ValidateSearchQuery(query string) error {
	trimmed := strings.TrimSpace(query)
	
	if trimmed == "" {
		return models.NewUserInputError("Search query cannot be empty", 
			"Enter at least 2 characters to search for foods")
	}
	
	if len(trimmed) < 2 {
		return models.NewUserInputError("Search query is too short", 
			"Enter at least 2 characters to get meaningful search results")
	}
	
	if len(trimmed) > 100 {
		return models.NewUserInputError("Search query is too long", 
			"Search query must be less than 100 characters")
	}
	
	return nil
}

// ValidateExportFormat checks if the export format is supported
func (iv *InputValidator) ValidateExportFormat(format models.ExportFormat) error {
	switch format {
	case models.JSON, models.CSV, models.XML:
		return nil
	default:
		return models.NewValidationError("export_format", 
			fmt.Sprintf("Unsupported export format: %d", int(format)),
			"Use 0 for JSON, 1 for CSV, or 2 for XML")
	}
}

// isValidFoodID checks if a food ID contains only allowed characters
// Food IDs should contain only alphanumeric characters, hyphens, and underscores
func (iv *InputValidator) isValidFoodID(id string) bool {
	if len(id) == 0 || len(id) > 50 {
		return false
	}
	
	for _, char := range id {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return false
		}
	}
	
	return true
}

// ValidateNumericInput validates that a string can be converted to a valid number
// Used for CLI input validation
func (iv *InputValidator) ValidateNumericInput(input string, fieldName string, min, max float64) error {
	trimmed := strings.TrimSpace(input)
	
	if trimmed == "" {
		return models.NewUserInputError(
			fmt.Sprintf("%s cannot be empty", fieldName),
			fmt.Sprintf("Enter a number between %.1f and %.1f", min, max))
	}
	
	// This is a basic validation - actual numeric conversion would be done elsewhere
	// Here we just check for obviously invalid formats
	if strings.Contains(trimmed, " ") {
		return models.NewUserInputError(
			fmt.Sprintf("Invalid %s format", fieldName),
			"Enter a valid number without spaces")
	}
	
	return nil
}

// GetValidationRules returns the current validation rules
func (iv *InputValidator) GetValidationRules() models.NutritionalDataValidation {
	return iv.validationRules
}

// SetValidationRules updates the validation rules
func (iv *InputValidator) SetValidationRules(rules models.NutritionalDataValidation) {
	iv.validationRules = rules
}

// ValidateNutritionalRange checks if a single nutritional value is within range
// This is a helper function for more granular validation
func (iv *InputValidator) ValidateNutritionalRange(value float64, min, max float64, fieldName string) *models.ValidationError {
	if value < min {
		return &models.ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("%s cannot be less than %.1f", fieldName, min),
			Min:     &min,
			Max:     &max,
		}
	}
	
	if value > max {
		return &models.ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("%s cannot exceed %.1f", fieldName, max),
			Min:     &min,
			Max:     &max,
		}
	}
	
	return nil
}