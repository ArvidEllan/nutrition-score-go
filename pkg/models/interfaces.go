package models

import (
	"context"
)

// NutritionalScorer defines the interface for nutritional scoring operations
// This interface encapsulates all scoring-related functionality
type NutritionalScorer interface {
	// CalculateScore computes the nutritional score for given data and food type
	// Returns the complete score breakdown including positive/negative points and letter grade
	CalculateScore(data NutritionalData, foodType ScoreType) (NutritionalScore, error)
	
	// ValidateNutritionalData checks if the provided nutritional data is within acceptable ranges
	// Returns a slice of validation errors if any values are invalid
	ValidateNutritionalData(data NutritionalData) []ValidationError
	
	// GetScoreGrade converts a numerical score to a letter grade (A-E)
	// Lower scores get better grades (A is best, E is worst)
	GetScoreGrade(score int) string
	
	// GetScoreThresholds returns the score thresholds for each letter grade
	// Used for displaying grade boundaries to users
	GetScoreThresholds() map[string]int
}

// ScoreCalculator defines the interface for detailed score calculation logic
// This interface handles the mathematical aspects of the Nutri-Score algorithm
type ScoreCalculator interface {
	// CalculateNegativePoints computes points from nutrients that should be limited
	// These include energy, sugars, saturated fat, and sodium
	CalculateNegativePoints(data NutritionalData) int
	
	// CalculatePositivePoints computes points from beneficial nutrients
	// These include fruits/vegetables/nuts, fiber, and protein
	CalculatePositivePoints(data NutritionalData, foodType ScoreType) int
	
	// GetFinalScore combines negative and positive points according to Nutri-Score rules
	// Different food types may have different calculation rules
	GetFinalScore(negative, positive int, foodType ScoreType) int
}

// FoodDatabase defines the interface for food database operations
// This interface handles access to the embedded food database
type FoodDatabase interface {
	// SearchFoods finds foods matching the given query string
	// Searches across food names, categories, and brands
	SearchFoods(ctx context.Context, query string) ([]Food, error)
	
	// GetFoodByID retrieves a specific food by its unique identifier
	GetFoodByID(ctx context.Context, id string) (Food, error)
	
	// GetAllFoods returns all foods in the database
	// Should support pagination for large datasets
	GetAllFoods(ctx context.Context) ([]Food, error)
	
	// GetFoodsByCategory returns all foods in a specific category
	GetFoodsByCategory(ctx context.Context, category string) ([]Food, error)
	
	// GetCategories returns all available food categories
	GetCategories(ctx context.Context) ([]string, error)
	
	// LoadDatabase initializes the food database from storage
	// Should be called during application startup
	LoadDatabase(ctx context.Context) error
}

// UserFoodRepository defines the interface for user-defined food management
// This interface handles CRUD operations for foods created by users
type UserFoodRepository interface {
	// SaveFood stores a new user-defined food or updates an existing one
	SaveFood(ctx context.Context, food Food) error
	
	// GetUserFoods retrieves all foods created by users
	GetUserFoods(ctx context.Context) ([]Food, error)
	
	// GetUserFoodByID retrieves a specific user-defined food by ID
	GetUserFoodByID(ctx context.Context, id string) (Food, error)
	
	// UpdateFood modifies an existing user-defined food
	UpdateFood(ctx context.Context, id string, food Food) error
	
	// DeleteFood removes a user-defined food from storage
	DeleteFood(ctx context.Context, id string) error
	
	// SearchUserFoods finds user-defined foods matching the query
	SearchUserFoods(ctx context.Context, query string) ([]Food, error)
}

// StorageService defines the interface for data persistence operations
// This interface handles all file-based storage operations
type StorageService interface {
	// SaveAnalysis stores a nutritional analysis result
	SaveAnalysis(ctx context.Context, analysis NutritionalAnalysis) error
	
	// GetAnalysisHistory retrieves historical analyses with optional filtering
	GetAnalysisHistory(ctx context.Context, filter HistoryFilter) ([]NutritionalAnalysis, error)
	
	// GetAnalysisByID retrieves a specific analysis by its ID
	GetAnalysisByID(ctx context.Context, id string) (NutritionalAnalysis, error)
	
	// DeleteAnalysis removes an analysis from storage
	DeleteAnalysis(ctx context.Context, id string) error
	
	// SaveComparison stores a food comparison result
	SaveComparison(ctx context.Context, comparison FoodComparison) error
	
	// GetComparisonHistory retrieves historical comparisons
	GetComparisonHistory(ctx context.Context, filter HistoryFilter) ([]FoodComparison, error)
	
	// ExportData exports data in the specified format
	ExportData(ctx context.Context, format ExportFormat, data interface{}) ([]byte, error)
	
	// InitializeStorage sets up the storage directories and files
	InitializeStorage(ctx context.Context) error
}

// CLIInterface defines the interface for command-line user interactions
// This interface handles all user input/output operations
type CLIInterface interface {
	// ShowMainMenu displays the main application menu and returns user choice
	ShowMainMenu() MenuChoice
	
	// GetNutritionalInput prompts user for nutritional data entry
	// Includes validation and retry logic for invalid inputs
	GetNutritionalInput() (NutritionalData, error)
	
	// DisplayScore shows the calculated nutritional score with detailed breakdown
	DisplayScore(result NutritionalAnalysis)
	
	// ShowFoodList displays a list of foods and allows user selection
	ShowFoodList(foods []Food) (Food, error)
	
	// ShowSearchInterface provides food search functionality
	ShowSearchInterface() (string, error)
	
	// DisplayComparison shows the results of food comparison
	DisplayComparison(comparison FoodComparison)
	
	// ShowHistoryInterface displays analysis history with filtering options
	ShowHistoryInterface(analyses []NutritionalAnalysis)
	
	// GetExportOptions prompts user for export format and options
	GetExportOptions() (ExportFormat, string, error)
	
	// ShowError displays error messages in a user-friendly format
	ShowError(err error)
	
	// ShowSuccess displays success messages
	ShowSuccess(message string)
	
	// ConfirmAction asks user to confirm an action (yes/no)
	ConfirmAction(message string) bool
}

// MenuChoice represents the available menu options in the CLI
// This enum defines all possible user actions in the main menu
type MenuChoice int

const (
	MenuCalculateScore MenuChoice = iota // Calculate nutritional score for a food
	MenuSearchFoods                      // Search for foods in database
	MenuManageUserFoods                  // Manage user-defined foods
	MenuCompareFoods                     // Compare multiple foods
	MenuViewHistory                      // View analysis history
	MenuExportData                       // Export analysis data
	MenuSettings                         // Application settings
	MenuExit                             // Exit the application
)

// String returns the string representation of MenuChoice for display
func (mc MenuChoice) String() string {
	switch mc {
	case MenuCalculateScore:
		return "Calculate Nutritional Score"
	case MenuSearchFoods:
		return "Search Foods Database"
	case MenuManageUserFoods:
		return "Manage User Foods"
	case MenuCompareFoods:
		return "Compare Foods"
	case MenuViewHistory:
		return "View Analysis History"
	case MenuExportData:
		return "Export Data"
	case MenuSettings:
		return "Settings"
	case MenuExit:
		return "Exit"
	default:
		return "Unknown"
	}
}

// InputValidator defines the interface for input validation operations
// This interface handles validation of user inputs and data integrity
type InputValidator interface {
	// ValidateNutritionalData validates nutritional data against defined rules
	ValidateNutritionalData(data NutritionalData) []ValidationError
	
	// ValidateFood validates a complete food item
	ValidateFood(food Food) []ValidationError
	
	// ValidateScoreType checks if the score type is valid
	ValidateScoreType(scoreType ScoreType) error
	
	// ValidateSearchQuery checks if a search query is valid
	ValidateSearchQuery(query string) error
	
	// ValidateExportFormat checks if the export format is supported
	ValidateExportFormat(format ExportFormat) error
}

// ComparisonEngine defines the interface for food comparison operations
// This interface handles the logic for comparing multiple foods
type ComparisonEngine interface {
	// CompareFoods analyzes multiple foods and determines the best/worst choices
	CompareFoods(ctx context.Context, foods []Food) (FoodComparison, error)
	
	// RankFoods sorts foods by their nutritional scores (best to worst)
	RankFoods(analyses []NutritionalAnalysis) []NutritionalAnalysis
	
	// GetBestChoice determines the food with the best nutritional profile
	GetBestChoice(analyses []NutritionalAnalysis) *Food
	
	// GetWorstChoice determines the food with the worst nutritional profile
	GetWorstChoice(analyses []NutritionalAnalysis) *Food
	
	// GenerateComparisonSummary creates a summary of the comparison results
	GenerateComparisonSummary(comparison FoodComparison) string
}

// DataExporter defines the interface for data export operations
// This interface handles converting data to various export formats
type DataExporter interface {
	// ExportAnalyses exports analysis data to the specified format
	ExportAnalyses(ctx context.Context, analyses []NutritionalAnalysis, format ExportFormat) ([]byte, error)
	
	// ExportFoods exports food data to the specified format
	ExportFoods(ctx context.Context, foods []Food, format ExportFormat) ([]byte, error)
	
	// ExportComparisons exports comparison data to the specified format
	ExportComparisons(ctx context.Context, comparisons []FoodComparison, format ExportFormat) ([]byte, error)
	
	// GetSupportedFormats returns the list of supported export formats
	GetSupportedFormats() []ExportFormat
	
	// ValidateExportData checks if the data is suitable for export
	ValidateExportData(data interface{}) error
}

// ConfigurationManager defines the interface for application configuration
// This interface handles application settings and preferences
type ConfigurationManager interface {
	// LoadConfiguration loads application settings from storage
	LoadConfiguration(ctx context.Context) error
	
	// SaveConfiguration saves current settings to storage
	SaveConfiguration(ctx context.Context) error
	
	// GetValidationRules returns the current validation rules
	GetValidationRules() NutritionalDataValidation
	
	// SetValidationRules updates the validation rules
	SetValidationRules(rules NutritionalDataValidation) error
	
	// GetDefaultScoreType returns the default score type for calculations
	GetDefaultScoreType() ScoreType
	
	// SetDefaultScoreType sets the default score type
	SetDefaultScoreType(scoreType ScoreType) error
	
	// GetExportDirectory returns the directory for export files
	GetExportDirectory() string
	
	// SetExportDirectory sets the directory for export files
	SetExportDirectory(directory string) error
}