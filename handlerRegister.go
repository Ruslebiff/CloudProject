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
			ingredients := GetIngredient(w, r)                 // calls func
			totalIngredients := strconv.Itoa(len(ingredients)) // sets total ingredients
			fmt.Fprintln(w, "Total ingredients: "+totalIngredients)
			json.NewEncoder(w).Encode(&ingredients)

		case "recipe":
			recipes := GetRecipe(w, r)                 //calls func
			totalRecipes := strconv.Itoa(len(recipes)) // sets total recipes
			fmt.Fprintln(w, "Total recipes: "+totalRecipes)
			json.NewEncoder(w).Encode(&recipes)
		}
		// Post either recipes or ingredients to firebase DB
	case http.MethodPost:
		authToken := Token{}
		resp, err := ioutil.ReadAll(r.Body) // reads body from request
		if err != nil {
			http.Error(w, "Couldn't read request: ", http.StatusBadRequest)
		}

		err = json.Unmarshal(resp, &authToken) // unmarshal the requested input
		if err != nil {
			http.Error(w, "Unable to unmarshal request body: ", http.StatusBadRequest)
		}

		//  To post either one, you have to post it with a POST request with a .json body i.e. Postman
		//  and include the authorization token given by the developers through mail inside the body
		//  Detailed instructions for registering is in the readme
		if DBCheckAuthorization(authToken.AuthToken) { // check for Authorization
			switch endpoint {
			case "ingredient": // Posts ingredient
				RegisterIngredient(w, resp)

			case "recipe": // Posts recipe
				RegisterRecipe(w, resp)
			}
		} else { // there is no authorization
			http.Error(w, "Not authorized to POST to DB: ", http.StatusBadRequest)
			break
		}
	}
	w.Header().Add("content-type", "application/json") // givse json format to output on webpage
}

// RegisterIngredient func saves the ingredient to its respective collection in our firestore DB
func RegisterIngredient(w http.ResponseWriter, respo []byte) {
	ing := Ingredient{}
	found := false // ingredient found or not in database
	err := json.Unmarshal(respo, &ing)
	if err != nil {
		http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
	}
	GetNutrients(&ing, w) // calls func

	allIngredients, err := DBReadAllIngredients() // reads all ingredients from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	for i := range allIngredients { // loops true all ingredients
		if ing.Name == allIngredients[i].Name { // check if name exists in database
			found = true // found ingredient in database
			http.Error(w, "Ingredient \""+ing.Name+"\" already in database.", http.StatusBadRequest)
			break
		}
	}

	if found == false { // if ingredient is not found in database
		err = DBSaveIngredient(&ing) // save it
		if err != nil {
			http.Error(w, "Could not save document to collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
		} else {
			// if saving didn't return error, call webhooks
			CallURL(IngredientCollection, &ing) // post a webhook to webhooks.site with information on what has been added
			fmt.Fprintln(w, "Ingredient \""+ing.Name+"\" saved successfully to database.")
		}
	}

}

// RegisterRecipe func saves the recipe to its respective collection in our firestore DB
func RegisterRecipe(w http.ResponseWriter, respo []byte) {
	rec := Recipe{}
	err := json.Unmarshal(respo, &rec)
	if err != nil {
		http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
	}

	recingredients := len(rec.Ingredients) // number of ingredients in recipe
	ingredientsfound := 0                  // number of ingredients in recipe found in database
	var missingingredients []string        // name of ingredients in recipe missing in database
	recipeNameInUse := false

	allRecipes, err := DBReadAllRecipes() // reads all recipes from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	allIngredients, err := DBReadAllIngredients() // reads all ingredients from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	for i := range allRecipes { // loops true all recipes
		if allRecipes[i].RecipeName == rec.RecipeName { // checks if recipe name is in use
			recipeNameInUse = true
		}
	}

	for i := range rec.Ingredients { // loops true all ingredients in recipe
		found := false
		for _, j := range allIngredients { // loops true all ingredients
			if rec.Ingredients[i].Name == j.Name { // check if ingredients in recipe exsists in database
				ingredientsfound = ingredientsfound + 1
				found = true
				break
			}
		}
		if found == false {
			missingingredients = append(missingingredients, rec.Ingredients[i].Name)
		}
	}

	// difference for printing
	diff := strconv.Itoa(recingredients - ingredientsfound)

	if ingredientsfound == recingredients && recipeNameInUse == false {
		err = GetRecipeNutrients(&rec, w)
		if err != nil {
			http.Error(w, "Could not get nutrients for recipe", http.StatusInternalServerError)
		}
		err = DBSaveRecipe(&rec)
		if err != nil {
			http.Error(w, "Could not save document to collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
		} else {
			//CallURL(RecipeCollection, &rec)
			fmt.Fprintln(w, "Recipe \""+rec.RecipeName+"\" saved successfully to database.")
		}

	} else if ingredientsfound != recingredients {
		// console print:
		fmt.Println("Registration error: Recipe with name \"" + rec.RecipeName + "\" is missing " + diff + " ingredient(s)")

		// http response:
		http.Error(w, "Cannot save recipe, missing ingredient(s) in database:", http.StatusBadRequest)
		for i := range missingingredients {
			// print all missing ingredients in http response
			fmt.Fprintln(w, "- "+missingingredients[i])
		}
		fmt.Fprintf(w, "\n Register these ingredients first!")

	} else if recipeNameInUse == true {
		// console print:
		fmt.Println("Registration error: Recipe with name \"" + rec.RecipeName + "\" - name already in use.")

		// http response:
		http.Error(w, "Cannot save recipe, name already in use.", http.StatusBadRequest)
	} else {
		http.Error(w, "Cannot save recipe, internal server error.", http.StatusInternalServerError)
	}
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

func GetNutrients(ing *Ingredient, w http.ResponseWriter) { // fix error return?
	client := http.DefaultClient
	APIURL := "http://api.edamam.com/api/nutrition-data?app_id=f1d62971&app_key=fd32917955dc051f73436739d92b374e&ingr="
	//APIURL += strconv.Itoa(ing.Quantity) // temp removed due to changing Quantity to Float64 type
	APIURL += "%20"
	if ing.Unit != "" {
		APIURL += ing.Unit
		APIURL += "%20"
	}
	APIURL += ing.Name
	fmt.Println(APIURL)
	r := DoRequest(APIURL, client, w)

	err := json.NewDecoder(r.Body).Decode(&ing)
	if err != nil {
		http.Error(w, "Could not HER BAJSER JEG PAA MEE decode response body "+err.Error(), http.StatusInternalServerError)
	}

}

// GetRecipeNutrients calculates total nutritients in a recipe
func GetRecipeNutrients(rec *Recipe, w http.ResponseWriter) error {
	for i := range rec.Ingredients {
		temptotalnutrients := CalcNutrition(rec.Ingredients[i], rec.Ingredients[i].Unit, rec.Ingredients[i].Quantity)
		// assign these to rec totalnutrients or something
	}

	return nil
}
