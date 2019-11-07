package cravings

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerMeal is the handler for getting which recipes you can make out of your ingredients
func HandlerMeal(w http.ResponseWriter, r *http.Request) {

	ingredientsList := strings.Split(r.URL.Query().Get("ingredients"), "_") //Array of ingredients
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	recipeCount := []RecipePrint{}

	for _, r := range recipeList { //Goes through all recipes
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = r.RecipeName
		recipeTemp.Ingredients.DontHave = r.Ingredients

		for i := range r.Ingredients { //All indexes of ingredients needed for the recipe
			for _, j := range ingredientsList { //Name of ingredients from query
				have := false
				if r.Ingredients[i].Name == j { //if it matches ingredient from recipe
					recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, r.Ingredients[i])
					fmt.Println("have " + j)
					have = true
					break
				}
				if have {
					fmt.Println("dont " + r.Ingredients[i].Name)
					recipeTemp.Ingredients.DontHave = append(recipeTemp.Ingredients.DontHave, r.Ingredients[i])
				}
			}
		}
		recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
	}
	fmt.Fprintln(w, recipeCount)
}
