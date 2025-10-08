package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nutritional-score/pkg/models"
)

// UserFoodData represents the structure of the user foods JSON file
type UserFoodData struct {
	Version     string        `json:"version"`
	LastUpdated time.Time     `json:"last_updated"`
	Foods       []models.Food `json:"foods"`
}

// JSONUserFoodRepository implements the UserFoodRepository interface using JSON file storage
type JSONUserFoodRepository struct {
	data     *UserFoodData
	filePath string
	loaded   bool
}

// NewJSONUserFoodRepository creates a new instance of the JSON user food repository
func NewJSONUserFoodRepository(filePath string) *JSONUserFoodRepository {
	return &JSONUserFoodRepository{
		filePath: filePath,
		loaded:   false,
	}
}

// loadData loads user food data from the JSON file
func (repo *JSONUserFoodRepository) loadData() error {
	// Check if file exists
	if _, err := os.Stat(repo.filePath); os.IsNotExist(err) {
		// Create empty data structure if file doesn't exist
		repo.data = &UserFoodData{
			Version:     "1.0",
			LastUpdated: time.Now(),
			Foods:       []models.Food{},
		}
		repo.loaded = true
		return repo.saveData()
	}

	// Read the file
	fileData, err := os.ReadFile(repo.filePath)
	if err != nil {
		return fmt.Errorf("failed to read user foods file: %w", err)
	}

	// Parse JSON data
	var data UserFoodData
	if err := json.Unmarshal(fileData, &data); err != nil {
		return fmt.Errorf("failed to parse user foods JSON: %w", err)
	}

	repo.data = &data
	repo.loaded = true
	return nil
}

// saveData saves user food data to the JSON file
func (repo *JSONUserFoodRepository) saveData() error {
	if repo.data == nil {
		return fmt.Errorf("no data to save")
	}

	// Update last modified time
	repo.data.LastUpdated = time.Now()

	// Ensure directory exists
	dir := filepath.Dir(repo.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal data to JSON
	jsonData, err := json.MarshalIndent(repo.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user foods data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(repo.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write user foods file: %w", err)
	}

	return nil
}

// ensureLoaded ensures that the data is loaded before performing operations
func (repo *JSONUserFoodRepository) ensureLoaded() error {
	if !repo.loaded {
		return repo.loadData()
	}
	return nil
}

// SaveFood stores a new user-defined food or updates an existing one
func (repo *JSONUserFoodRepository) SaveFood(ctx context.Context, food models.Food) error {
	if err := repo.ensureLoaded(); err != nil {
		return err
	}

	// Generate ID if not provided
	if food.ID == "" {
		food.ID = uuid.New().String()
	}

	// Set user-defined flag and timestamps
	food.IsUserDefined = true
	now := time.Now()
	
	// Check if food already exists (update case)
	for i, existingFood := range repo.data.Foods {
		if existingFood.ID == food.ID {
			food.CreatedAt = existingFood.CreatedAt // Preserve original creation time
			food.UpdatedAt = now
			repo.data.Foods[i] = food
			return repo.saveData()
		}
	}

	// New food case
	food.CreatedAt = now
	food.UpdatedAt = now
	repo.data.Foods = append(repo.data.Foods, food)

	return repo.saveData()
}

// GetUserFoods retrieves all foods created by users
func (repo *JSONUserFoodRepository) GetUserFoods(ctx context.Context) ([]models.Food, error) {
	if err := repo.ensureLoaded(); err != nil {
		return nil, err
	}

	// Return a copy of the foods slice to prevent external modification
	foods := make([]models.Food, len(repo.data.Foods))
	copy(foods, repo.data.Foods)

	return foods, nil
}

// GetUserFoodByID retrieves a specific user-defined food by ID
func (repo *JSONUserFoodRepository) GetUserFoodByID(ctx context.Context, id string) (models.Food, error) {
	if err := repo.ensureLoaded(); err != nil {
		return models.Food{}, err
	}

	if id == "" {
		return models.Food{}, fmt.Errorf("food ID cannot be empty")
	}

	for _, food := range repo.data.Foods {
		if food.ID == id {
			return food, nil
		}
	}

	return models.Food{}, fmt.Errorf("user food not found with ID: %s", id)
}

// UpdateFood modifies an existing user-defined food
func (repo *JSONUserFoodRepository) UpdateFood(ctx context.Context, id string, food models.Food) error {
	if err := repo.ensureLoaded(); err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("food ID cannot be empty")
	}

	// Find and update the food
	for i, existingFood := range repo.data.Foods {
		if existingFood.ID == id {
			// Preserve ID, creation time, and user-defined flag
			food.ID = id
			food.CreatedAt = existingFood.CreatedAt
			food.IsUserDefined = true
			food.UpdatedAt = time.Now()
			
			repo.data.Foods[i] = food
			return repo.saveData()
		}
	}

	return fmt.Errorf("user food not found with ID: %s", id)
}

// DeleteFood removes a user-defined food from storage
func (repo *JSONUserFoodRepository) DeleteFood(ctx context.Context, id string) error {
	if err := repo.ensureLoaded(); err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("food ID cannot be empty")
	}

	// Find and remove the food
	for i, food := range repo.data.Foods {
		if food.ID == id {
			// Remove the food from the slice
			repo.data.Foods = append(repo.data.Foods[:i], repo.data.Foods[i+1:]...)
			return repo.saveData()
		}
	}

	return fmt.Errorf("user food not found with ID: %s", id)
}

// SearchUserFoods finds user-defined foods matching the query
func (repo *JSONUserFoodRepository) SearchUserFoods(ctx context.Context, query string) ([]models.Food, error) {
	if err := repo.ensureLoaded(); err != nil {
		return nil, err
	}

	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []models.Food

	for _, food := range repo.data.Foods {
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

// GetUserFoodCount returns the number of user-defined foods
func (repo *JSONUserFoodRepository) GetUserFoodCount(ctx context.Context) (int, error) {
	if err := repo.ensureLoaded(); err != nil {
		return 0, err
	}

	return len(repo.data.Foods), nil
}

// GetDefaultUserFoodsPath returns the default path for user foods storage
func GetDefaultUserFoodsPath() string {
	return filepath.Join("data", "user_foods.json")
}