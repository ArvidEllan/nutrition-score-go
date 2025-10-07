package main

import (
	"nutritional-score/internal/core"
	"nutritional-score/pkg/models"
)

// Legacy type aliases for backward compatibility with existing main.go
// These will be replaced when we update the main application structure
type ScoreType = models.ScoreType
type NutritionalScore = models.NutritionalScore
type NutritionalData = models.NutritionalData

// Legacy constants for backward compatibility
const (
	Food     = models.FoodType
	Beverage = models.BeverageType
	Water    = models.WaterType
	Cheese   = models.CheeseType
)

// Legacy type aliases for nutritional components
// These maintain backward compatibility while using the new models
type EnergyKJ = models.EnergyKJ
type SugarGram = models.SugarGram
type SaturatedFattyAcids = models.SaturatedFattyAcids
type SodiumMilligram = models.SodiumMilligram
type FruitsPercent = models.FruitsPercent
type FibreGram = models.FibreGram
type ProteinGram = models.ProteinGram

// GetNutritionalScore calculates the nutritional score using the enhanced scoring engine
// This function now uses the official Nutri-Score algorithm with proper validation
func GetNutritionalScore(n NutritionalData, st ScoreType) NutritionalScore {
	// Create the enhanced nutritional scorer with official algorithm
	scorer := core.NewNutritionalScorer()
	
	// Calculate the score using the official Nutri-Score algorithm
	// This includes proper validation and grade assignment
	result, err := scorer.CalculateScore(n, st)
	if err != nil {
		// If there's a validation error, return a default score with error indication
		// In a real application, this error should be handled properly
		return NutritionalScore{
			Value:     999, // High value indicates error
			Grade:     "E", // Worst grade for invalid data
			Positive:  0,
			Negative:  0,
			ScoreType: st,
		}
	}
	
	return result
}
// 
ValidateNutritionalData validates nutritional data and returns user-friendly error messages
// This function provides a simple interface for validation in the CLI
func ValidateNutritionalData(n NutritionalData) []string {
	validator := core.NewInputValidator()
	validationErrors := validator.ValidateNutritionalData(n)
	
	// Convert validation errors to simple string messages for CLI display
	var messages []string
	for _, err := range validationErrors {
		messages = append(messages, err.Message)
	}
	
	return messages
}

// GetScoreGrade converts a numerical score to a letter grade using official thresholds
// This provides a simple interface to the enhanced scoring system
func GetScoreGrade(score int) string {
	scorer := core.NewNutritionalScorer()
	return scorer.GetScoreGrade(score)
}

// GetScoreThresholds returns the official Nutri-Score grade thresholds
// Useful for displaying grade information to users
func GetScoreThresholds() map[string]int {
	scorer := core.NewNutritionalScorer()
	return scorer.GetScoreThresholds()
}