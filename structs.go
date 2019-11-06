package cravings

import (
	"context"

	"cloud.google.com/go/firestore"
)

//  Struct for a recipe containting ingredients used in firebase.go and register.go
type Recipe struct {
	ID          string       `json:"id"`
	Ingredients []Ingredient `json:"ingredients"`
}

//  Struct for an ingredient used in firebase.go and register.go
type Ingredient struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	Name     string `json:"name"`
	Calories int    `json:"kcal"`
	Weight   int    `json:"weight"`
}

type Nutrient struct {
	Label    string `json:"label"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

type TotalNutrients struct {
	AllNutrients struct{
		Fat Nutrient `json:"FAT"`
		Protein Nutrient `json:"PROCNT"`
		Sugar Nutrient `json:"SUGAR"`
		Energy Nutrient `json:"ENERC_KCAL"`
	} `json:"totalNutrients"`
}

type TotalDaily struct {
	AllNutrients struct{
		Fat Nutrient `json:"FAT"`
		Protein Nutrient `json:"PROCNT"`
		Energy Nutrient `json:"ENERC_KCAL"`
	} `json:"totalDaily"`
}

// FirestoreDatabase implements our Database access through Firestore
type FirestoreDatabase struct {
	Ctx    context.Context
	Client *firestore.Client
}
