package main

import (
	"fmt"
	"nutritional-score/pkg/models"
)

// This is a simple test file to demonstrate the enhanced scoring system
// Run with: go run test_scoring.go nutritionalscore.go
func testScoring() {
	// Test data for an apple (healthy food)
	apple := models.NutritionalData{
		Energy:              models.EnergyKJ(218),    // 52 kcal = ~218 kJ
		Sugars:              models.SugarGram(10.4),  // Natural fruit sugars
		SaturatedFattyAcids: models.SaturatedFattyAcids(0.1), // Very low
		Sodium:              models.SodiumMilligram(1),       // Very low
		Fruits:              models.FruitsPercent(100),       // 100% fruit
		Fibre:               models.FibreGram(2.4),           // Good fiber content
		Protein:             models.ProteinGram(0.3),         // Low protein
	}

	// Test data for a chocolate bar (unhealthy food)
	chocolate := models.NutritionalData{
		Energy:              models.EnergyKJ(2200),   // ~525 kcal
		Sugars:              models.SugarGram(47),    // High sugar
		SaturatedFattyAcids: models.SaturatedFattyAcids(18), // High saturated fat
		Sodium:              models.SodiumMilligram(24),     // Low sodium
		Fruits:              models.FruitsPercent(0),        // No fruits
		Fibre:               models.FibreGram(7),            // Some fiber
		Protein:             models.ProteinGram(8),          // Some protein
	}

	fmt.Println("=== Enhanced Nutritional Scoring Test ===")
	
	// Test apple scoring
	appleScore := GetNutritionalScore(apple, models.FoodType)
	fmt.Printf("\nApple (Healthy Food):\n")
	fmt.Printf("  Score: %d (Grade: %s)\n", appleScore.Value, appleScore.Grade)
	fmt.Printf("  Positive Points: %d\n", appleScore.Positive)
	fmt.Printf("  Negative Points: %d\n", appleScore.Negative)
	
	// Test chocolate scoring
	chocolateScore := GetNutritionalScore(chocolate, models.FoodType)
	fmt.Printf("\nChocolate Bar (Unhealthy Food):\n")
	fmt.Printf("  Score: %d (Grade: %s)\n", chocolateScore.Value, chocolateScore.Grade)
	fmt.Printf("  Positive Points: %d\n", chocolateScore.Positive)
	fmt.Printf("  Negative Points: %d\n", chocolateScore.Negative)
	
	// Test validation
	fmt.Printf("\n=== Validation Test ===\n")
	invalidData := models.NutritionalData{
		Energy:              models.EnergyKJ(-100),   // Invalid: negative energy
		Sugars:              models.SugarGram(150),   // Invalid: too high
		SaturatedFattyAcids: models.SaturatedFattyAcids(5), // Valid
		Sodium:              models.SodiumMilligram(500),   // Valid
		Fruits:              models.FruitsPercent(50),      // Valid
		Fibre:               models.FibreGram(3),           // Valid
		Protein:             models.ProteinGram(10),        // Valid
	}
	
	validationErrors := ValidateNutritionalData(invalidData)
	if len(validationErrors) > 0 {
		fmt.Printf("Validation errors found:\n")
		for _, err := range validationErrors {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Printf("No validation errors found.\n")
	}
	
	// Show grade thresholds
	fmt.Printf("\n=== Nutri-Score Grade Thresholds ===\n")
	thresholds := GetScoreThresholds()
	fmt.Printf("Grade A: Score ≤ %d (Best)\n", thresholds["A"])
	fmt.Printf("Grade B: Score ≤ %d\n", thresholds["B"])
	fmt.Printf("Grade C: Score ≤ %d\n", thresholds["C"])
	fmt.Printf("Grade D: Score ≤ %d\n", thresholds["D"])
	fmt.Printf("Grade E: Score ≥ %d (Worst)\n", thresholds["E"])
}

func main() {
	testScoring()
}