package cravings

import (
	"context"
	"fmt"
	"log"
	firebase "firebase.google.com/go"
)

//  Struct for a recipe containting ingredients
type Recipe struct {
	RecipeID string `json:"id"`
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