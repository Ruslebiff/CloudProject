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
		var quantity float64
		var err error
		if len(ingredient) < 2 { //checks if quantity is set for this ingredient
			quantity = 1.0 //Sets quantity to 'default' if not defined
		} else {
			quantity, err = strconv.ParseFloat(ingredient[1], 64)
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

// CalcNutrition calculates nutritional info for given ingredient
func CalcNutrition(ing Ingredient, unit string, quantity float64) Ingredient {
	var grams float64
	var litres float64

	if unit == "l" {
		litres += quantity
	}
	if unit == "dl" {
		litres += quantity / 10
	}
	if unit == "cl" {
		litres += quantity / 100
	}
	if unit == "ml" {
		litres += quantity / 1000
	}
	if unit == "g" {
		grams += quantity
	}
	if unit == "kg" {
		grams += quantity * 1000
	}

	if grams > 0 {
		ing.Unit = "g"
		ing.Nutrients.Energy.Quantity *= grams
		ing.Nutrients.Fat.Quantity *= grams
		ing.Nutrients.Carbohydrate.Quantity *= grams
		ing.Nutrients.Protein.Quantity *= grams
		ing.Nutrients.Sugar.Quantity *= grams
	} else if litres > 0 {
		ing.Unit = "l"
		ing.Nutrients.Energy.Quantity *= litres
		ing.Nutrients.Fat.Quantity *= litres
		ing.Nutrients.Carbohydrate.Quantity *= litres
		ing.Nutrients.Protein.Quantity *= litres
		ing.Nutrients.Sugar.Quantity *= litres
	}

	return ing
}
