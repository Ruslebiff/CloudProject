package cravings

import (
	"encoding/json"
	"net/http"
	"time"
)

// HandlerStatus handles the status endpoint
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	var S Status

	// Sets status for Edemam ***************************************
	resp, err := http.Get("https://api.edamam.com/api/nutrition-details") // gets api
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	S.Edemam = resp.StatusCode // sets status code for api

	// Sets staus for database ***************************************
	resp, err = http.Get("https://firebase.google.com") // gets link
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	S.Database = resp.StatusCode // sets status code for link

	// Sets total of recipes *****************************************
	statusRecipe, err := DBReadAllRecipes() // gets all recipes from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	S.TotalRecipe = len(statusRecipe) // sets totat for recipes

	// Sets status for Ingredients ***********************************
	statusIngredients, err := DBReadAllIngredients() // gets all ingredients from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	S.TotalIngredients = len(statusIngredients) // sers total for ingredients

	// Sets status for uptime ****************************************
	elapse := time.Since(StartTime) //sets run time
	S.Uptime = elapse.Seconds()     //convert run time to seconds

	S.Version = "v1"

	http.Header.Add(w.Header(), "Content-Type", "application/json") // makes the print look good

	json.NewEncoder(w).Encode(S)
}
