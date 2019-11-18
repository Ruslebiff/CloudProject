package cravings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// HandlerMeal is the handler for getting the recipes you can make out of your ingredients
func HandlerMeal(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "Content-Type", "application/json")

	ingredientsList := []Ingredient{}

	var err error

	switch r.Method { //sets the list of remaining ingredients from either a post or get request
	case http.MethodPost:
		{ //  case Post posts a meal and decodes it
			err := json.NewDecoder(r.Body).Decode(&ingredientsList)

			if err != nil {
				http.Error(w, "Failed to decode body "+err.Error(), http.StatusBadRequest)
				return
			}

			for _, i := range ingredientsList { //checks if unit is an allowed unit
				allowed := false

				for _, unit := range AllowedUnit {
					if i.Unit == unit {
						allowed = true
					}
				}

				if !allowed {
					http.Error(w, i.Unit+" is not an allowed unit.", http.StatusBadRequest)
					return
				}

				if i.Quantity <= 0 {
					http.Error(w, "Not a valid quantity.", http.StatusBadRequest)
					return
				}
			}
		}
	case http.MethodGet:
		{ //  Case get reads the ingredients which is in the URL query, each ingredient is separated by '_'
			ingredientsList, err = ReadIngredients(strings.Split(QueryGet("ingredients", "", r), "_"), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	recipeList, err := DBReadAllRecipes(w) //list of all recipes from firebase

	if err != nil {
		http.Error(w, "Failed to retrieve recipes "+err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range ingredientsList {
		ingredientsList[i], err = CalcNutrition(ingredientsList[i], w)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	//  Contains all the recipes the user can make with ingredients at hand,, including the ones which the user
	//  potentially could make
	recipeCount := []RecipePrint{}

	for _, list := range recipeList { //Goes through all recipes
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = list.RecipeName //  Appends the remaining ingredients to a list
		recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining, ingredientsList...)

		for _, i := range list.Ingredients { //i is the ingredient needed for the recipe
			found := false //sets found to true if ingredient is in recipe

			for n, j := range recipeTemp.Ingredients.Remaining { //Name|quantity of ingredients from query
				if j.Name == i.Name { //if it matches ingredient from recipe
					found = true       //found ingredient
					tempUnit := i.Unit //saves the unit the recipe is based on

					j = CalcRemaining(j, i, false) //calculates nutritional value for j

					if strings.Contains(i.Unit, "spoon") { //specialcase: if recipe uses tablespoon or teaspoon as unit
						noOfSpoons := j.Calories / (i.Calories / i.Quantity) //Amount we have/the value of calories from 1 spoon
						unitPerSpoon := j.Quantity / noOfSpoons              //calculates the amount of units stored per spoon

						if noOfSpoons <= i.Quantity { // if less or equal to what is needed from recipe
							tempOriginalUnit := j.Unit
							j.Unit = i.Unit         //set unit to recipes unit (...spoon)
							j.Quantity = noOfSpoons //Quantity to number of spoons
							recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, j)
							// Delete the ingredient from remaining:
							recipeTemp.Ingredients.Remaining =
								append(recipeTemp.Ingredients.Remaining[:n], recipeTemp.Ingredients.Remaining[n+1:]...)
							i.Quantity -= j.Quantity //  Calculates the amount the recipe needs after subtracting what we have

							if i.Quantity > 0 { //  If the recipe still needs more of the ingredient we have
								i.Unit = tempOriginalUnit
								i.Quantity *= unitPerSpoon     //total units for spoons
								i = CalcRemaining(i, j, false) //calculate nutrition with new quantity
								i.Unit = tempUnit
								i.Quantity /= unitPerSpoon //calculates back to spoon quantity
								recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
							}
						} else {
							recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
							j = CalcRemaining(j, i, true)
							recipeTemp.Ingredients.Remaining[n] = j
						}
					} else {
						ConvertUnit(&j, tempUnit) //sets both ingredients to the recipes unit

						if j.Quantity <= i.Quantity { //If recipe needs more than what was sent
							//adds the ingredients sent to 'have'
							recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, j)
							//deletes the ingredient from remaining:
							recipeTemp.Ingredients.Remaining =
								append(recipeTemp.Ingredients.Remaining[:n], recipeTemp.Ingredients.Remaining[n+1:]...)

							i.Quantity -= j.Quantity //calculates the 'missing' quantities

							if i.Quantity > 0 {
								i = CalcRemaining(i, j, false) //calculate nutrition with new quantity
								ConvertUnit(&i, tempUnit)      //set unit back to recipes unit
								recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
							}
						} else {
							recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
							j = CalcRemaining(j, i, true) //removes i's quantity from j and calculates the new nutrition value
							recipeTemp.Ingredients.Remaining[n] = j
						}
						break //break out after finding matching name
					}
				}
			}

			if !found { //adds the ingredient to 'missing' if not found
				recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
			}
		} //  Allow missing determines if we want to see the recipes we can make even though we're missing some ingredients

		allowMissing, err := strconv.ParseBool(r.URL.Query().Get("allowMissing")) //reads the allowMissing bool from query

		if err != nil {
			allowMissing = true //sets to true if not set or set to non-boolean
		}

		if allowMissing || len(recipeTemp.Ingredients.Missing) == 0 {
			//appends the recipe if it is allowed to be missing, or there are no ingredients missing
			recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
		}
	}

	sortBy := strings.ToLower(QueryGet("sortBy", "missing", r))

	switch sortBy {
	case "have": //  Sorts the recipes in an descending order of the most ingredients in "have" to least in the recipes
		sort.Slice(recipeCount, func(i, j int) bool {
			return len(recipeCount[i].Ingredients.Have) > len(recipeCount[j].Ingredients.Have)
		})
	case "remaining":
		//  Sorts the recipes in an ascending order of the least ingredients in "remaining" to most in the recipes
		sort.Slice(recipeCount, func(i, j int) bool {
			return len(recipeCount[i].Ingredients.Remaining) < len(recipeCount[j].Ingredients.Remaining)
		})
	default: //sorts by missing if not defined
		//  Sorts the recipes in an ascending order of the least ingredients in "missing" to most in the recipes
		sort.Slice(recipeCount, func(i, j int) bool {
			return len(recipeCount[i].Ingredients.Missing) < len(recipeCount[j].Ingredients.Missing)
		})
	}

	limit, err := strconv.Atoi(QueryGet("limit", "5", r)) //reads limit if sendt, else set it to 5

	if err != nil {
		limit = 5
	}

	if limit < len(recipeCount) { //if there are more than limit
		recipeCount = recipeCount[:limit] //sets recipecount to cut off all recipes after the value of limit
	}

	err = json.NewEncoder(w).Encode(recipeCount)

	if err != nil {
		http.Error(w, "Couldn't encode response: "+err.Error(), http.StatusBadRequest)
	}
}
