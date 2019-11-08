package cravings

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

// Recipe Struct for a recipe containting ingredients used in firebase.go and register.go
type Recipe struct {
	ID           string       `json:"id"`
	RecipeName   string       `json:"recipeName"`
	Ingredients  []Ingredient `json:"ingredients"`
	AllNutrients TotalNutrients
}

//RecipePrint struct containing the ingredients the user has, needs and what remains after using the recipe
type RecipePrint struct {
	RecipeName  string `json:"recipeName"`
	Ingredients struct {
		Have      []Ingredient `json:"have"`      //Ingredients that fits the recipe
		Missing   []Ingredient `json:"missing"`   //Missing ingredients for recipe
		Remaining []Ingredient `json:"remaining"` //Remaining ingredients after using recipe
	} `json:"ingredients"`
}

// Ingredient Struct for an ingredient used in firebase.go and register.go
type Ingredient struct {
	ID        string         `json:"id"`
	Quantity  float64        `json:"quantity"`
	Unit      string         `json:"unit"`
	Name      string         `json:"name"`
	Calories  int            `json:"calories"`
	Weight    float64        `json:"totalWeight"`
	Nutrients TotalNutrients `json:"totalNutrients"`
}

// Webhook Struct for an webhook used in firebase.go and webhooks.go
type Webhook struct {
	ID    string    `json:"id"`
	Event string    `json:"event"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time"`
}

//Nutrient Struct for nutrient from Edamam
type Nutrient struct {
	Label    string  `json:"label"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

//TotalNutrients Struct that stores the different nutrients from Edamam
type TotalNutrients struct {
	Fat          Nutrient `json:"FAT"`
	Protein      Nutrient `json:"PROCNT"`
	Carbohydrate Nutrient `json:"CHOCDF"`
	Sugar        Nutrient `json:"SUGAR"`
	Energy       Nutrient `json:"ENERC_KCAL"`
}

// VI FÅR SE OM DETTE BLIR MED
//TotalDaily Struct that stores the % of the daily nutrition the recipe contains (?-)
// type TotalDaily struct {
// 	AllNutrients struct {
// 		Fat          Nutrient `json:"FAT"`
// 		Protein      Nutrient `json:"PROCNT"`
// 		Carbohydrate Nutrient `json:"CHOCDF"`
// 		Energy       Nutrient `json:"ENERC_KCAL"`
// 	} `json:"totalDaily"`
// }

type RecipeAnalysisPost struct {
	Title       string   `json:"title"`
	Ingredients []string `json:"ingr"`
}

// FirestoreDatabase implements our Database access through Firestore
type FirestoreDatabase struct {
	Ctx    context.Context
	Client *firestore.Client
}

type Token struct {
	ID        string `json:"id"`
	AuthToken string `json:"token"`
}

// Status struct for status endpoint
type Status struct {
	Edemam           int     `json:"edemam"`
	Database         int     `json:"database"`
	TotalRecipe      int     `json:"total recipes"`
	TotalIngredients int     `json:"total ingredients"`
	Uptime           float64 `json:"uptime"`
	Version          string  `json:"version"`
}
