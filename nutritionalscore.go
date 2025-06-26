package main

import (
	"fmt"
	"os"
)

type Scoretype int

const (
	Food Scoretype = iota
	Beverage
	Water
	Cheese
)

type NutritionalScore struct {
	Value     int
	Positive  int
	Negative  int
	Scoretype Scoretype
}

type EnergyKJ float64

func (e EnergyKJ) GetPoints() int {
	// Placeholder logic
	if e > 1000 {
		return 5
	}
	return 1
}

type SugarGram float64

func (s SugarGram) GetPoints() int {
	if s > 10 {
		return 5
	}
	return 1
}

type SaturatedFattyAcids float64

func (s SaturatedFattyAcids) GetPoints() int {
	if s > 5 {
		return 5
	}
	return 1
}

type SodiumMilligram float64

func (s SodiumMilligram) GetPoints() int {
	if s > 500 {
		return 5
	}
	return 1
}

type FruitsPercent float64

func (f FruitsPercent) GetPoints() int {
	if f > 80 {
		return 5
	}
	return 1
}

type FibreGram float64

func (f FibreGram) GetPoints() int {
	if f > 5 {
		return 5
	}
	return 1
}

type ProteinGram float64

func (p ProteinGram) GetPoints() int {
	if p > 5 {
		return 5
	}
	return 1
}

type NutritionalData struct {
	Energy              EnergyKJ
	Sugars              SugarGram
	SaturatedFattyAcids SaturatedFattyAcids
	Sodium              SodiumMilligram
	Fruits              FruitsPercent
	Fibre               FibreGram
	Protein             ProteinGram
}

func GetNutNutritionalScore(n NutritionalData, st Scoretype) NutritionalScore {
	value := 0
	positive := 0
	negative := 0

	if st != Water {
		fruitPoints := n.Fruits.GetPoints()
		fibrePoints := n.Fibre.GetPoints()
		negative = n.Energy.GetPoints() + n.Sugars.GetPoints() + n.SaturatedFattyAcids.GetPoints() + n.Sodium.GetPoints()
		positive = fruitPoints + fibrePoints + n.Protein.GetPoints()
		value = negative - positive
	}

	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		Scoretype: st,
	}
}

func main() {
	// Simple CLI for demonstration
	var n NutritionalData
	var st int
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
	fmt.Println("Enter Scoretype (0:Food, 1:Beverage, 2:Water, 3:Cheese):")
	fmt.Scan(&st)
	if st < 0 || st > 3 {
		fmt.Println("Invalid Scoretype")
		os.Exit(1)
	}
	result := GetNutNutritionalScore(n, Scoretype(st))
	fmt.Printf("Nutritional Score: %+v\n", result)
}
