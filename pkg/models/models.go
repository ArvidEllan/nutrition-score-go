package models

import (
	"time"
)

// ScoreType represents different types of food/beverage categories for nutritional scoring
// This enum is used to apply different scoring rules based on the food category
type ScoreType int

const (
	FoodType     ScoreType = iota // Regular food items
	BeverageType                  // Liquid beverages
	WaterType                     // Water (special case with no scoring)
	CheeseType                    // Cheese products (may have different scoring rules)
)

// String returns the string representation of ScoreType for better display
func (st ScoreType) String() string {
	switch st {
	case FoodType:
		return "Food"
	case BeverageType:
		return "Beverage"
	case WaterType:
		return "Water"
	case CheeseType:
		return "Cheese"
	default:
		return "Unknown"
	}
}

// NutritionalScore holds the calculated nutritional score and its components
// This struct contains the final score calculation results and breakdown
type NutritionalScore struct {
	Value     int       `json:"value"`      // Final calculated score (negative - positive)
	Grade     string    `json:"grade"`      // Letter grade (A, B, C, D, E)
	Positive  int       `json:"positive"`   // Sum of positive nutritional points (beneficial nutrients)
	Negative  int       `json:"negative"`   // Sum of negative nutritional points (nutrients to limit)
	ScoreType ScoreType `json:"score_type"` // Category of the food/beverage being scored
}

// EnergyKJ represents energy content in kilojoules
// Higher energy content contributes to negative (unhealthy) points
type EnergyKJ float64

// SugarGram represents sugar content in grams
// Higher sugar content contributes to negative (unhealthy) points
type SugarGram float64

// SaturatedFattyAcids represents saturated fat content in grams
// Higher saturated fat contributes to negative (unhealthy) points
type SaturatedFattyAcids float64

// SodiumMilligram represents sodium content in milligrams
// Higher sodium content contributes to negative (unhealthy) points
type SodiumMilligram float64

// FruitsPercent represents the percentage of fruits/vegetables/nuts
// Higher fruit/vegetable content contributes to positive (healthy) points
type FruitsPercent float64

// FibreGram represents fiber content in grams
// Higher fiber content contributes to positive (healthy) points
type FibreGram float64

// ProteinGram represents protein content in grams
// Higher protein content contributes to positive (healthy) points
type ProteinGram float64

// NutritionalData contains all the nutritional information needed for scoring
// This struct holds the complete nutritional profile of a food item per 100g
type NutritionalData struct {
	Energy              EnergyKJ            `json:"energy"`                // Energy content in kJ per 100g
	Sugars              SugarGram           `json:"sugars"`                // Sugar content in grams per 100g
	SaturatedFattyAcids SaturatedFattyAcids `json:"saturated_fatty_acids"` // Saturated fat content in grams per 100g
	Sodium              SodiumMilligram     `json:"sodium"`                // Sodium content in milligrams per 100g
	Fruits              FruitsPercent       `json:"fruits"`                // Fruits/vegetables/nuts percentage
	Fibre               FibreGram           `json:"fibre"`                 // Fiber content in grams per 100g
	Protein             ProteinGram         `json:"protein"`               // Protein content in grams per 100g
}

// Food represents a food item with its nutritional data and metadata
// This struct can represent both database foods and user-defined foods
type Food struct {
	ID               string          `json:"id"`                 // Unique identifier for the food
	Name             string          `json:"name"`               // Display name of the food
	Category         string          `json:"category"`           // Food category (e.g., "Fruits", "Dairy", "Grains")
	Brand            string          `json:"brand,omitempty"`    // Brand name (optional, for packaged foods)
	NutritionalData  NutritionalData `json:"nutritional_data"`   // Complete nutritional profile
	IsUserDefined    bool            `json:"is_user_defined"`    // True if created by user, false if from database
	CreatedAt        time.Time       `json:"created_at"`         // When the food was added to the system
	UpdatedAt        time.Time       `json:"updated_at"`         // When the food was last modified
	Source           string          `json:"source,omitempty"`   // Data source (e.g., "USDA", "User Input")
}

// NutritionalAnalysis represents a complete analysis of a food item
// This struct contains the food data, calculated score, and analysis metadata
type NutritionalAnalysis struct {
	ID              string           `json:"id"`               // Unique identifier for this analysis
	Food            Food             `json:"food"`             // The food item that was analyzed
	Score           NutritionalScore `json:"score"`            // Calculated nutritional score and breakdown
	AnalyzedAt      time.Time        `json:"analyzed_at"`      // When the analysis was performed
	Notes           string           `json:"notes,omitempty"`  // Optional user notes about the analysis
	ServingSize     float64          `json:"serving_size"`     // Serving size in grams (default 100g)
	UserID          string           `json:"user_id,omitempty"` // User who performed the analysis (for multi-user systems)
}

// FoodComparison represents a comparison between multiple food items
// This struct contains the foods being compared and the comparison results
type FoodComparison struct {
	ID              string                `json:"id"`               // Unique identifier for this comparison
	Foods           []Food                `json:"foods"`            // List of foods being compared
	Analyses        []NutritionalAnalysis `json:"analyses"`         // Analysis results for each food
	BestChoice      *Food                 `json:"best_choice,omitempty"` // Food with the best nutritional score (lowest value)
	WorstChoice     *Food                 `json:"worst_choice,omitempty"` // Food with the worst nutritional score (highest value)
	ComparedAt      time.Time             `json:"compared_at"`      // When the comparison was performed
	ComparisonNotes string                `json:"comparison_notes,omitempty"` // Optional notes about the comparison
	UserID          string                `json:"user_id,omitempty"` // User who performed the comparison
}

// HistoryFilter represents filtering options for analysis history
// This struct is used to filter historical analyses by various criteria
type HistoryFilter struct {
	StartDate    *time.Time `json:"start_date,omitempty"`    // Filter analyses after this date
	EndDate      *time.Time `json:"end_date,omitempty"`      // Filter analyses before this date
	FoodCategory string     `json:"food_category,omitempty"` // Filter by food category
	ScoreRange   *ScoreRange `json:"score_range,omitempty"`  // Filter by score range
	UserID       string     `json:"user_id,omitempty"`       // Filter by user (for multi-user systems)
	Limit        int        `json:"limit,omitempty"`         // Maximum number of results to return
}

// ScoreRange represents a range of nutritional scores for filtering
// Used to filter analyses by score values (e.g., only show healthy foods)
type ScoreRange struct {
	Min *int `json:"min,omitempty"` // Minimum score (inclusive)
	Max *int `json:"max,omitempty"` // Maximum score (inclusive)
}

// ExportFormat represents the available export formats
// This enum defines the supported formats for data export
type ExportFormat int

const (
	JSON ExportFormat = iota // JSON format export
	CSV                      // CSV format export
	XML                      // XML format export (future enhancement)
)

// String returns the string representation of ExportFormat
func (ef ExportFormat) String() string {
	switch ef {
	case JSON:
		return "JSON"
	case CSV:
		return "CSV"
	case XML:
		return "XML"
	default:
		return "Unknown"
	}
}

// FileExtension returns the appropriate file extension for the export format
func (ef ExportFormat) FileExtension() string {
	switch ef {
	case JSON:
		return ".json"
	case CSV:
		return ".csv"
	case XML:
		return ".xml"
	default:
		return ".txt"
	}
}

// ExportData represents data prepared for export
// This struct contains the data and metadata for export operations
type ExportData struct {
	Format      ExportFormat `json:"format"`       // Export format used
	ExportedAt  time.Time    `json:"exported_at"`  // When the export was performed
	DataType    string       `json:"data_type"`    // Type of data being exported (e.g., "analyses", "foods")
	RecordCount int          `json:"record_count"` // Number of records in the export
	UserID      string       `json:"user_id,omitempty"` // User who performed the export
}



// NutritionalDataValidation contains validation rules for nutritional data
// This struct defines the acceptable ranges for each nutritional component
type NutritionalDataValidation struct {
	EnergyMin              float64 `json:"energy_min"`                // Minimum energy in kJ per 100g
	EnergyMax              float64 `json:"energy_max"`                // Maximum energy in kJ per 100g
	SugarsMin              float64 `json:"sugars_min"`                // Minimum sugar in g per 100g
	SugarsMax              float64 `json:"sugars_max"`                // Maximum sugar in g per 100g
	SaturatedFatMin        float64 `json:"saturated_fat_min"`         // Minimum saturated fat in g per 100g
	SaturatedFatMax        float64 `json:"saturated_fat_max"`         // Maximum saturated fat in g per 100g
	SodiumMin              float64 `json:"sodium_min"`                // Minimum sodium in mg per 100g
	SodiumMax              float64 `json:"sodium_max"`                // Maximum sodium in mg per 100g
	FruitsMin              float64 `json:"fruits_min"`                // Minimum fruits percentage
	FruitsMax              float64 `json:"fruits_max"`                // Maximum fruits percentage
	FibreMin               float64 `json:"fibre_min"`                 // Minimum fiber in g per 100g
	FibreMax               float64 `json:"fibre_max"`                 // Maximum fiber in g per 100g
	ProteinMin             float64 `json:"protein_min"`               // Minimum protein in g per 100g
	ProteinMax             float64 `json:"protein_max"`               // Maximum protein in g per 100g
}

// DefaultValidationRules returns the default validation rules for nutritional data
// These rules are based on realistic ranges for food nutritional content
func DefaultValidationRules() NutritionalDataValidation {
	return NutritionalDataValidation{
		EnergyMin:       0,     // 0 kJ per 100g (minimum)
		EnergyMax:       4000,  // 4000 kJ per 100g (very high energy foods like oils)
		SugarsMin:       0,     // 0g per 100g
		SugarsMax:       100,   // 100g per 100g (pure sugar)
		SaturatedFatMin: 0,     // 0g per 100g
		SaturatedFatMax: 100,   // 100g per 100g (pure fat)
		SodiumMin:       0,     // 0mg per 100g
		SodiumMax:       10000, // 10000mg per 100g (very high sodium foods)
		FruitsMin:       0,     // 0% fruits/vegetables/nuts
		FruitsMax:       100,   // 100% fruits/vegetables/nuts
		FibreMin:        0,     // 0g per 100g
		FibreMax:        50,    // 50g per 100g (very high fiber foods)
		ProteinMin:      0,     // 0g per 100g
		ProteinMax:      100,   // 100g per 100g (pure protein)
	}
}