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
func CallURL(event string, s interface{}) error {

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

	return nil
}

//ReadIngredients splits up the ingredient name from the quantity
func ReadIngredients(ingredients []string) []Ingredient {
	IngredientList := []Ingredient{}
	defVal := 1.0

	for i := range ingredients {
		ingredient := strings.Split(ingredients[i], "|")
		var err error
		ingredientTemp := Ingredient{}
		ingredientTemp.Quantity = defVal //sets to defVal

		if len(ingredient) == 2 {
			ingredientTemp.Quantity, err = strconv.ParseFloat(ingredient[1], 64)

			if err != nil { //if error set to defVal
				ingredientTemp.Quantity = defVal
			}
		}

		if len(ingredient) == 3 {
			ingredientTemp.Quantity, err = strconv.ParseFloat(ingredient[1], 64)

			if err != nil { //if error set to defVal
				ingredientTemp.Quantity = defVal
			}
			ingredientTemp.Unit = ingredient[2]
		}

		ingredientTemp.Name = ingredient[0] //name of the ingredient
		ingredientTemp = CalcNutrition(ingredientTemp, ingredientTemp.Unit, ingredientTemp.Quantity)
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

// CalcNutrition calculates nutritional info for given ingredient
func CalcNutrition(ing Ingredient, unit string, quantity float64) Ingredient { // maybe only get ingredient as parameter
	temping, err := DBReadIngredientByName(ing.Name)
	if err != nil {
		fmt.Println("Cound not read ingredient by name")
	}
	ing.Nutrients = temping.Nutrients // reset nutrients to 1g or 1l
	ing.ID = temping.ID               // add ID to ing since it's a copy

	//temping = ConvertUnit(&ing)
	ConvertUnit(&temping)
	ing.Unit = temping.Unit
	ing.Quantity = temping.Quantity

	ing.Nutrients.Energy.Quantity *= temping.Quantity
	ing.Nutrients.Fat.Quantity *= temping.Quantity
	ing.Nutrients.Carbohydrate.Quantity *= temping.Quantity
	ing.Nutrients.Protein.Quantity *= temping.Quantity
	ing.Nutrients.Sugar.Quantity *= temping.Quantity

	return ing
}

// ConvertUnit converts units for ingredients, and changes their quantity respectively.
func ConvertUnit(ing *Ingredient) {
	switch ing.Unit {
	case "dl":
		ing.Quantity = ing.Quantity / 10
		ing.Unit = "l"
	case "cl":
		ing.Quantity = ing.Quantity / 100
		ing.Unit = "l"
	case "ml":
		ing.Quantity = ing.Quantity / 1000
		ing.Unit = "l"
	case "kg":
		ing.Quantity = ing.Quantity * 1000
		ing.Unit = "g"
	}
}
