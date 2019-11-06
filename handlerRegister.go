package cravings

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Function which registers either an ingredient or a recipe
func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[3]
	//fmt.Println(endpoint)
	switch r.Method {
	// Gets either recipes or ingredients
	case http.MethodGet:
		switch endpoint {
		case "Ingredient":

		case "Recipe":
		}
		// Post either recipes or ingredients to firebase DB
	case http.MethodPost:
		switch endpoint {
		case "Ingredient": // Posts ingredient
			RegisterIngredient(w, r)

		case "Recipe": // Posts recipe
			RegisterRecipe(w, r)
		}
	}
}

// RegisterIngredient func saves the ingredient to its respective collection in our firestore DB
func RegisterIngredient(w http.ResponseWriter, r *http.Request) {
	i := Ingredient{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, "Could not decode body of request"+err.Error(), http.StatusBadRequest)
	}
	w.Header().Add("content-type", "application/json")

	err = DBSaveIngredient(&i)
	if err != nil {
		http.Error(w, "Could not save document to collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

}

// RegisterRecipe func saves the recipe to its respective collection in our firestore DB
func RegisterRecipe(w http.ResponseWriter, r *http.Request) {
	rec := Recipe{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, "Could not decode body of request"+err.Error(), http.StatusBadRequest)
	}
	w.Header().Add("content-type", "application/json")

	err = DBSaveRecipe(&rec)
	if err != nil {
		http.Error(w, "Could not save document to collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
}

func GetRecipe(w http.ResponseWriter, r *http.Request) {
	var allRecipes []Recipe
	allRecipes, err := DBReadAllRecipes
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
}

func GetIngredient(w http.ResponseWriter, r *http.Request) {
	var allIngredients []Ingredient
	allIngredients, err := DBReadAllIngredients
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}
}
