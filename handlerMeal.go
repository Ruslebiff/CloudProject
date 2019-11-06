package cravings

import (
	"fmt"
	"net/http"
	"strings"
)

func HandlerMeal(w http.ResponseWriter, r *http.Request) {

	ingredientsList := strings.Split(r.URL.Query().Get("ingredients"), "_") //Array of ingredients
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for r := range recipeList {
		for i := range ingredientsList {
			fmt.Fprintln(w, r, i)
		}
	}
}
