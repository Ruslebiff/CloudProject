package cravings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// HandlerFood which registers or view either an ingredient or a recipe
// Whenever calling this endpoint in the browser, it is only possible to view the food,
// to register food, one has to post the .json body
func HandlerFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json") // JSON http header
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[3] // Store the query which represents either recipe or ingredient
	name := ""
	if len(parts) > 4 {
		name = parts[4] // The name of the ingredient or recipe
	}

	switch r.Method {
	case http.MethodGet: // Gets either recipes or ingredients
		switch endpoint {
		case "ingredient":
			if name != "" { //  If ingredient name is specified in URL
				ingr, err := DBReadIngredientByName(name, w) //  Get that ingredient
				if err != nil {
					http.Error(w, "Couldn't retrieve ingredient: "+err.Error(), http.StatusInternalServerError)
				}
				err = json.NewEncoder(w).Encode(&ingr)
				if err != nil {
					http.Error(w, "Couldn't encode response: "+err.Error(), http.StatusInternalServerError)
				}
			} else {
				ingredients, err := GetAllIngredients(w, r) //  Else retrieve all ingredients
				if err != nil {
					http.Error(w, "Couldn't retrieve ingredients: "+err.Error(), http.StatusBadRequest)
				}
				err = json.NewEncoder(w).Encode(&ingredients)
				if err != nil {
					http.Error(w, "Couldn't encode response: "+err.Error(), http.StatusInternalServerError)
				}
			}
		case "recipe":
			if name != "" { //  If user wrote in query for name of recipe
				re := Recipe{}
				re, err := DBReadRecipeByName(name, w) //  Get that recipe
				if err != nil {
					http.Error(w, "Couldn't retrieve recipe: "+err.Error(), http.StatusBadRequest)
				}

				err = json.NewEncoder(w).Encode(&re)
				if err != nil {
					http.Error(w, "Couldn't encode response: "+err.Error(), http.StatusInternalServerError)
				}
			} else {
				recipes, err := GetAllRecipes(w, r) //  Else get all recipes
				if err != nil {
					http.Error(w, "Couldn't retrieve recipes: "+err.Error(), http.StatusBadRequest)
				}
				err = json.NewEncoder(w).Encode(&recipes)
				if err != nil {
					http.Error(w, "Couldn't encode response: "+err.Error(), http.StatusInternalServerError)
				}
			}
		}

		// Post either recipes or ingredients to firebase DB
	case http.MethodPost: //  Func DBCheck checks if the user has posted a valid token, returns a bool
		authorised, resp := DBCheckAuthorization(w, r)

		//  To post either one, you have to post it with a POST request with a .json body i.e. Postman
		//  and include the authorization token given by the developers through mail inside the body
		//  Detailed instructions for registering is in the readme
		if authorised {
			switch endpoint {
			case "ingredient": // Posts ingredient
				RegisterIngredient(w, resp)

			case "recipe": // Posts recipe
				RegisterRecipe(w, resp)
			}
		} else {
			http.Error(w, "Not authorized to POST to DB: ", http.StatusBadRequest)
		}
	case http.MethodDelete:
		authorised, resp := DBCheckAuthorization(w, r)

		if authorised {
			switch endpoint {
			case "ingredient":

				ing := Ingredient{}
				err := json.Unmarshal(resp, &ing)
				if err != nil {
					http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
				}
				ing, err = DBReadIngredientByName(ing.Name, w) //  Get that ingredient
				if err != nil {
					http.Error(w, "Couldn't retrieve ingredient: "+err.Error(), http.StatusBadRequest)
				}
				err = DBDelete(ing.ID, IngredientCollection, w)
				if err != nil {
					http.Error(w, "Failed to delete ingredient: "+err.Error(), http.StatusInternalServerError)
				} else {
					fmt.Fprintln(w, "Successfully deleted ingredient", http.StatusOK)
				}

			case "recipe":

				rec := Recipe{}
				err := json.Unmarshal(resp, &rec)
				if err != nil {
					http.Error(w, "Could not unmarshal body of request"+err.Error(), http.StatusBadRequest)
				}
				rec, err = DBReadRecipeByName(rec.RecipeName, w) //  Get that recipe
				if err != nil {
					http.Error(w, "Couldn't retrieve recipe: "+err.Error(), http.StatusBadRequest)
				}

				err = DBDelete(rec.ID, RecipeCollection, w)
				if err != nil {
					http.Error(w, "Failed to delete recipe: "+err.Error(), http.StatusInternalServerError)
				} else {
					fmt.Fprintln(w, "Successfully deleted recipe"+rec.RecipeName, http.StatusOK)
				}

			}
		} else {
			fmt.Fprintln(w, "Not authorised to DELETE from DB:", http.StatusBadRequest)
		}
	default:
		fmt.Fprintln(w, "Invalid method "+r.Method, http.StatusBadRequest)
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
	ing.Name = strings.ToLower(ing.Name) // force lowercase ingredient name

	if ing.Unit == "" {
		http.Error(w, "Could not save ingredient, missing \"unit\"", http.StatusBadRequest)
	} else {
		unitParam := ing.Unit //  Checks if the posted unit is one of the legal measurements
		inList := false
		for _, v := range AllowedUnit { //  Loops through the allowed units
			if unitParam == v {
				inList = true
			}
		} //  If it is one of the allowed units, cast it into g or l
		if inList {
			if strings.Contains(unitParam, "g") {
				unitParam = "g"
			} else {
				unitParam = "l"
			}
		} else { //  Prints the allowed units for an ingridient
			http.Error(w, "Unit has to be of one of the values ", http.StatusBadRequest)
			for _, v := range AllowedUnit {
				fmt.Fprintln(w, v) // Print allowed units
			}
		}

		allIngredients, err := DBReadAllIngredients(w) // temporary list of all ingredients in database
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
		if !found { // if ingredient is not found in database
			ConvertUnit(&ing, unitParam) // convert unit to "g" or "l"
			ing.Quantity = 1             // force quantity to 1
			err = GetNutrients(&ing, w)  // get nutrients for the ingredient
			if err != nil {
				http.Error(w, "Couldn't get nutritional values: "+err.Error(), http.StatusInternalServerError)
			}

			if ing.Nutrients.Energy.Label == "" {
				// check if it got nutrients from db.
				//All ingredients will get this label if GetNutrients is ok
				http.Error(w, "ERROR: Failed to get nutrients for ingredient."+
					"Ingredient was not saved.", http.StatusInternalServerError)
			} else {
				err = DBSaveIngredient(&ing, w) // save it to database
				if err != nil {                 // if DBSaveIngredient return error
					http.Error(w, "Could not save document to collection "+
						IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
				} else { // DBSaveIngredient did not return error
					err := CallURL(IngredientCollection, &ing, w) // Call webhooks
					if err != nil {
						fmt.Fprintln(w, "Could not post to webhooks.site: "+
							err.Error(), http.StatusBadRequest)
					}
					fmt.Fprintln(w, "Ingredient \""+ing.Name+"\" saved successfully to database.") // Success!
				}
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
	allRecipes, err := DBReadAllRecipes(w)
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	//  Retrieves all the ingredients to get the ones missing for the recipe
	allIngredients, err := DBReadAllIngredients(w)
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	//  If the name of the one created matches any of the ones in the DB
	for i := range allRecipes {
		if allRecipes[i].RecipeName == rec.RecipeName {
			recipeNameInUse = true
		}
	}

	unitOk := false                  // Check to see if user has posted with the equivalent unit as the ingredient has in the DB
	for i := range rec.Ingredients { // Loops through all the ingredients
		found := false                     // Reset if current ingredient is found or not
		for _, j := range allIngredients { // If the ingredient is found the loop breaks and found is set to true
			if rec.Ingredients[i].Name == j.Name {
				ingredientsfound++
				found = true
				unitOk = UnitCheck(rec.Ingredients[i].Unit, j.Unit)
				if !unitOk {
					http.Error(w, rec.Ingredients[i].Name+" can't be saved with unit "+j.Unit, http.StatusBadRequest)
				}
				break
			}
		}
		if !found {
			missingingredients = append(missingingredients, rec.Ingredients[i].Name)
		}
	}

	//  If the ingredient found matchets that of the recipe, the name is available and the unit of legal value
	if ingredientsfound == recingredients && !recipeNameInUse && unitOk {
		err = GetRecipeNutrients(&rec, w) //  Collect the nutrients of that recipe
		if err != nil {
			http.Error(w, "Could not get nutrients for recipe", http.StatusInternalServerError)
		}
		err = DBSaveRecipe(&rec, w) //  Saves the recipe
		if err != nil {
			http.Error(w, "Could not save document to collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
		} else {
			err := CallURL(RecipeCollection, &rec, w) // Invokes the url
			if err != nil {
				fmt.Fprintln(w, "Could not post to webhooks.site: "+err.Error(), http.StatusBadRequest)
			}
			fmt.Fprintln(w, "Recipe \""+rec.RecipeName+"\" saved successfully to database.")
		}

	} else if ingredientsfound != recingredients {
		// console print:
		fmt.Fprintln(w, "Registration error: Recipe with name \""+rec.RecipeName+"\" is missing "+
			strconv.Itoa(recingredients-ingredientsfound)+" ingredient(s) "+err.Error(), http.StatusBadRequest)

		// http response:
		http.Error(w, "Cannot save recipe, missing ingredient(s) in database:", http.StatusBadRequest)
		for i := range missingingredients {
			// print all missing ingredients in http response
			fmt.Fprintln(w, "- "+missingingredients[i])
		}
		fmt.Fprintf(w, "\n Register these ingredients first!")
		//  If the name is in use
	} else if recipeNameInUse {
		// console print:
		fmt.Fprintln(w, "Registration error: Recipe with name \""+rec.RecipeName+"\" - name already in use. "+
			err.Error(), http.StatusBadRequest)

		// http response:
		http.Error(w, "Cannot save recipe, name already in use.", http.StatusBadRequest)
	} else if !unitOk {
		//  Error message when posting with mismatched units, i.e liquid with kg or solid with ml
		http.Error(w, "Couldn't save recipe due to unit mismatch", http.StatusBadRequest)
	} else {
		http.Error(w, "Cannot save recipe, internal server error.", http.StatusInternalServerError)
	}
}

// GetAllRecipes returns all recipes from database using the DBReadAllRecipes function
func GetAllRecipes(w http.ResponseWriter, r *http.Request) ([]Recipe, error) {
	var allRecipes []Recipe
	allRecipes, err := DBReadAllRecipes(w)
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	return allRecipes, err
}

// GetAllIngredients returns all ingredients from database using the DBReadAllIngredients function
func GetAllIngredients(w http.ResponseWriter, r *http.Request) ([]Ingredient, error) {
	var allIngredients []Ingredient
	allIngredients, err := DBReadAllIngredients(w)
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}

	return allIngredients, err
}

// GetNutrients gets nutritional info from external API for the ingredient. Returns http error if it fails
func GetNutrients(ing *Ingredient, w http.ResponseWriter) error {
	client := http.DefaultClient

	APIURL := "http://api.edamam.com/api/nutrition-data?app_id="
	APIURL += AppID
	APIURL += "&app_key="
	APIURL += AppKey
	APIURL += "&ingr="
	APIURL += strings.ReplaceAll(ing.Name, " ", "%20") // substitute spaces with "%20" so URL to API works with spaces in ingredient name
	if ing.Unit != "pc" {
		APIURL += "%20"
		APIURL += ing.Unit
	}
	resp, err := DoRequest(APIURL, client)
	if err != nil {
		http.Error(w, "Unable to get "+APIURL+err.Error(), http.StatusBadRequest)
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	err = json.NewDecoder(resp.Body).Decode(&ing)
	if err != nil {
		http.Error(w, "Could not decode response body "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// GetRecipeNutrients calculates total nutritients in a recipe
func GetRecipeNutrients(rec *Recipe, w http.ResponseWriter) error {
	// Set all the labels for the recipe
	rec.AllNutrients.Energy.Label = "Energy"
	rec.AllNutrients.Energy.Unit = "kcal"
	rec.AllNutrients.Fat.Label = "Fat"
	rec.AllNutrients.Fat.Unit = "g"
	rec.AllNutrients.Carbohydrate.Label = "Carbs"
	rec.AllNutrients.Carbohydrate.Unit = "g"
	rec.AllNutrients.Sugar.Label = "Sugar"
	rec.AllNutrients.Sugar.Unit = "g"
	rec.AllNutrients.Protein.Label = "Protein"
	rec.AllNutrients.Protein.Unit = "g"

	//  Loops through each ingredient in the recipe and adds up the nutritional information from each
	//  to a total amount of nutrients for the recipe as a whol
	for i := range rec.Ingredients {
		temptotalnutrients := CalcNutrition(rec.Ingredients[i], w)

		rec.AllNutrients.Energy.Quantity += temptotalnutrients.Nutrients.Energy.Quantity
		rec.AllNutrients.Fat.Quantity += temptotalnutrients.Nutrients.Fat.Quantity
		rec.AllNutrients.Carbohydrate.Quantity += temptotalnutrients.Nutrients.Carbohydrate.Quantity
		rec.AllNutrients.Sugar.Quantity += temptotalnutrients.Nutrients.Sugar.Quantity
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
