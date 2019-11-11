package cravings

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

// HandlerMeal is the handler for getting the recipes you can make out of your ingredients
func HandlerMeal(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "Content-Type", "application/json")

	ingredientsList := []Ingredient{}
	switch r.Method { //sets the list of remaining ingredients from either a post or get request
	case http.MethodPost:
		{
			err := json.NewDecoder(r.Body).Decode(&ingredientsList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	case http.MethodGet:
		{
			ingredientsList = ReadIngredients(strings.Split(r.URL.Query().Get("ingredients"), "_"))
		}
	}
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
					i = CalcNutrition(i)
					j = CalcNutrition(j) //adds nutritional value and makes the ingredient the same unit

					if j.Quantity <= i.Quantity { //If recipe needs more than what was sendt
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, j)                                                       //adds the ingredients sendt to 'have'
						recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining[:n], recipeTemp.Ingredients.Remaining[n+1:]...) //deletes the ingredient from remaining

						i.Quantity -= j.Quantity //sets the quantity to 'missing' value
						if i.Quantity > 0 {
							recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, CalcNutrition(i))
						}
					} else {

						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
						j.Quantity -= i.Quantity
						recipeTemp.Ingredients.Remaining[n] = CalcNutrition(j)
					}
					break //break out after finding matching name
				}
			}
		}
		recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
	}
	sort.Slice(recipeCount, func(i, j int) bool {
		return len(recipeCount[i].Ingredients.Missing) > len(recipeCount[j].Ingredients.Missing)
	})

	json.NewEncoder(w).Encode(recipeCount)
}
