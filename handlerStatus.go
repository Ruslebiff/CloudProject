package cravings

import (
	"encoding/json"
	"net/http"
	"time"
)

// HandlerStatus handles the status endpoint
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	var S Status

	APIURL := "https://api.edamam.com/api/nutrition-data?app_id="
	APIURL += App_id
	APIURL += "&app_key="
	APIURL += App_key
	APIURL += "&ingr=1%20large%20apple"

	// Sets status for Edamam ***************************************
	resp, err := http.Get(APIURL) // gets api
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	S.Edamam = resp.StatusCode // sets status code for api

	// Sets staus for database ***************************************
	resp, err = http.Get("https://firebase.google.com") // gets link
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	S.Database = resp.StatusCode // sets status code for link

	// Sets total of recipes *****************************************
	statusRecipe, err := DBReadAllRecipes(w) // gets all recipes from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+RecipeCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	S.TotalRecipe = len(statusRecipe) // sets totat for recipes

	// Sets status for Ingredients ***********************************
	statusIngredients, err := DBReadAllIngredients(w) // gets all ingredients from database
	if err != nil {
		http.Error(w, "Could not retrieve collection "+IngredientCollection+" "+err.Error(), http.StatusInternalServerError)
	}
	S.TotalIngredients = len(statusIngredients) // sers total for ingredients

	// Sets status for uptime ****************************************
	elapse := time.Since(StartTime) //sets run time
	S.Uptime = elapse.Seconds()     //convert run time to seconds

	S.Version = "v1"

	http.Header.Add(w.Header(), "Content-Type", "application/json") // makes the print look good

	err = json.NewEncoder(w).Encode(S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
