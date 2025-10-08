package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nutritional-score/pkg/models"
)

// FoodDatabaseData represents the structure of the embedded food database JSON file
type FoodDatabaseData struct {
	Version     string        `json:"version"`
	LastUpdated time.Time     `json:"last_updated"`
	Description string        `json:"description"`
	Foods       []models.Food `json:"foods"`
}

// EmbeddedFoodDatabase implements the FoodDatabase interface for the embedded food database
type EmbeddedFoodDatabase struct {
	data         *FoodDatabaseData
	databasePath string
	loaded       bool
}

// NewEmbeddedFoodDatabase creates a new instance of the embedded food database
func NewEmbeddedFoodDatabase(databasePath string) *EmbeddedFoodDatabase {
	return &EmbeddedFoodDatabase{
		databasePath: databasePath,
		loaded:       false,
	}
}

// LoadDatabase initializes the food database from storage
func (db *EmbeddedFoodDatabase) LoadDatabase(ctx context.Context) error {
	// Check if database file exists
	if _, err := os.Stat(db.databasePath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found: %s", db.databasePath)
	}

	// Read the database file
	fileData, err := os.ReadFile(db.databasePath)
	if err != nil {
		return fmt.Errorf("failed to read database file: %w", err)
	}

	// Parse JSON data
	var data FoodDatabaseData
	if err := json.Unmarshal(fileData, &data); err != nil {
		return fmt.Errorf("failed to parse database JSON: %w", err)
	}

	// Validate that we have foods
	if len(data.Foods) == 0 {
		return fmt.Errorf("database contains no foods")
	}

	// Store the loaded data
	db.data = &data
	db.loaded = true

	return nil
}

// SearchFoods finds foods matching the given query string
func (db *EmbeddedFoodDatabase) SearchFoods(ctx context.Context, query string) ([]models.Food, error) {
	if !db.loaded {
		return nil, fmt.Errorf("database not loaded")
	}

	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []models.Food

	for _, food := range db.data.Foods {
		// Search in food name
		if strings.Contains(strings.ToLower(food.Name), query) {
			results = append(results, food)
			continue
		}

		// Search in category
		if strings.Contains(strings.ToLower(food.Category), query) {
			results = append(results, food)
			continue
		}

		// Search in brand (if not empty)
		if food.Brand != "" && strings.Contains(strings.ToLower(food.Brand), query) {
			results = append(results, food)
			continue
		}
	}

	return results, nil
}

// GetFoodByID retrieves a specific food by its unique identifier
func (db *EmbeddedFoodDatabase) GetFoodByID(ctx context.Context, id string) (models.Food, error) {
	if !db.loaded {
		return models.Food{}, fmt.Errorf("database not loaded")
	}

	if id == "" {
		return models.Food{}, fmt.Errorf("food ID cannot be empty")
	}

	for _, food := range db.data.Foods {
		if food.ID == id {
			return food, nil
		}
	}

	return models.Food{}, fmt.Errorf("food not found with ID: %s", id)
}

// GetAllFoods returns all foods in the database
func (db *EmbeddedFoodDatabase) GetAllFoods(ctx context.Context) ([]models.Food, error) {
	if !db.loaded {
		return nil, fmt.Errorf("database not loaded")
	}

	// Return a copy of the foods slice to prevent external modification
	foods := make([]models.Food, len(db.data.Foods))
	copy(foods, db.data.Foods)

	return foods, nil
}

// GetFoodsByCategory returns all foods in a specific category
func (db *EmbeddedFoodDatabase) GetFoodsByCategory(ctx context.Context, category string) ([]models.Food, error) {
	if !db.loaded {
		return nil, fmt.Errorf("database not loaded")
	}

	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	category = strings.ToLower(strings.TrimSpace(category))
	var results []models.Food

	for _, food := range db.data.Foods {
		if strings.ToLower(food.Category) == category {
			results = append(results, food)
		}
	}

	return results, nil
}

// GetCategories returns all available food categories
func (db *EmbeddedFoodDatabase) GetCategories(ctx context.Context) ([]string, error) {
	if !db.loaded {
		return nil, fmt.Errorf("database not loaded")
	}

	categoryMap := make(map[string]bool)
	for _, food := range db.data.Foods {
		categoryMap[food.Category] = true
	}

	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}

// GetDatabaseInfo returns information about the loaded database
func (db *EmbeddedFoodDatabase) GetDatabaseInfo() (string, time.Time, int, error) {
	if !db.loaded {
		return "", time.Time{}, 0, fmt.Errorf("database not loaded")
	}

	return db.data.Version, db.data.LastUpdated, len(db.data.Foods), nil
}

// IsLoaded returns whether the database has been loaded
func (db *EmbeddedFoodDatabase) IsLoaded() bool {
	return db.loaded
}

// GetDefaultDatabasePath returns the default path for the embedded food database
func GetDefaultDatabasePath() string {
	return filepath.Join("data", "foods_database.json")
}