package main

// ScoreType represents different types of food/beverage categories for nutritional scoring
// This enum is used to apply different scoring rules based on the food category
type ScoreType int

const (
	Food     ScoreType = iota // Regular food items
	Beverage                  // Liquid beverages
	Water                     // Water (special case with no scoring)
	Cheese                    // Cheese products (may have different scoring rules)
)

// NutritionalScore holds the calculated nutritional score and its components
// This struct contains the final score calculation results and breakdown
type NutritionalScore struct {
	Value     int       // Final calculated score (negative - positive)
	Positive  int       // Sum of positive nutritional points (beneficial nutrients)
	Negative  int       // Sum of negative nutritional points (nutrients to limit)
	ScoreType ScoreType // Category of the food/beverage being scored
}

// EnergyKJ represents energy content in kilojoules
// Higher energy content contributes to negative (unhealthy) points
type EnergyKJ float64

// GetPoints calculates negative points based on energy content
// Uses simplified thresholds - will be enhanced with official Nutri-Score algorithm later
func (e EnergyKJ) GetPoints() int {
	// Placeholder logic - simplified scoring for now
	// Higher energy = more negative points
	if e > 1000 {
		return 5
	}
	return 1
}

// SugarGram represents sugar content in grams
// Higher sugar content contributes to negative (unhealthy) points
type SugarGram float64

// GetPoints calculates negative points based on sugar content
// Simplified scoring - will be enhanced with official thresholds later
func (s SugarGram) GetPoints() int {
	// Higher sugar = more negative points
	if s > 10 {
		return 5
	}
	return 1
}

// SaturatedFattyAcids represents saturated fat content in grams
// Higher saturated fat contributes to negative (unhealthy) points
type SaturatedFattyAcids float64

// GetPoints calculates negative points based on saturated fatty acids content
// Simplified scoring - will be enhanced with official thresholds later
func (s SaturatedFattyAcids) GetPoints() int {
	// Higher saturated fat = more negative points
	if s > 5 {
		return 5
	}
	return 1
}

// SodiumMilligram represents sodium content in milligrams
// Higher sodium content contributes to negative (unhealthy) points
type SodiumMilligram float64

// GetPoints calculates negative points based on sodium content
// Simplified scoring - will be enhanced with official thresholds later
func (s SodiumMilligram) GetPoints() int {
	// Higher sodium = more negative points
	if s > 500 {
		return 5
	}
	return 1
}

// FruitsPercent represents the percentage of fruits/vegetables/nuts
// Higher fruit/vegetable content contributes to positive (healthy) points
type FruitsPercent float64

// GetPoints calculates positive points based on fruits/vegetables/nuts percentage
// Simplified scoring - will be enhanced with official thresholds later
func (f FruitsPercent) GetPoints() int {
	// Higher fruit/vegetable percentage = more positive points
	if f > 80 {
		return 5
	}
	return 1
}

// FibreGram represents fiber content in grams
// Higher fiber content contributes to positive (healthy) points
type FibreGram float64

// GetPoints calculates positive points based on fiber content
// Simplified scoring - will be enhanced with official thresholds later
func (f FibreGram) GetPoints() int {
	// Higher fiber = more positive points
	if f > 5 {
		return 5
	}
	return 1
}

// ProteinGram represents protein content in grams
// Higher protein content contributes to positive (healthy) points
type ProteinGram float64

// GetPoints calculates positive points based on protein content
// Simplified scoring - will be enhanced with official thresholds later
func (p ProteinGram) GetPoints() int {
	// Higher protein = more positive points
	if p > 5 {
		return 5
	}
	return 1
}

// NutritionalData contains all the nutritional information needed for scoring
// This struct holds the complete nutritional profile of a food item
type NutritionalData struct {
	Energy              EnergyKJ            // Energy content in kJ
	Sugars              SugarGram           // Sugar content in grams
	SaturatedFattyAcids SaturatedFattyAcids // Saturated fat content in grams
	Sodium              SodiumMilligram     // Sodium content in milligrams
	Fruits              FruitsPercent       // Fruits/vegetables/nuts percentage
	Fibre               FibreGram           // Fiber content in grams
	Protein             ProteinGram         // Protein content in grams
}

// GetNutritionalScore calculates the nutritional score based on the provided data and score type
// This is the main scoring function that implements the nutritional scoring algorithm
func GetNutritionalScore(n NutritionalData, st ScoreType) NutritionalScore {
	// Initialize scoring variables
	value := 0    // Final score value
	positive := 0 // Positive points (beneficial nutrients)
	negative := 0 // Negative points (nutrients to limit)

	// Water doesn't get scored, all other types do
	// This follows the Nutri-Score standard where water has no nutritional score
	if st != Water {
		// Calculate positive points from beneficial nutrients
		fruitPoints := n.Fruits.GetPoints()
		fibrePoints := n.Fibre.GetPoints()
		proteinPoints := n.Protein.GetPoints()
		positive = fruitPoints + fibrePoints + proteinPoints
		
		// Calculate negative points from nutrients to limit
		// These are nutrients that should be consumed in moderation
		negative = n.Energy.GetPoints() + n.Sugars.GetPoints() + 
				  n.SaturatedFattyAcids.GetPoints() + n.Sodium.GetPoints()
		
		// Final score is negative points minus positive points
		// Lower scores are better (more healthy)
		value = negative - positive
	}

	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}