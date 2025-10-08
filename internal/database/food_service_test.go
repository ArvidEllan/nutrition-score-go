package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nutritional-score/pkg/models"
)

func TestFoodService_SearchAllFoods(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create embedded database
	embeddedDBPath := filepath.Join(tempDir, "embedded_foods.json")
	embeddedData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test embedded database",
		"foods": [
			{
				"id": "embedded-apple-001",
				"name": "Apple, red",
				"category": "Fruits",
				"brand": "",
				"nutritional_data": {
					"energy": 218,
					"sugars": 10.4,
					"saturated_fatty_acids": 0.1,
					"sodium": 1,
					"fruits": 100,
					"fibre": 2.4,
					"protein": 0.3
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "USDA"
			}
		]
	}`
	
	err := os.WriteFile(embeddedDBPath, []byte(embeddedData), 0644)
	if err != nil {
		t.Fatalf("Failed to create embedded database file: %v", err)
	}
	
	// Create user foods repository
	userFoodsPath := filepath.Join(tempDir, "user_foods.json")
	
	// Initialize services
	embeddedDB := NewEmbeddedFoodDatabase(embeddedDBPath)
	userRepo := NewJSONUserFoodRepository(userFoodsPath)
	foodService := NewFoodService(embeddedDB, userRepo)
	
	ctx := context.Background()
	
	// Load embedded database
	err = foodService.InitializeDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Add a user food
	userFood := models.Food{
		Name:     "My Custom Apple",
		Category: "Fruits",
		NutritionalData: models.NutritionalData{
			Energy: 200,
			Sugars: 8.0,
		},
	}
	
	err = foodService.SaveUserFood(ctx, userFood)
	if err != nil {
		t.Fatalf("Failed to save user food: %v", err)
	}
	
	// Test search across both sources
	results, err := foodService.SearchAllFoods(ctx, "apple")
	if err != nil {
		t.Fatalf("Failed to search all foods: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'apple', got %d", len(results))
	}
	
	// Verify that embedded foods come first (due to sorting)
	embeddedFirst := false
	userFirst := false
	
	for _, food := range results {
		if !food.IsUserDefined {
			embeddedFirst = true
			break
		}
	}
	
	for _, food := range results {
		if food.IsUserDefined {
			userFirst = true
			break
		}
	}
	
	if !embeddedFirst || !userFirst {
		t.Error("Search should return both embedded and user foods")
	}
}

func TestFoodService_GetFoodByID(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create embedded database
	embeddedDBPath := filepath.Join(tempDir, "embedded_foods.json")
	embeddedData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test embedded database",
		"foods": [
			{
				"id": "embedded-apple-001",
				"name": "Apple, red",
				"category": "Fruits",
				"brand": "",
				"nutritional_data": {
					"energy": 218,
					"sugars": 10.4,
					"saturated_fatty_acids": 0.1,
					"sodium": 1,
					"fruits": 100,
					"fibre": 2.4,
					"protein": 0.3
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "USDA"
			}
		]
	}`
	
	err := os.WriteFile(embeddedDBPath, []byte(embeddedData), 0644)
	if err != nil {
		t.Fatalf("Failed to create embedded database file: %v", err)
	}
	
	// Create user foods repository
	userFoodsPath := filepath.Join(tempDir, "user_foods.json")
	
	// Initialize services
	embeddedDB := NewEmbeddedFoodDatabase(embeddedDBPath)
	userRepo := NewJSONUserFoodRepository(userFoodsPath)
	foodService := NewFoodService(embeddedDB, userRepo)
	
	ctx := context.Background()
	
	// Load embedded database
	err = foodService.InitializeDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Test getting embedded food by ID
	food, err := foodService.GetFoodByID(ctx, "embedded-apple-001")
	if err != nil {
		t.Fatalf("Failed to get embedded food by ID: %v", err)
	}
	
	if food.Name != "Apple, red" {
		t.Errorf("Expected 'Apple, red', got '%s'", food.Name)
	}
	
	if food.IsUserDefined {
		t.Error("Embedded food should not be marked as user-defined")
	}
	
	// Add a user food and test getting it by ID
	userFood := models.Food{
		Name:     "My Custom Apple",
		Category: "Fruits",
	}
	
	err = foodService.SaveUserFood(ctx, userFood)
	if err != nil {
		t.Fatalf("Failed to save user food: %v", err)
	}
	
	// Get user foods to find the ID
	userFoods, err := foodService.GetUserFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get user foods: %v", err)
	}
	
	if len(userFoods) != 1 {
		t.Fatalf("Expected 1 user food, got %d", len(userFoods))
	}
	
	userFoodID := userFoods[0].ID
	
	// Test getting user food by ID
	retrievedUserFood, err := foodService.GetFoodByID(ctx, userFoodID)
	if err != nil {
		t.Fatalf("Failed to get user food by ID: %v", err)
	}
	
	if retrievedUserFood.Name != "My Custom Apple" {
		t.Errorf("Expected 'My Custom Apple', got '%s'", retrievedUserFood.Name)
	}
	
	if !retrievedUserFood.IsUserDefined {
		t.Error("User food should be marked as user-defined")
	}
}

func TestFoodService_GetAllCategories(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create embedded database with multiple categories
	embeddedDBPath := filepath.Join(tempDir, "embedded_foods.json")
	embeddedData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test embedded database",
		"foods": [
			{
				"id": "embedded-apple-001",
				"name": "Apple",
				"category": "Fruits",
				"brand": "",
				"nutritional_data": {
					"energy": 218,
					"sugars": 10.4,
					"saturated_fatty_acids": 0.1,
					"sodium": 1,
					"fruits": 100,
					"fibre": 2.4,
					"protein": 0.3
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "USDA"
			},
			{
				"id": "embedded-chicken-001",
				"name": "Chicken",
				"category": "Meat",
				"brand": "",
				"nutritional_data": {
					"energy": 540,
					"sugars": 0.0,
					"saturated_fatty_acids": 1.0,
					"sodium": 74,
					"fruits": 0,
					"fibre": 0.0,
					"protein": 23.1
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "USDA"
			}
		]
	}`
	
	err := os.WriteFile(embeddedDBPath, []byte(embeddedData), 0644)
	if err != nil {
		t.Fatalf("Failed to create embedded database file: %v", err)
	}
	
	// Create user foods repository
	userFoodsPath := filepath.Join(tempDir, "user_foods.json")
	
	// Initialize services
	embeddedDB := NewEmbeddedFoodDatabase(embeddedDBPath)
	userRepo := NewJSONUserFoodRepository(userFoodsPath)
	foodService := NewFoodService(embeddedDB, userRepo)
	
	ctx := context.Background()
	
	// Load embedded database
	err = foodService.InitializeDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Add a user food with a new category
	userFood := models.Food{
		Name:     "My Custom Pasta",
		Category: "Grains",
	}
	
	err = foodService.SaveUserFood(ctx, userFood)
	if err != nil {
		t.Fatalf("Failed to save user food: %v", err)
	}
	
	// Get all categories
	categories, err := foodService.GetAllCategories(ctx)
	if err != nil {
		t.Fatalf("Failed to get all categories: %v", err)
	}
	
	if len(categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(categories))
	}
	
	// Check that all expected categories are present
	categoryMap := make(map[string]bool)
	for _, category := range categories {
		categoryMap[category] = true
	}
	
	expectedCategories := []string{"Fruits", "Meat", "Grains"}
	for _, expected := range expectedCategories {
		if !categoryMap[expected] {
			t.Errorf("Expected category '%s' not found", expected)
		}
	}
}

func TestFoodService_GetFoodStats(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create embedded database
	embeddedDBPath := filepath.Join(tempDir, "embedded_foods.json")
	embeddedData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test embedded database",
		"foods": [
			{
				"id": "embedded-apple-001",
				"name": "Apple",
				"category": "Fruits",
				"brand": "",
				"nutritional_data": {
					"energy": 218,
					"sugars": 10.4,
					"saturated_fatty_acids": 0.1,
					"sodium": 1,
					"fruits": 100,
					"fibre": 2.4,
					"protein": 0.3
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "USDA"
			}
		]
	}`
	
	err := os.WriteFile(embeddedDBPath, []byte(embeddedData), 0644)
	if err != nil {
		t.Fatalf("Failed to create embedded database file: %v", err)
	}
	
	// Create user foods repository
	userFoodsPath := filepath.Join(tempDir, "user_foods.json")
	
	// Initialize services
	embeddedDB := NewEmbeddedFoodDatabase(embeddedDBPath)
	userRepo := NewJSONUserFoodRepository(userFoodsPath)
	foodService := NewFoodService(embeddedDB, userRepo)
	
	ctx := context.Background()
	
	// Load embedded database
	err = foodService.InitializeDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Add user foods
	userFood1 := models.Food{Name: "User Food 1", Category: "Fruits"}
	userFood2 := models.Food{Name: "User Food 2", Category: "Vegetables"}
	
	err = foodService.SaveUserFood(ctx, userFood1)
	if err != nil {
		t.Fatalf("Failed to save user food 1: %v", err)
	}
	
	err = foodService.SaveUserFood(ctx, userFood2)
	if err != nil {
		t.Fatalf("Failed to save user food 2: %v", err)
	}
	
	// Get stats
	stats, err := foodService.GetFoodStats(ctx)
	if err != nil {
		t.Fatalf("Failed to get food stats: %v", err)
	}
	
	// Check embedded foods count
	if embeddedCount, ok := stats["embedded_foods_count"].(int); !ok || embeddedCount != 1 {
		t.Errorf("Expected embedded_foods_count to be 1, got %v", stats["embedded_foods_count"])
	}
	
	// Check user foods count
	if userCount, ok := stats["user_foods_count"].(int); !ok || userCount != 2 {
		t.Errorf("Expected user_foods_count to be 2, got %v", stats["user_foods_count"])
	}
	
	// Check categories count (should be 2: Fruits and Vegetables)
	if categoriesCount, ok := stats["categories_count"].(int); !ok || categoriesCount != 2 {
		t.Errorf("Expected categories_count to be 2, got %v", stats["categories_count"])
	}
}