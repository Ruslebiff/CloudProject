package cravings

import (
	"fmt"
	"net/http"
	"strings"
)

func HandlerMeal(w http.ResponseWriter, r *http.Request) {

	ingredientsList := strings.Split(r.URL.Query().Get("ingredients"), "_") //Array of ingredients

	for i := range ingredientsList {
		fmt.Println(i)
	}
}
