package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nutritional-score/pkg/models"
)

func TestJSONUserFoodRepository_SaveFood(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "user_foods.json")
	
	repo := NewJSONUserFoodRepository(testFilePath)
	ctx := context.Background()
	
	// Create a test food
	testFood := models.Food{
		Name:     "My Custom Apple",
		Category: "Fruits",
		Brand:    "Home Grown",
		NutritionalData: models.NutritionalData{
			Energy:              200,
			Sugars:              8.0,
			SaturatedFattyAcids: 0.1,
			Sodium:              1,
			Fruits:              100,
			Fibre:               3.0,
			Protein:             0.5,
		},
	}
	
	// Save the food
	err := repo.SaveFood(ctx, testFood)
	if err != nil {
		t.Fatalf("Failed to save food: %v", err)
	}
	
	// Verify the file was created
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Error("User foods file was not created")
	}
	
	// Get all user foods to verify it was saved
	foods, err := repo.GetUserFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get user foods: %v", err)
	}
	
	if len(foods) != 1 {
		t.Errorf("Expected 1 food, got %d", len(foods))
	}
	
	savedFood := foods[0]
	if savedFood.Name != testFood.Name {
		t.Errorf("Expected name '%s', got '%s'", testFood.Name, savedFood.Name)
	}
	
	if !savedFood.IsUserDefined {
		t.Error("Food should be marked as user-defined")
	}
	
	if savedFood.ID == "" {
		t.Error("Food should have been assigned an ID")
	}
	
	if savedFood.CreatedAt.IsZero() {
		t.Error("Food should have a creation timestamp")
	}
}

func TestJSONUserFoodRepository_UpdateFood(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "user_foods.json")
	
	repo := NewJSONUserFoodRepository(testFilePath)
	ctx := context.Background()
	
	// Create and save a test food
	testFood := models.Food{
		Name:     "Original Name",
		Category: "Fruits",
		NutritionalData: models.NutritionalData{
			Energy: 200,
			Sugars: 8.0,
		},
	}
	
	err := repo.SaveFood(ctx, testFood)
	if err != nil {
		t.Fatalf("Failed to save food: %v", err)
	}
	
	// Get the saved food to get its ID
	foods, err := repo.GetUserFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get user foods: %v", err)
	}
	
	savedFood := foods[0]
	originalCreatedAt := savedFood.CreatedAt
	
	// Update the food
	updatedFood := models.Food{
		Name:     "Updated Name",
		Category: "Vegetables",
		NutritionalData: models.NutritionalData{
			Energy: 150,
			Sugars: 5.0,
		},
	}
	
	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)
	
	err = repo.UpdateFood(ctx, savedFood.ID, updatedFood)
	if err != nil {
		t.Fatalf("Failed to update food: %v", err)
	}
	
	// Verify the update
	updatedFoodFromRepo, err := repo.GetUserFoodByID(ctx, savedFood.ID)
	if err != nil {
		t.Fatalf("Failed to get updated food: %v", err)
	}
	
	if updatedFoodFromRepo.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", updatedFoodFromRepo.Name)
	}
	
	if updatedFoodFromRepo.Category != "Vegetables" {
		t.Errorf("Expected category 'Vegetables', got '%s'", updatedFoodFromRepo.Category)
	}
	
	if updatedFoodFromRepo.CreatedAt != originalCreatedAt {
		t.Error("Created timestamp should be preserved during update")
	}
	
	if !updatedFoodFromRepo.UpdatedAt.After(originalCreatedAt) {
		t.Error("Updated timestamp should be newer than created timestamp")
	}
}

func TestJSONUserFoodRepository_DeleteFood(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "user_foods.json")
	
	repo := NewJSONUserFoodRepository(testFilePath)
	ctx := context.Background()
	
	// Create and save test foods
	testFood1 := models.Food{
		Name:     "Food 1",
		Category: "Fruits",
	}
	
	testFood2 := models.Food{
		Name:     "Food 2",
		Category: "Vegetables",
	}
	
	err := repo.SaveFood(ctx, testFood1)
	if err != nil {
		t.Fatalf("Failed to save food 1: %v", err)
	}
	
	err = repo.SaveFood(ctx, testFood2)
	if err != nil {
		t.Fatalf("Failed to save food 2: %v", err)
	}
	
	// Get the foods to get their IDs
	foods, err := repo.GetUserFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get user foods: %v", err)
	}
	
	if len(foods) != 2 {
		t.Fatalf("Expected 2 foods, got %d", len(foods))
	}
	
	// Delete the first food
	err = repo.DeleteFood(ctx, foods[0].ID)
	if err != nil {
		t.Fatalf("Failed to delete food: %v", err)
	}
	
	// Verify only one food remains
	remainingFoods, err := repo.GetUserFoods(ctx)
	if err != nil {
		t.Fatalf("Failed to get remaining foods: %v", err)
	}
	
	if len(remainingFoods) != 1 {
		t.Errorf("Expected 1 remaining food, got %d", len(remainingFoods))
	}
	
	// Verify the correct food was deleted
	if remainingFoods[0].ID == foods[0].ID {
		t.Error("Wrong food was deleted")
	}
	
	// Try to get the deleted food
	_, err = repo.GetUserFoodByID(ctx, foods[0].ID)
	if err == nil {
		t.Error("Expected error when getting deleted food")
	}
}

func TestJSONUserFoodRepository_SearchUserFoods(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "user_foods.json")
	
	repo := NewJSONUserFoodRepository(testFilePath)
	ctx := context.Background()
	
	// Create test foods
	foods := []models.Food{
		{
			Name:     "Red Apple",
			Category: "Fruits",
			Brand:    "Farm Fresh",
		},
		{
			Name:     "Green Apple",
			Category: "Fruits",
			Brand:    "Organic Co",
		},
		{
			Name:     "Banana Split",
			Category: "Desserts",
			Brand:    "Sweet Treats",
		},
	}
	
	// Save all foods
	for _, food := range foods {
		err := repo.SaveFood(ctx, food)
		if err != nil {
			t.Fatalf("Failed to save food '%s': %v", food.Name, err)
		}
	}
	
	// Test search by name
	results, err := repo.SearchUserFoods(ctx, "apple")
	if err != nil {
		t.Fatalf("Failed to search foods: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'apple', got %d", len(results))
	}
	
	// Test search by category
	results, err = repo.SearchUserFoods(ctx, "fruits")
	if err != nil {
		t.Fatalf("Failed to search by category: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'fruits', got %d", len(results))
	}
	
	// Test search by brand
	results, err = repo.SearchUserFoods(ctx, "organic")
	if err != nil {
		t.Fatalf("Failed to search by brand: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'organic', got %d", len(results))
	}
	
	if results[0].Name != "Green Apple" {
		t.Errorf("Expected 'Green Apple', got '%s'", results[0].Name)
	}
	
	// Test empty query
	_, err = repo.SearchUserFoods(ctx, "")
	if err == nil {
		t.Error("Expected error for empty query")
	}
	
	// Test no results
	results, err = repo.SearchUserFoods(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Failed to search for nonexistent food: %v", err)
	}
	
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}
}