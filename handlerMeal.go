package cravings

import (
	"fmt"
	"net/http"
	"strconv"
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
	fmt.Println(ingredientsList)
	recipeCount := []RecipePrint{}

	for _, r := range recipeList { //Goes through all recipes
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = r.RecipeName

		for _, i := range r.Ingredients { //i is name of the ingredients needed for the recipe
			have := false //if there are enough of the ingredient to the recipe, this becomes true

			for _, j := range ingredientsList { //Name of ingredients from query
				ingredient := strings.Split(j, "|")
				var quantity int

				if len(ingredient) < 2 { //checks if quantity is set for this ingredient
					quantity = 1 //Sets quantity to 'default' if not defined
				} else {
					t, err := strconv.Atoi(ingredient[1])
					if err != nil { //if error set to 1
						quantity = 1
					} else { //set to given value
						quantity = t
					}
				}

				if i.Name == j { //if it matches ingredient from recipe

					if quantity >= i.Quantity { //if there are more of an ingredient than needed
						have = true                   //have enough ingredients
						rest := quantity - i.Quantity //quantity remaining after using recipe
						i.Quantity = quantity
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
						if rest != 0 { //if there is a rest,add it to the remaining list
							i.Quantity = rest
							recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining, i)
						}
					} else { //adds the ingredients that was sendt
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
						i.Quantity -= quantity //sets quantity to the 'missing' value
					}
					fmt.Println("have " + j)
					break
				}
			}
			if !have {
				fmt.Println("dont " + i.Name)
				recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
			}
		}
		recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
	}
	fmt.Println(w, recipeCount)
}
