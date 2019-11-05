package cravings

import (
	"context"

	"cloud.google.com/go/firestore"
)

//  Struct for a recipe containting ingredients
type Recpie struct {
	RecipeID    string       `json:"id"`
	Ingredients []Ingredient `json:"ingredients"`
}

//  Struct for an ingredient
type Ingredient struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	Name     string `json:"name"`
	Calories int    `json:"kcal"`
	Weight   int    `json:"weight"`
}

// FirestoreDatabase implements our Database access through Firestore
type FirestoreDatabase struct {
	Ctx    context.Context
	Client *firestore.Client
}
