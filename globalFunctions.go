package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func DoRequest(url string, c *http.Client, w http.ResponseWriter) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}

	resp, err := c.Do(req)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}

	return resp
}

//QueryGet func to read  variable for link
func QueryGet(s string, w http.ResponseWriter, r *http.Request) string {

	test := r.URL.Query().Get(s) // gets app key or app id
	if test == "" {              // check if it is empty
		fmt.Fprintln(w, s+" is missing")
	}
	return test

}

// CallURL post webhooks to webhooks.site
func CallURL(event string, s interface{}) {

	webhooks, err := DBReadAllWebhooks() // gets all webhooks
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for i := range webhooks { // loops true all webhooks
		if webhooks[i].Event == event { // see if webhooks.event is same as event
			var request = s

			requestBody, err := json.Marshal(request)
			if err != nil {
				fmt.Println("Can not encode: " + err.Error())
			}

			fmt.Println("Attempting invoation of URL " + webhooks[i].URL + "...")

			resp, err := http.Post(webhooks[i].URL, "json", bytes.NewReader([]byte(requestBody))) // post webhook to webhooks.site
			if err != nil {
				fmt.Println("Error in HTTP request: " + err.Error())
			}

			defer resp.Body.Close() // close body

		}

	}

}

//ReadIngredients splits up the ingredient name from the quantity
func ReadIngredients(ingredients []string) []Ingredient {
	IngredientList := []Ingredient{}
	for i := range ingredients {
		ingredient := strings.Split(ingredients[i], "|")
		var quantity int
		var err error
		if len(ingredient) < 2 { //checks if quantity is set for this ingredient
			quantity = 1.0 //Sets quantity to 'default' if not defined
		} else {
			quantity, err = strconv.Atoi(ingredient[1])
			if err != nil { //if error set to 1
				quantity = 1.0
			}
		}
		ingredientTemp := Ingredient{}
		ingredientTemp.Name = ingredient[0] //name of the ingredient
		ingredientTemp.Quantity = quantity  //quantity of the ingredient

		IngredientList = append(IngredientList, ingredientTemp)
	}
	return IngredientList
}

func RemoveIngredient(list []Ingredient, ingredient Ingredient) []Ingredient {
	for n, i := range list {
		if i.Name == ingredient.Name {
			fmt.Println(i.Quantity, " : ", ingredient.Quantity)
			if i.Quantity <= ingredient.Quantity {
				fmt.Println("Sletter: " + i.Name)

				list = append(list[:n], list[n+1:]...)
			} else {
				fmt.Println("Tar vekk: ", ingredient.Quantity, "fra: "+i.Name)
				i.Quantity = i.Quantity - ingredient.Quantity
			}
			return list
		}
	}
	return list
}
