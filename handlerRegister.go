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
		case "Ingredient":
			RegisterIngredient(w, r)

		case "Recipe":
			RegisterRecipe(w, r)
		}
	}
}

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

func GetRecipe() {

}

func GetIngredient() {

}

//	parts :=  Siste part er enteen ingredient eller recipe
// 	Kjør switch og kall på respektiv handler
