package database

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/nutritional-score/pkg/models"
)

// FoodService provides a unified interface for accessing both embedded and user-defined foods
type FoodService struct {
	embeddedDB   models.FoodDatabase
	userFoodRepo models.UserFoodRepository
}

// NewFoodService creates a new food service with embedded database and user food repository
func NewFoodService(embeddedDB models.FoodDatabase, userFoodRepo models.UserFoodRepository) *FoodService {
	return &FoodService{
		embeddedDB:   embeddedDB,
		userFoodRepo: userFoodRepo,
	}
}

// SearchAllFoods searches across both embedded database and user-defined foods
func (fs *FoodService) SearchAllFoods(ctx context.Context, query string) ([]models.Food, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	var allResults []models.Food

	// Search embedded database
	embeddedResults, err := fs.embeddedDB.SearchFoods(ctx, query)
	if err != nil {
		// Log error but continue with user foods search
		fmt.Printf("Warning: embedded database search failed: %v\n", err)
	} else {
		allResults = append(allResults, embeddedResults...)
	}

	// Search user-defined foods
	userResults, err := fs.userFoodRepo.SearchUserFoods(ctx, query)
	if err != nil {
		// Log error but continue with embedded results
		fmt.Printf("Warning: user foods search failed: %v\n", err)
	} else {
		allResults = append(allResults, userResults...)
	}

	// Sort results by relevance (exact matches first, then partial matches)
	fs.sortSearchResults(allResults, query)

	return allResults, nil
}

// GetFoodByID retrieves a food by ID from either embedded database or user foods
func (fs *FoodService) GetFoodByID(ctx context.Context, id string) (models.Food, error) {
	if id == "" {
		return models.Food{}, fmt.Errorf("food ID cannot be empty")
	}

	// Try embedded database first
	food, err := fs.embeddedDB.GetFoodByID(ctx, id)
	if err == nil {
		return food, nil
	}

	// Try user foods
	food, err = fs.userFoodRepo.GetUserFoodByID(ctx, id)
	if err == nil {
		return food, nil
	}

	return models.Food{}, fmt.Errorf("food not found with ID: %s", id)
}

// GetAllFoods returns all foods from both embedded database and user foods
func (fs *FoodService) GetAllFoods(ctx context.Context) ([]models.Food, error) {
	var allFoods []models.Food

	// Get embedded foods
	embeddedFoods, err := fs.embeddedDB.GetAllFoods(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get embedded foods: %v\n", err)
	} else {
		allFoods = append(allFoods, embeddedFoods...)
	}

	// Get user foods
	userFoods, err := fs.userFoodRepo.GetUserFoods(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get user foods: %v\n", err)
	} else {
		allFoods = append(allFoods, userFoods...)
	}

	// Sort by name for consistent ordering
	sort.Slice(allFoods, func(i, j int) bool {
		return strings.ToLower(allFoods[i].Name) < strings.ToLower(allFoods[j].Name)
	})

	return allFoods, nil
}

// GetFoodsByCategory returns foods from a specific category from both sources
func (fs *FoodService) GetFoodsByCategory(ctx context.Context, category string) ([]models.Food, error) {
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	var allFoods []models.Food

	// Get embedded foods by category
	embeddedFoods, err := fs.embeddedDB.GetFoodsByCategory(ctx, category)
	if err != nil {
		fmt.Printf("Warning: failed to get embedded foods by category: %v\n", err)
	} else {
		allFoods = append(allFoods, embeddedFoods...)
	}

	// Get all user foods and filter by category
	userFoods, err := fs.userFoodRepo.GetUserFoods(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get user foods: %v\n", err)
	} else {
		categoryLower := strings.ToLower(strings.TrimSpace(category))
		for _, food := range userFoods {
			if strings.ToLower(food.Category) == categoryLower {
				allFoods = append(allFoods, food)
			}
		}
	}

	// Sort by name
	sort.Slice(allFoods, func(i, j int) bool {
		return strings.ToLower(allFoods[i].Name) < strings.ToLower(allFoods[j].Name)
	})

	return allFoods, nil
}

// GetAllCategories returns all unique categories from both embedded database and user foods
func (fs *FoodService) GetAllCategories(ctx context.Context) ([]string, error) {
	categoryMap := make(map[string]bool)

	// Get embedded categories
	embeddedCategories, err := fs.embeddedDB.GetCategories(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get embedded categories: %v\n", err)
	} else {
		for _, category := range embeddedCategories {
			categoryMap[category] = true
		}
	}

	// Get user food categories
	userFoods, err := fs.userFoodRepo.GetUserFoods(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get user foods for categories: %v\n", err)
	} else {
		for _, food := range userFoods {
			if food.Category != "" {
				categoryMap[food.Category] = true
			}
		}
	}

	// Convert map to sorted slice
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories, nil
}

// GetEmbeddedFoods returns only foods from the embedded database
func (fs *FoodService) GetEmbeddedFoods(ctx context.Context) ([]models.Food, error) {
	return fs.embeddedDB.GetAllFoods(ctx)
}

// GetUserFoods returns only user-defined foods
func (fs *FoodService) GetUserFoods(ctx context.Context) ([]models.Food, error) {
	return fs.userFoodRepo.GetUserFoods(ctx)
}

// SaveUserFood saves a user-defined food
func (fs *FoodService) SaveUserFood(ctx context.Context, food models.Food) error {
	return fs.userFoodRepo.SaveFood(ctx, food)
}

// UpdateUserFood updates a user-defined food
func (fs *FoodService) UpdateUserFood(ctx context.Context, id string, food models.Food) error {
	return fs.userFoodRepo.UpdateFood(ctx, id, food)
}

// DeleteUserFood deletes a user-defined food
func (fs *FoodService) DeleteUserFood(ctx context.Context, id string) error {
	return fs.userFoodRepo.DeleteFood(ctx, id)
}

// InitializeDatabase loads the embedded database
func (fs *FoodService) InitializeDatabase(ctx context.Context) error {
	return fs.embeddedDB.LoadDatabase(ctx)
}

// sortSearchResults sorts search results by relevance
func (fs *FoodService) sortSearchResults(foods []models.Food, query string) {
	queryLower := strings.ToLower(query)
	
	sort.Slice(foods, func(i, j int) bool {
		foodI := foods[i]
		foodJ := foods[j]
		
		nameI := strings.ToLower(foodI.Name)
		nameJ := strings.ToLower(foodJ.Name)
		
		// Exact name matches first
		if nameI == queryLower && nameJ != queryLower {
			return true
		}
		if nameJ == queryLower && nameI != queryLower {
			return false
		}
		
		// Name starts with query
		startsI := strings.HasPrefix(nameI, queryLower)
		startsJ := strings.HasPrefix(nameJ, queryLower)
		if startsI && !startsJ {
			return true
		}
		if startsJ && !startsI {
			return false
		}
		
		// User-defined foods after embedded foods (for same relevance)
		if foodI.IsUserDefined != foodJ.IsUserDefined {
			return !foodI.IsUserDefined // embedded foods first
		}
		
		// Alphabetical order as final tiebreaker
		return nameI < nameJ
	})
}

// GetFoodStats returns statistics about the food database
func (fs *FoodService) GetFoodStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get embedded food count
	embeddedFoods, err := fs.embeddedDB.GetAllFoods(ctx)
	if err != nil {
		stats["embedded_foods_count"] = 0
		stats["embedded_foods_error"] = err.Error()
	} else {
		stats["embedded_foods_count"] = len(embeddedFoods)
	}
	
	// Get user food count
	userFoods, err := fs.userFoodRepo.GetUserFoods(ctx)
	if err != nil {
		stats["user_foods_count"] = 0
		stats["user_foods_error"] = err.Error()
	} else {
		stats["user_foods_count"] = len(userFoods)
	}
	
	// Get total categories
	categories, err := fs.GetAllCategories(ctx)
	if err != nil {
		stats["categories_count"] = 0
		stats["categories_error"] = err.Error()
	} else {
		stats["categories_count"] = len(categories)
		stats["categories"] = categories
	}
	
	return stats, nil
}