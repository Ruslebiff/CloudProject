package cravings

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// HandlerMeal is the handler for getting the recipes you can make out of your ingredients
func HandlerMeal(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "Content-Type", "application/json")

	ingredientsList := []Ingredient{}
	switch r.Method { //sets the list of remaining ingredients from either a post or get request
	case http.MethodPost:
		{ //  case Post posts a meal and ceodes it
			err := json.NewDecoder(r.Body).Decode(&ingredientsList)
			if err != nil {
				http.Error(w, "Failed to post meal "+err.Error(), http.StatusBadRequest)
				return
			}
		}
	case http.MethodGet:
		{ //  Case get reads the ingredients which is in the URL, each ingredient is separated by '_'
			ingredientsList = ReadIngredients(strings.Split(r.URL.Query().Get("ingredients"), "_"), w)
		}
	}
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, "Failed to retrieve recipes "+err.Error(), http.StatusInternalServerError)
		return
	}
	recipeCount := []RecipePrint{}

	for _, list := range recipeList { //Goes through all recipes
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = list.RecipeName
		recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining, ingredientsList...)

		for _, i := range list.Ingredients { //i is the ingredient needed for the recipe
			for n, j := range recipeTemp.Ingredients.Remaining { //Name|quantity of ingredients from query
				if j.Name == i.Name { //if it matches ingredient from recipe

					tempUnit := i.Unit //saves the unit the recipe is based on
					i = CalcNutrition(i, w)
					j = CalcNutrition(j, w) //calculates nutritional value
					ConvertUnit(&i, tempUnit)
					ConvertUnit(&j, tempUnit)     //sets both ingredients to the recipes unit
					if j.Quantity <= i.Quantity { //If recipe needs more than what was sendt

						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, j)                                                       //adds the ingredients sendt to 'have'
						recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining[:n], recipeTemp.Ingredients.Remaining[n+1:]...) //deletes the ingredient from remaining

						i.Quantity -= j.Quantity //calculates the 'missing' quantities
						if i.Quantity > 0 {
							CalcNutrition(i, w)       //calculate nutrition with new quantity
							ConvertUnit(&i, tempUnit) //set unit back to recipes unit
							recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
						}
					} else {

						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
						j.Quantity -= i.Quantity //calculates 'remaining' quantities
						CalcNutrition(j, w)
						ConvertUnit(&j, tempUnit)
						recipeTemp.Ingredients.Remaining[n] = j
					}
					break //break out after finding matching name
				}
			}
		}
		allowMissing, err := strconv.ParseBool(r.URL.Query().Get("allowMissing"))
		if err != nil {
			allowMissing = true //sets to true if not set or set to non-boolean in query
		}
		if allowMissing || len(recipeTemp.Ingredients.Missing) == 0 { //appends the recipe if it is allowed to be missing, or there are no ingredients missing
			recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
		}
	}
	sort.Slice(recipeCount, func(i, j int) bool {
		return len(recipeCount[i].Ingredients.Missing) > len(recipeCount[j].Ingredients.Missing)
	})

	json.NewEncoder(w).Encode(recipeCount)
}
