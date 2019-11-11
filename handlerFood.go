package cravings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// HandlerFood which registers or view either an ingredient or a recipe
// Whenever calling this endpoint in the browser, it is only possible to view the food,
// to register food, one has to post the .json body
func HandlerFood(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[3] // Store the query which represents either recipe or ingredient
	name := ""
	if len(parts) > 4 {
		name = parts[4] // The name of the ingredient or recipe
	}

	w.Header().Add("content-type", "application/json")

	switch r.Method {
	// Gets either recipes or ingredients
	case http.MethodGet:
		switch endpoint {
		case "ingredient":
			if name != "" { //  If user wrote in query for name of ingredient
				ingr := Ingredient{}
				ingr, err := DBReadIngredientByName(name) //  Get that ingredient
				if err != nil {
					http.Error(w, "Couldn't retrieve ingredient: "+err.Error(), http.StatusBadRequest)
				}
				json.NewEncoder(w).Encode(&ingr)
			} else { //  Else retireve all ingredients
				ingredients := GetAllIngredients(w, r)
				totalIngredients := strconv.Itoa(len(ingredients)) // With the number of total ingredients
				fmt.Fprintln(w, "Total ingredients: "+totalIngredients)
				json.NewEncoder(w).Encode(&ingredients)
			}
		case "recipe":
			if name != "" { //  If user wrote in query for name of recipe
				re := Recipe{}
				re, err := DBReadRecipeByName(name) //  Get that recipe
				if err != nil {
					http.Error(w, "Couldn't retrieve recipe: "+err.Error(), http.StatusBadRequest)
				}
				json.NewEncoder(w).Encode(&re)
			} else { //  Else get all recipes
				recipes := GetAllRecipes(w, r)
				totalRecipes := strconv.Itoa(len(recipes))
				fmt.Fprintln(w, "Total recipes: "+totalRecipes) // With the number of total recipes
				json.NewEncoder(w).Encode(&recipes)
			}
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

}

// RegisterIngredient func saves the ingredient to its respective collection in our firestore DB
func RegisterIngredient(w http.ResponseWriter, respo []byte) {
	ing := Ingredient{}
	found := false // ingredient found or not in database
	err := json.Unmarshal(respo, &ing)
	if err != nil {
		http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
	}
	ing.Quantity = 1
	ing.Name = strings.ToLower(ing.Name)

	if ing.Unit == "" {
		http.Error(w, "Could not save ingredient, missing \"unit\"", http.StatusBadRequest)
	} else {
		unitParam := ing.Unit
		if strings.Contains(unitParam, "g") {
			unitParam = "g"
		} else {
			unitParam = "l"
		}
		ConvertUnit(&ing, unitParam) // testing reference instead
		GetNutrients(&ing, w)        // calls func

		allIngredients, err := DBReadAllIngredients()
		if err != nil {
			http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
		}
		//  Check to see if the ingredient is already in the DB
		for i := range allIngredients {
			if ing.Name == allIngredients[i].Name {
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
				err := CallURL(IngredientCollection, &ing)
				if err != nil {
					fmt.Println("could not post to webhooks.site: ", err)
				}
				fmt.Fprintln(w, "Ingredient \""+ing.Name+"\" saved successfully to database.")
			}
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

	//  Retrieve all recipes and ingredients to see if the one the user is trying to register already exists
	allRecipes, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	//  Retrieves all the ingredients to get the ones missing for the recipe
	allIngredients, err := DBReadAllIngredients()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	//  If the name of the one created matches any of the ones in the DB
	for i := range allRecipes {
		if allRecipes[i].RecipeName == rec.RecipeName {
			recipeNameInUse = true
		}
	}

	for i := range rec.Ingredients { // Loops through all the ingredients
		found := false
		for _, j := range allIngredients { // If the ingredient is found the loop breaks and found is set to true
			if rec.Ingredients[i].Name == j.Name {
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
			err := CallURL(RecipeCollection, &rec)
			if err != nil {
				fmt.Println("could not post to webhooks.site: ", err)
			}
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
func GetAllRecipes(w http.ResponseWriter, r *http.Request) []Recipe {
	var allRecipes []Recipe
	allRecipes, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	return allRecipes
}

// GetIngredient returns all ingredients from database using the DBReadAllIngredients function
func GetAllIngredients(w http.ResponseWriter, r *http.Request) []Ingredient {
	var allIngredients []Ingredient
	allIngredients, err := DBReadAllIngredients()
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	return allIngredients
}

func GetNutrients(ing *Ingredient, w http.ResponseWriter) error {
	client := http.DefaultClient

	APIURL := "http://api.edamam.com/api/nutrition-data?app_id="
	APIURL += App_id
	APIURL += "&app_key="
	APIURL += App_key
	APIURL += "&ingr="
	if ing.Unit != "" {
		APIURL += ing.Unit
		APIURL += "%20"
	}
	APIURL += ing.Name
	r := DoRequest(APIURL, client, w)

	err := json.NewDecoder(r.Body).Decode(&ing)
	if err != nil {
		http.Error(w, "Could not decode response body "+err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// //  This is meant for when each ingredient is 100g, change later
// GetRecipeNutrients calculates total nutritients in a recipe
func GetRecipeNutrients(rec *Recipe, w http.ResponseWriter) error {
	//  Loops through each ingredient in the recipe and adds up the nutritional information from each
	//  to a total amount of nutrients for the recipe as a whol
	for i := range rec.Ingredients {
		temptotalnutrients := CalcNutrition(rec.Ingredients[i])

		rec.AllNutrients.Energy.Label = temptotalnutrients.Nutrients.Energy.Label
		rec.AllNutrients.Energy.Unit = temptotalnutrients.Nutrients.Energy.Unit
		rec.AllNutrients.Energy.Quantity += temptotalnutrients.Nutrients.Energy.Quantity

		rec.AllNutrients.Fat.Label = temptotalnutrients.Nutrients.Fat.Label
		rec.AllNutrients.Fat.Unit = temptotalnutrients.Nutrients.Fat.Unit
		rec.AllNutrients.Fat.Quantity += temptotalnutrients.Nutrients.Fat.Quantity

		rec.AllNutrients.Carbohydrate.Label = temptotalnutrients.Nutrients.Carbohydrate.Label
		rec.AllNutrients.Carbohydrate.Unit = temptotalnutrients.Nutrients.Carbohydrate.Unit
		rec.AllNutrients.Carbohydrate.Quantity += temptotalnutrients.Nutrients.Carbohydrate.Quantity

		rec.AllNutrients.Sugar.Label = temptotalnutrients.Nutrients.Sugar.Label
		rec.AllNutrients.Sugar.Unit = temptotalnutrients.Nutrients.Sugar.Unit
		rec.AllNutrients.Sugar.Quantity += temptotalnutrients.Nutrients.Sugar.Quantity

		rec.AllNutrients.Protein.Label = temptotalnutrients.Nutrients.Protein.Label
		rec.AllNutrients.Protein.Unit = temptotalnutrients.Nutrients.Protein.Unit
		rec.AllNutrients.Protein.Quantity += temptotalnutrients.Nutrients.Protein.Quantity

		rec.Ingredients[i].Nutrients.Energy = temptotalnutrients.Nutrients.Energy
		rec.Ingredients[i].Nutrients.Fat = temptotalnutrients.Nutrients.Fat
		rec.Ingredients[i].Nutrients.Carbohydrate = temptotalnutrients.Nutrients.Carbohydrate
		rec.Ingredients[i].Nutrients.Sugar = temptotalnutrients.Nutrients.Sugar
		rec.Ingredients[i].Nutrients.Protein = temptotalnutrients.Nutrients.Protein

		rec.Ingredients[i].Calories = temptotalnutrients.Nutrients.Energy.Quantity
		rec.Ingredients[i].ID = temptotalnutrients.ID
	}
	return nil
}
