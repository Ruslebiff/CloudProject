package cravings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	ingredientsList := strings.Split(r.URL.Query().Get("ingredients"), "_") //Array of ingredients
	recipeList, err := DBReadAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(recipeList)
	recipeCount := []RecipePrint{}

	for _, r := range recipeList { //Goes through all recipes
		recipeTemp := RecipePrint{}
		recipeTemp.RecipeName = r.RecipeName
		//recipeTemp.Ingredients.Remaining =
		for _, i := range r.Ingredients { //i is name of the ingredients needed for the recipe
			have := false //if there are enough of the ingredient to the recipe, this becomes true

			fmt.Println(i.Name + ":")

			for _, j := range ingredientsList { //Name|quantity of ingredients from query
				ingredient := strings.Split(j, "|")
				var quantity int
				if len(ingredient) < 2 { //checks if quantity is set for this ingredient
					quantity = 1 //Sets quantity to 'default' if not defined
				} else {
					quantity, err = strconv.Atoi(ingredient[1])
					if err != nil { //if error set to 1
						quantity = 1
					}
				}
				fmt.Println("\t"+ingredient[0], quantity, ":")
				if i.Name == ingredient[0] { //if it matches ingredient from recipe

					if quantity >= i.Quantity { //if there are more of an ingredient than needed
						have = true //have enough ingredients
						fmt.Println("\t\tHave: ")
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)

						rest := quantity - i.Quantity //quantity remaining after using recipe
						if rest > 0 {                 //if there is a rest,add it to the remaining list
							i.Quantity = rest
							recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining, i)
						}
					} else { //adds the ingredients that was sendt
						temp := i.Quantity
						i.Quantity = quantity //the quantity that user has
						recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
						i.Quantity = temp - quantity //sets quantity to the 'missing' value
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
