package main

import (
	"fmt"
	"os"
)

// main function - entry point for the nutritional score calculator
// This is a temporary main that will be replaced with the enhanced application structure
func main() {
	// Simple CLI for demonstration - this will be enhanced in later tasks
	var n NutritionalData
	var st int
	
	fmt.Println("=== Nutritional Score Calculator ===")
	
	// Collect nutritional data from user input with clear prompts
	fmt.Println("Enter Energy (kJ):")
	fmt.Scan(&n.Energy)
	fmt.Println("Enter Sugars (g):")
	fmt.Scan(&n.Sugars)
	fmt.Println("Enter Saturated Fatty Acids (g):")
	fmt.Scan(&n.SaturatedFattyAcids)
	fmt.Println("Enter Sodium (mg):")
	fmt.Scan(&n.Sodium)
	fmt.Println("Enter Fruits (%):")
	fmt.Scan(&n.Fruits)
	fmt.Println("Enter Fibre (g):")
	fmt.Scan(&n.Fibre)
	fmt.Println("Enter Protein (g):")
	fmt.Scan(&n.Protein)
	
	// Get score type from user with validation
	fmt.Println("Enter Scoretype (0:Food, 1:Beverage, 2:Water, 3:Cheese):")
	fmt.Scan(&st)
	
	// Validate score type input range
	if st < 0 || st > 3 {
		fmt.Println("Invalid Scoretype")
		os.Exit(1)
	}
	
	// Calculate and display the nutritional score using the corrected function name
	result := GetNutritionalScore(n, ScoreType(st))
	fmt.Printf("Nutritional Score: %+v\n", result)
}