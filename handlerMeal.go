package cravings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// HandlerMeal is the handler for getting the recipes you can make out of your ingredients
func HandlerMeal(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "Content-Type", "application/json")
	/*
		switch(r.Method){
			case http.MethodPost{

				err := json.NewDecoder(r.Body).Decode(&web)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}*/
	ingredientsQuery := strings.Split(r.URL.Query().Get("ingredients"), "_") //Array of all the ingredients|quantity
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(recipeList)
	recipeCount := []RecipePrint{}

	for _, r := range recipeList { //Goes through all recipes
		ingredientsList := ReadIngredients(ingredientsQuery)
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = r.RecipeName
		recipeTemp.Ingredients.Remaining = ingredientsList
		for _, i := range r.Ingredients { //i is name of the ingredients needed for the recipe
			have := false //if there are enough of the ingredient to the recipe, this becomes true

			fmt.Println(i.Name + ":")

			for _, j := range ingredientsList { //Name|quantity of ingredients from query

				fmt.Println("\t"+j.Name, j.Quantity, ":")
				if i.Name == j.Name { //if it matches ingredient from recipe
					fmt.Println("\tFør: ", recipeTemp.Ingredients.Remaining)
					recipeTemp.Ingredients.Remaining = RemoveIngredient(recipeTemp.Ingredients.Remaining, i) //removes from remaining
					fmt.Println("\tEtt: ", recipeTemp.Ingredients.Remaining)

					if j.Quantity >= i.Quantity { //if there are more of an ingredient than needed
						have = true //have enough ingredients
						fmt.Println("\t\tHave: ")
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)

						rest := j.Quantity - i.Quantity //quantity remaining after using recipe
						/*if rest > 0 {                   //if there is a rest,add it to the remaining list
								i.Quantity = rest
								recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining, i)
							}
						} else*/

						if rest <= 0 { //adds the ingredients that was sendt
							temp := i.Quantity
							i.Quantity = j.Quantity //the quantity that user has
							recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
							i.Quantity = temp - j.Quantity //sets quantity to the 'missing' value
						}
					}
					//fmt.Println("have " + ingredient[0])
					//break //break out of loop since the ingredient was found
				}
			}
			if !have { //if ingredient was not found, put it in the missing list
				fmt.Println("\t\tDont: " + i.Name)
				recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
			}
		}
		recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
	}
	json.NewEncoder(w).Encode(recipeCount)
}
