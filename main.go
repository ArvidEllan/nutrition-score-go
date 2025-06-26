package main

import (
	"fmt"

)
func main(){
	ns : GetNutritionalScore(NutritionalData{
		Energy: EnergyFromKcal(),
		Sugars: SugarGram(),
		SarturatedFattyAcids: SarturatedFattyAcids() ,
		Sodium: SodiumMilligram(),
		Fruits: FruitsPercent(),
		Fibre: FibreGram(),
		Protein: ProteinGram(),
	},Food ,
	)
	fmt.Printf("Nutritional Score:%d\n", ns.Value)
}