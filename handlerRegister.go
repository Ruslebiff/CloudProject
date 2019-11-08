package cravings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// HandlerRegister which registers either an ingredient or a recipe
func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[3]

	switch r.Method {
	// Gets either recipes or ingredients
	case http.MethodGet:
		switch endpoint {
		case "ingredient":
			ingredients := GetIngredient(w, r)
			totalIngredients := strconv.Itoa(len(ingredients))
			fmt.Fprintln(w, "Total ingredients: "+totalIngredients)
			json.NewEncoder(w).Encode(&ingredients)

		case "recipe":
			recipes := GetRecipe(w, r)
			totalRecipes := strconv.Itoa(len(recipes))
			fmt.Fprintln(w, "Total recipes: "+totalRecipes)
			json.NewEncoder(w).Encode(&recipes)
		}
		// Post either recipes or ingredients to firebase DB
	case http.MethodPost:
		authToken := Token{}
		resp, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Couldn't read request: ", http.StatusBadRequest)
		}

		err = json.Unmarshal(resp, &authToken)
		if err != nil {
			http.Error(w, "Unable to unmarshal request body: ", http.StatusBadRequest)
		}

		//  To post either one, you have to post it with a POST request with a .json body i.e. Postman
		//  and include the authorization token given by the developers through mail inside the body
		//  Detailed instructions for registering is in the readme
		if DBCheckAuthorization(authToken.AuthToken) {
			switch endpoint {
			case "ingredient": // Posts ingredient
				RegisterIngredient(w, resp)

			case "recipe": // Posts recipe
				RegisterRecipe(w, resp)
			}
		} else {
			http.Error(w, "Not authorized to POST to DB: ", http.StatusBadRequest)
			break
		}
	}
	w.Header().Add("content-type", "application/json")
}

// RegisterIngredient func saves the ingredient to its respective collection in our firestore DB
func RegisterIngredient(w http.ResponseWriter, respo []byte) {
	i := Ingredient{}
	err := json.Unmarshal(respo, &i)
	if err != nil {
		http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
	}

	err = DBSaveIngredient(&i)
	if err != nil {
		http.Error(w, "Could not save document to collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("content-type", "application/json")

	CallURL(IngredientCollection, &i)
}

// RegisterRecipe func saves the recipe to its respective collection in our firestore DB
func RegisterRecipe(w http.ResponseWriter, respo []byte) {
	rec := Recipe{}
	err := json.Unmarshal(respo, &rec)
	if err != nil {
		http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
	}

	err = DBSaveRecipe(&rec)
	if err != nil {
		http.Error(w, "Could not save document to collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("content-type", "application/json")

	CallURL(RecipeCollection, &rec)

}

// GetRecipe returns all recipes from database using the DBReadAllRecipes function
func GetRecipe(w http.ResponseWriter, r *http.Request) []Recipe {
	var allRecipes []Recipe
	allRecipes, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	return allRecipes
}

// GetIngredient returns all ingredients from database using the DBReadAllIngredients function
func GetIngredient(w http.ResponseWriter, r *http.Request) []Ingredient {
	var allIngredients []Ingredient
	allIngredients, err := DBReadAllIngredients()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	return allIngredients
}
