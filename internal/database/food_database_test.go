package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedFoodDatabase_LoadDatabase(t *testing.T) {
	// Create a temporary test database file
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_foods.json")
	
	// Create test data
	testData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test database",
		"foods": [
			{
				"id": "test-apple-001",
				"name": "Test Apple",
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
				"source": "Test"
			}
		]
	}`
	
	err := os.WriteFile(testDBPath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}
	
	// Test loading the database
	db := NewEmbeddedFoodDatabase(testDBPath)
	ctx := context.Background()
	
	err = db.LoadDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	
	if !db.IsLoaded() {
		t.Error("Database should be marked as loaded")
	}
	
	// Test getting all foods
	foods, err := db.GetAllFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get all foods: %v", err)
	}
	
	if len(foods) != 1 {
		t.Errorf("Expected 1 food, got %d", len(foods))
	}
	
	if foods[0].Name != "Test Apple" {
		t.Errorf("Expected food name 'Test Apple', got '%s'", foods[0].Name)
	}
}

func TestEmbeddedFoodDatabase_SearchFoods(t *testing.T) {
	// Create a temporary test database file with multiple foods
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_foods.json")
	
	testData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test database",
		"foods": [
			{
				"id": "test-apple-001",
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
				"source": "Test"
			},
			{
				"id": "test-banana-001",
				"name": "Banana, yellow",
				"category": "Fruits",
				"brand": "",
				"nutritional_data": {
					"energy": 371,
					"sugars": 12.2,
					"saturated_fatty_acids": 0.1,
					"sodium": 1,
					"fruits": 100,
					"fibre": 2.6,
					"protein": 1.1
				},
				"is_user_defined": false,
				"created_at": "2025-01-08T00:00:00Z",
				"updated_at": "2025-01-08T00:00:00Z",
				"source": "Test"
			}
		]
	}`
	
	err := os.WriteFile(testDBPath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}
	
	db := NewEmbeddedFoodDatabase(testDBPath)
	ctx := context.Background()
	
	err = db.LoadDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	
	// Test search by name
	results, err := db.SearchFoods(ctx, "apple")
	if err != nil {
		t.Fatalf("Failed to search foods: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'apple', got %d", len(results))
	}
	
	if results[0].Name != "Apple, red" {
		t.Errorf("Expected 'Apple, red', got '%s'", results[0].Name)
	}
	
	// Test search by category
	results, err = db.SearchFoods(ctx, "fruits")
	if err != nil {
		t.Fatalf("Failed to search foods by category: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'fruits', got %d", len(results))
	}
	
	// Test empty query
	_, err = db.SearchFoods(ctx, "")
	if err == nil {
		t.Error("Expected error for empty query")
	}
}

func TestEmbeddedFoodDatabase_GetFoodByID(t *testing.T) {
	// Create a temporary test database file
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_foods.json")
	
	testData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test database",
		"foods": [
			{
				"id": "test-apple-001",
				"name": "Test Apple",
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
				"source": "Test"
			}
		]
	}`
	
	err := os.WriteFile(testDBPath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}
	
	db := NewEmbeddedFoodDatabase(testDBPath)
	ctx := context.Background()
	
	err = db.LoadDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	
	// Test getting existing food
	food, err := db.GetFoodByID(ctx, "test-apple-001")
	if err != nil {
		t.Fatalf("Failed to get food by ID: %v", err)
	}
	
	if food.Name != "Test Apple" {
		t.Errorf("Expected 'Test Apple', got '%s'", food.Name)
	}
	
	// Test getting non-existent food
	_, err = db.GetFoodByID(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent food ID")
	}
	
	// Test empty ID
	_, err = db.GetFoodByID(ctx, "")
	if err == nil {
		t.Error("Expected error for empty food ID")
	}
}

func TestEmbeddedFoodDatabase_GetCategories(t *testing.T) {
	// Create a temporary test database file with multiple categories
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_foods.json")
	
	testData := `{
		"version": "1.0",
		"last_updated": "2025-01-08T00:00:00Z",
		"description": "Test database",
		"foods": [
			{
				"id": "test-apple-001",
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
				"source": "Test"
			},
			{
				"id": "test-chicken-001",
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
				"source": "Test"
			}
		]
	}`
	
	err := os.WriteFile(testDBPath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}
	
	db := NewEmbeddedFoodDatabase(testDBPath)
	ctx := context.Background()
	
	err = db.LoadDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	
	categories, err := db.GetCategories(ctx)
	if err != nil {
		t.Fatalf("Failed to get categories: %v", err)
	}
	
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
	
	// Check that both categories are present
	categoryMap := make(map[string]bool)
	for _, category := range categories {
		categoryMap[category] = true
	}
	
	if !categoryMap["Fruits"] {
		t.Error("Expected 'Fruits' category")
	}
	
	if !categoryMap["Meat"] {
		t.Error("Expected 'Meat' category")
	}
}