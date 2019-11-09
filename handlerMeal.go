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

		for _, i := range r.Ingredients { //i is the ingredient needed for the recipe
			//have := false //if there are enough of the ingredient to the recipe, this becomes true
			fmt.Println(i.Name + ":")

			for _, j := range ingredientsList { //Name|quantity of ingredients from query

				fmt.Println("\t"+j.Name, j.Quantity, ":")
				if i.Name == j.Name { //if it matches ingredient from recipe

					fmt.Println("\tFør: ", recipeTemp.Ingredients.Remaining)
					//recipeTemp = RemoveIngredient(recipeTemp, i) //removes from remaining
					for n := range recipeTemp.Ingredients.Remaining {

						if recipeTemp.Ingredients.Remaining[n].Name == i.Name {
							fmt.Println(recipeTemp.Ingredients.Remaining[n], " : ", i.Quantity)

							if recipeTemp.Ingredients.Remaining[n].Quantity <= i.Quantity { //If recipe needs more than what was sendt
								fmt.Println("Sletter: " + recipeTemp.Ingredients.Remaining[n].Name)

								recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, j) //adds the ingredients sendt to have
								i.Quantity -= j.Quantity                                             //sets the quantity to 'missing' value
								recipeTemp.Ingredients.Missing = append(recipeTemp.Ingredients.Missing, i)
								recipeTemp.Ingredients.Remaining = append(recipeTemp.Ingredients.Remaining[:n], recipeTemp.Ingredients.Remaining[n+1:]...) //deletes the ingredient from remaining

							} else {
								fmt.Println("Tar vekk: ", i.Quantity, "fra: "+recipeTemp.Ingredients.Remaining[n].Name)
								recipeTemp.Ingredients.Have = append(recipeTemp.Ingredients.Have, i)
								recipeTemp.Ingredients.Remaining[n].Quantity = recipeTemp.Ingredients.Remaining[n].Quantity - i.Quantity
							}
							break //break out of search
						}
					}
					fmt.Println("\tEtt: ", recipeTemp.Ingredients.Remaining)
					break //break out after finding matching name
				}
			}
		}
		recipeCount = append(recipeCount, recipeTemp) //adds recipeTemp in the recipeCount
	}
	json.NewEncoder(w).Encode(recipeCount)
}
