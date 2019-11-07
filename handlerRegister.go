package cravings

import (
	"encoding/json"
	"fmt"
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
		var r2 *http.Request
		r2.Body = r.Body
		authToken := Token{}
		json.NewDecoder(r2.Body).Decode(&authToken)

		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa")
		//  Dettte funker ikke, !!!!! For at det skal funke kommenter ut DBCHECKAUTHORIZATION
		i := Ingredient{}
		json.NewDecoder(r.Body).Decode(&i)
		fmt.Println(i.Name)
		// //  To post either one, you have to post it with a POST request with a .json body i.e. Postman
		// //  and include the authorization token given by the developers through mail inside the body
		// //  Detailed instructions for registering is in the readme
		// if DBCheckAuthorization(authToken.AuthToken) {
		// 	switch endpoint {
		// 	case "ingredient": // Posts ingredient
		// 		RegisterIngredient(w, r)
		// 	case "recipe": // Posts recipe
		// 		RegisterRecipe(w, r)
		// 	}
		// } else {
		// 	http.Error(w, "Not authorized to POST to DB: ", http.StatusBadRequest)
		// 	break
		// }
	}
	w.Header().Add("content-type", "application/json")
}

// RegisterIngredient func saves the ingredient to its respective collection in our firestore DB
func RegisterIngredient(w http.ResponseWriter, r *http.Request) {
	i := Ingredient{}
	json.NewDecoder(r.Body).Decode(&i)
	fmt.Println(i.Name)
	// err := json.NewDecoder(r.Body).Decode(&i)
	// if err != nil {
	// 	http.Error(w, "Could heiheiheihnot save document to collection  "+err.Error(), http.StatusBadRequest)
	// }

	// err = DBSaveIngredient(&i)
	// if err != nil {
	// 	http.Error(w, "Could not save document to collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	// }

	w.Header().Add("content-type", "application/json")
}

// RegisterRecipe func saves the recipe to its respective collection in our firestore DB
func RegisterRecipe(w http.ResponseWriter, r *http.Request) {
	rec := Recipe{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, "Could not decode body of request"+err.Error(), http.StatusBadRequest)
	}

	err = DBSaveRecipe(&rec)
	if err != nil {
		http.Error(w, "Could not save document to collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("content-type", "application/json")
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
