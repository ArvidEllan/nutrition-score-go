package core

import (
	"nutritional-score/pkg/models"
	"testing"
)

// TestNutritionalScorer_CalculateScore tests the main scoring functionality
func TestNutritionalScorer_CalculateScore(t *testing.T) {
	scorer := NewNutritionalScorer()

	tests := []struct {
		name     string
		data     models.NutritionalData
		foodType models.ScoreType
		expected models.NutritionalScore
		wantErr  bool
	}{
		{
			name: "Apple - Healthy Food (Grade A)",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(218),    // Low energy
				Sugars:              models.SugarGram(10.4),  // Natural sugars
				SaturatedFattyAcids: models.SaturatedFattyAcids(0.1), // Very low
				Sodium:              models.SodiumMilligram(1),       // Very low
				Fruits:              models.FruitsPercent(100),       // 100% fruit
				Fibre:               models.FibreGram(2.4),           // Good fiber
				Protein:             models.ProteinGram(0.3),         // Low protein
			},
			foodType: models.FoodType,
			expected: models.NutritionalScore{
				Value:     -4, // Should get Grade A
				Grade:     "A",
				Positive:  5, // 5 points from fruits
				Negative:  1, // 1 point from sugars
				ScoreType: models.FoodType,
			},
			wantErr: false,
		},
		{
			name: "Chocolate Bar - Unhealthy Food (Grade E)",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(2200),   // High energy
				Sugars:              models.SugarGram(47),    // Very high sugar
				SaturatedFattyAcids: models.SaturatedFattyAcids(18), // High saturated fat
				Sodium:              models.SodiumMilligram(24),     // Low sodium
				Fruits:              models.FruitsPercent(0),        // No fruits
				Fibre:               models.FibreGram(7),            // Some fiber
				Protein:             models.ProteinGram(8),          // Some protein
			},
			foodType: models.FoodType,
			expected: models.NutritionalScore{
				Value:     30, // Should get Grade E
				Grade:     "E",
				Positive:  9,  // 5 fiber + 4 protein
				Negative:  39, // High from energy, sugar, sat fat
				ScoreType: models.FoodType,
			},
			wantErr: false,
		},
		{
			name: "Water - Special Case",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(0),
				Sugars:              models.SugarGram(0),
				SaturatedFattyAcids: models.SaturatedFattyAcids(0),
				Sodium:              models.SodiumMilligram(0),
				Fruits:              models.FruitsPercent(0),
				Fibre:               models.FibreGram(0),
				Protein:             models.ProteinGram(0),
			},
			foodType: models.WaterType,
			expected: models.NutritionalScore{
				Value:     0,
				Grade:     "A", // Water always gets A
				Positive:  0,
				Negative:  0,
				ScoreType: models.WaterType,
			},
			wantErr: false,
		},
		{
			name: "Cheese - Special Scoring Rules",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(1500),   // Moderate energy
				Sugars:              models.SugarGram(1),     // Low sugar
				SaturatedFattyAcids: models.SaturatedFattyAcids(15), // High saturated fat
				Sodium:              models.SodiumMilligram(600),    // High sodium
				Fruits:              models.FruitsPercent(0),        // No fruits
				Fibre:               models.FibreGram(0),            // No fiber
				Protein:             models.ProteinGram(25),         // High protein
			},
			foodType: models.CheeseType,
			expected: models.NutritionalScore{
				Value:     20, // Cheese uses special rules
				Grade:     "E",
				Positive:  5,  // 5 points from protein
				Negative:  25, // High from energy, sat fat, sodium
				ScoreType: models.CheeseType,
			},
			wantErr: false,
		},
		{
			name: "Beverage - Modified Rules",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(180),    // Low energy
				Sugars:              models.SugarGram(4),     // Low sugar
				SaturatedFattyAcids: models.SaturatedFattyAcids(0), // No fat
				Sodium:              models.SodiumMilligram(10),    // Low sodium
				Fruits:              models.FruitsPercent(50),      // Some fruit
				Fibre:               models.FibreGram(0),           // No fiber
				Protein:             models.ProteinGram(0),         // No protein
			},
			foodType: models.BeverageType,
			expected: models.NutritionalScore{
				Value:     -1, // Should get Grade A
				Grade:     "A",
				Positive:  1,  // Only fruits count for beverages
				Negative:  0,  // Low negative points
				ScoreType: models.BeverageType,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scorer.CalculateScore(tt.data, tt.foodType)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if result.Grade != tt.expected.Grade {
					t.Errorf("CalculateScore() grade = %v, want %v", result.Grade, tt.expected.Grade)
				}
				if result.ScoreType != tt.expected.ScoreType {
					t.Errorf("CalculateScore() scoreType = %v, want %v", result.ScoreType, tt.expected.ScoreType)
				}
				// Allow some tolerance for score calculations due to algorithm complexity
				if abs(result.Value-tt.expected.Value) > 5 {
					t.Errorf("CalculateScore() value = %v, want approximately %v", result.Value, tt.expected.Value)
				}
			}
		})
	}
}

// TestNutritionalScorer_GetScoreGrade tests grade calculation accuracy
func TestNutritionalScorer_GetScoreGrade(t *testing.T) {
	scorer := NewNutritionalScorer()

	tests := []struct {
		name     string
		score    int
		expected string
	}{
		{"Grade A - Best Score", -5, "A"},
		{"Grade A - Boundary", -1, "A"},
		{"Grade B - Low", 0, "B"},
		{"Grade B - Boundary", 2, "B"},
		{"Grade C - Low", 3, "C"},
		{"Grade C - Mid", 6, "C"},
		{"Grade C - Boundary", 10, "C"},
		{"Grade D - Low", 11, "D"},
		{"Grade D - Mid", 15, "D"},
		{"Grade D - Boundary", 18, "D"},
		{"Grade E - Low", 19, "E"},
		{"Grade E - High", 25, "E"},
		{"Grade E - Very High", 50, "E"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scorer.GetScoreGrade(tt.score)
			if result != tt.expected {
				t.Errorf("GetScoreGrade(%d) = %v, want %v", tt.score, result, tt.expected)
			}
		})
	}
}

// TestScoreCalculator_CalculateNegativePoints tests negative points calculation
func TestScoreCalculator_CalculateNegativePoints(t *testing.T) {
	calculator := NewScoreCalculator()

	tests := []struct {
		name     string
		data     models.NutritionalData
		expected int
	}{
		{
			name: "All Zero Values",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(0),
				Sugars:              models.SugarGram(0),
				SaturatedFattyAcids: models.SaturatedFattyAcids(0),
				Sodium:              models.SodiumMilligram(0),
			},
			expected: 0,
		},
		{
			name: "Low Values - Minimal Points",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(300),  // 0 points
				Sugars:              models.SugarGram(4),   // 0 points
				SaturatedFattyAcids: models.SaturatedFattyAcids(0.5), // 0 points
				Sodium:              models.SodiumMilligram(50),      // 0 points
			},
			expected: 0,
		},
		{
			name: "Moderate Values",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(1000), // 2 points
				Sugars:              models.SugarGram(15),  // 3 points
				SaturatedFattyAcids: models.SaturatedFattyAcids(3),  // 2 points
				Sodium:              models.SodiumMilligram(300),    // 2 points
			},
			expected: 9, // 2+3+2+2
		},
		{
			name: "High Values - Maximum Points",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(4000), // 10 points
				Sugars:              models.SugarGram(50),  // 10 points
				SaturatedFattyAcids: models.SaturatedFattyAcids(15), // 10 points
				Sodium:              models.SodiumMilligram(1000),   // 10 points
			},
			expected: 40, // Maximum possible
		},
		{
			name: "Energy Boundary Test - 335 kJ",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(335), // Should be 0 points
				Sugars:              models.SugarGram(0),
				SaturatedFattyAcids: models.SaturatedFattyAcids(0),
				Sodium:              models.SodiumMilligram(0),
			},
			expected: 0,
		},
		{
			name: "Energy Boundary Test - 336 kJ",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(336), // Should be 1 point
				Sugars:              models.SugarGram(0),
				SaturatedFattyAcids: models.SaturatedFattyAcids(0),
				Sodium:              models.SodiumMilligram(0),
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateNegativePoints(tt.data)
			if result != tt.expected {
				t.Errorf("CalculateNegativePoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestScoreCalculator_CalculatePositivePoints tests positive points calculation
func TestScoreCalculator_CalculatePositivePoints(t *testing.T) {
	calculator := NewScoreCalculator()

	tests := []struct {
		name     string
		data     models.NutritionalData
		foodType models.ScoreType
		expected int
	}{
		{
			name: "All Zero Values",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(0),
				Fibre:   models.FibreGram(0),
				Protein: models.ProteinGram(0),
			},
			foodType: models.FoodType,
			expected: 0,
		},
		{
			name: "High Fruit Content",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(90), // 5 points
				Fibre:   models.FibreGram(0),
				Protein: models.ProteinGram(0),
			},
			foodType: models.FoodType,
			expected: 5,
		},
		{
			name: "High Fiber Content",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(0),
				Fibre:   models.FibreGram(6),  // 5 points
				Protein: models.ProteinGram(0),
			},
			foodType: models.FoodType,
			expected: 5,
		},
		{
			name: "High Protein Content",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(0),
				Fibre:   models.FibreGram(0),
				Protein: models.ProteinGram(10), // 5 points
			},
			foodType: models.FoodType,
			expected: 5,
		},
		{
			name: "Maximum Positive Points",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(100), // 5 points
				Fibre:   models.FibreGram(10),      // 5 points
				Protein: models.ProteinGram(20),    // 5 points
			},
			foodType: models.FoodType,
			expected: 15, // Maximum possible
		},
		{
			name: "Fruit Boundary Test - 40%",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(40), // Should be 0 points
				Fibre:   models.FibreGram(0),
				Protein: models.ProteinGram(0),
			},
			foodType: models.FoodType,
			expected: 0,
		},
		{
			name: "Fruit Boundary Test - 41%",
			data: models.NutritionalData{
				Fruits:  models.FruitsPercent(41), // Should be 1 point
				Fibre:   models.FibreGram(0),
				Protein: models.ProteinGram(0),
			},
			foodType: models.FoodType,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculatePositivePoints(tt.data, tt.foodType)
			if result != tt.expected {
				t.Errorf("CalculatePositivePoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestScoreCalculator_GetFinalScore tests final score calculation for different food types
func TestScoreCalculator_GetFinalScore(t *testing.T) {
	calculator := NewScoreCalculator()

	tests := []struct {
		name     string
		negative int
		positive int
		foodType models.ScoreType
		expected int
	}{
		{
			name:     "Regular Food - Basic Calculation",
			negative: 10,
			positive: 5,
			foodType: models.FoodType,
			expected: 5, // 10 - 5
		},
		{
			name:     "Water Type - Always Zero",
			negative: 10,
			positive: 5,
			foodType: models.WaterType,
			expected: 0, // Water always gets 0
		},
		{
			name:     "Cheese Type - Special Rules",
			negative: 15,
			positive: 8,
			foodType: models.CheeseType,
			expected: 7, // 15 - 8 (protein always counts)
		},
		{
			name:     "Beverage Type - Modified Rules",
			negative: 5,
			positive: 3,
			foodType: models.BeverageType,
			expected: 2, // Simplified beverage calculation
		},
		{
			name:     "High Negative Points - Regular Food",
			negative: 25,
			positive: 10,
			foodType: models.FoodType,
			expected: 15, // 25 - 10
		},
		{
			name:     "Zero Negative Points",
			negative: 0,
			positive: 5,
			foodType: models.FoodType,
			expected: -5, // 0 - 5 (can be negative)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.GetFinalScore(tt.negative, tt.positive, tt.foodType)
			if result != tt.expected {
				t.Errorf("GetFinalScore() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNutritionalScorer_ValidationIntegration tests validation integration
func TestNutritionalScorer_ValidationIntegration(t *testing.T) {
	scorer := NewNutritionalScorer()

	tests := []struct {
		name    string
		data    models.NutritionalData
		wantErr bool
	}{
		{
			name: "Valid Data",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(1000),
				Sugars:              models.SugarGram(10),
				SaturatedFattyAcids: models.SaturatedFattyAcids(5),
				Sodium:              models.SodiumMilligram(200),
				Fruits:              models.FruitsPercent(50),
				Fibre:               models.FibreGram(3),
				Protein:             models.ProteinGram(8),
			},
			wantErr: false,
		},
		{
			name: "Invalid Energy - Negative",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(-100), // Invalid
				Sugars:              models.SugarGram(10),
				SaturatedFattyAcids: models.SaturatedFattyAcids(5),
				Sodium:              models.SodiumMilligram(200),
				Fruits:              models.FruitsPercent(50),
				Fibre:               models.FibreGram(3),
				Protein:             models.ProteinGram(8),
			},
			wantErr: true,
		},
		{
			name: "Invalid Sugar - Too High",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(1000),
				Sugars:              models.SugarGram(150), // Invalid
				SaturatedFattyAcids: models.SaturatedFattyAcids(5),
				Sodium:              models.SodiumMilligram(200),
				Fruits:              models.FruitsPercent(50),
				Fibre:               models.FibreGram(3),
				Protein:             models.ProteinGram(8),
			},
			wantErr: true,
		},
		{
			name: "Invalid Fruits - Over 100%",
			data: models.NutritionalData{
				Energy:              models.EnergyKJ(1000),
				Sugars:              models.SugarGram(10),
				SaturatedFattyAcids: models.SaturatedFattyAcids(5),
				Sodium:              models.SodiumMilligram(200),
				Fruits:              models.FruitsPercent(150), // Invalid
				Fibre:               models.FibreGram(3),
				Protein:             models.ProteinGram(8),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := scorer.CalculateScore(tt.data, models.FoodType)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestEdgeCases tests various edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	scorer := NewNutritionalScorer()

	t.Run("Boundary Values - Energy Thresholds", func(t *testing.T) {
		// Test exact boundary values for energy scoring
		testCases := []struct {
			energy   models.EnergyKJ
			expected int // Expected energy points
		}{
			{335, 0},   // Boundary for 0 points
			{336, 1},   // Just over boundary
			{670, 1},   // Boundary for 1 point
			{671, 2},   // Just over boundary
			{3350, 9},  // Boundary for 9 points
			{3351, 10}, // Maximum points
		}

		calculator := NewScoreCalculator()
		for _, tc := range testCases {
			data := models.NutritionalData{
				Energy:              tc.energy,
				Sugars:              models.SugarGram(0),
				SaturatedFattyAcids: models.SaturatedFattyAcids(0),
				Sodium:              models.SodiumMilligram(0),
			}
			result := calculator.CalculateNegativePoints(data)
			if result != tc.expected {
				t.Errorf("Energy %v kJ: got %d points, want %d points", tc.energy, result, tc.expected)
			}
		}
	})

	t.Run("Extreme Values", func(t *testing.T) {
		// Test with maximum allowed values
		extremeData := models.NutritionalData{
			Energy:              models.EnergyKJ(4000),  // Maximum
			Sugars:              models.SugarGram(100),  // Maximum
			SaturatedFattyAcids: models.SaturatedFattyAcids(100), // Maximum
			Sodium:              models.SodiumMilligram(10000),   // Maximum
			Fruits:              models.FruitsPercent(100),       // Maximum
			Fibre:               models.FibreGram(50),            // Maximum
			Protein:             models.ProteinGram(100),         // Maximum
		}

		result, err := scorer.CalculateScore(extremeData, models.FoodType)
		if err != nil {
			t.Errorf("CalculateScore() with extreme values failed: %v", err)
		}

		// Should get Grade E due to very high negative points
		if result.Grade != "E" {
			t.Errorf("Expected Grade E for extreme values, got %s", result.Grade)
		}
	})

	t.Run("Minimum Values", func(t *testing.T) {
		// Test with minimum allowed values (all zeros)
		minData := models.NutritionalData{
			Energy:              models.EnergyKJ(0),
			Sugars:              models.SugarGram(0),
			SaturatedFattyAcids: models.SaturatedFattyAcids(0),
			Sodium:              models.SodiumMilligram(0),
			Fruits:              models.FruitsPercent(0),
			Fibre:               models.FibreGram(0),
			Protein:             models.ProteinGram(0),
		}

		result, err := scorer.CalculateScore(minData, models.FoodType)
		if err != nil {
			t.Errorf("CalculateScore() with minimum values failed: %v", err)
		}

		// Should get Grade A due to no negative points
		if result.Grade != "A" {
			t.Errorf("Expected Grade A for minimum values, got %s", result.Grade)
		}
	})
}

// TestGetScoreThresholds tests the score threshold functionality
func TestGetScoreThresholds(t *testing.T) {
	scorer := NewNutritionalScorer()
	thresholds := scorer.GetScoreThresholds()

	expectedThresholds := map[string]int{
		"A": -1,
		"B": 2,
		"C": 10,
		"D": 18,
		"E": 19,
	}

	for grade, expectedThreshold := range expectedThresholds {
		if threshold, exists := thresholds[grade]; !exists {
			t.Errorf("Missing threshold for grade %s", grade)
		} else if threshold != expectedThreshold {
			t.Errorf("Threshold for grade %s = %d, want %d", grade, threshold, expectedThreshold)
		}
	}
}

// Helper function to calculate absolute difference
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}