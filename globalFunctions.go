package cravings

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// DoRequest sends a new http request
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

// QueryGet func to read  variable for link
func QueryGet(s string, w http.ResponseWriter, r *http.Request) string {

	test := r.URL.Query().Get(s) // gets app key or app id
	if test == "" {              // check if it is empty
		fmt.Fprintln(w, s+" is missing")
	}
	return test

}

// CallURL post webhooks to webhooks.site
func CallURL(event string, s interface{}, w http.ResponseWriter) error {

	webhooks, err := DBReadAllWebhooks(w) // gets all webhooks
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for i := range webhooks { // loops true all webhooks
		if webhooks[i].Event == event { // see if webhooks.event is same as event
			var request = s

			requestBody, err := json.Marshal(request)
			if err != nil {
				fmt.Fprintln(w, "Can not encode: "+err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintln(w, "Attempting invoation of URL "+webhooks[i].URL+"...")

			resp, err := http.Post(webhooks[i].URL, "json", bytes.NewReader([]byte(requestBody))) // post webhook to webhooks.site
			if err != nil {
				fmt.Fprintln(w, "Error in HTTP request: "+err.Error(), http.StatusBadRequest)
			}

			defer resp.Body.Close() // close body

		}

	}

	return nil
}

// ReadIngredients splits up the ingredient name from the quantity from the URL
func ReadIngredients(ingredients []string, w http.ResponseWriter) []Ingredient {
	IngredientList := []Ingredient{}
	defVal := 1.0 //default value for quantity

	for i := range ingredients {
		ingredient := strings.Split(ingredients[i], "|") //splits up the string 'name|quantity|unit'
		var err error
		ingredientTemp := Ingredient{}
		ingredientTemp.Quantity = defVal //sets to defVal

		if len(ingredient) == 2 { //if quantity value is set
			ingredientTemp.Quantity, err = strconv.ParseFloat(ingredient[1], 64)

			if err != nil { //if error set to defVal
				ingredientTemp.Quantity = defVal
			}
		}

		if len(ingredient) == 3 { //if unit value is set
			ingredientTemp.Quantity, err = strconv.ParseFloat(ingredient[1], 64)

			if err != nil { //if error set to defVal
				ingredientTemp.Quantity = defVal
			}
			ingredientTemp.Unit = ingredient[2]
		}

		ingredientTemp.Name = ingredient[0] //name of the ingredient
		ingredientTemp = CalcNutrition(ingredientTemp, w)
		IngredientList = append(IngredientList, ingredientTemp)
	}
	return IngredientList
}

// CalcNutrition calculates nutritional info for given ingredient
func CalcNutrition(ing Ingredient, w http.ResponseWriter) Ingredient {
	temping, err := DBReadIngredientByName(ing.Name, w) //gets the ingredient with the same name from firebase
	if err != nil {
		fmt.Fprintln(w, "Could not read ingredient by name "+err.Error(), http.StatusBadRequest)
	}

	ing.ID = temping.ID               // add ID to ing since it's a copy
	ing.Nutrients = temping.Nutrients // reset nutrients to nutrients for 1g or 1l
	if ing.Unit == "kg" || ing.Unit == "g" {
		ConvertUnit(&ing, "g") // convert unit to g
	} else if ing.Unit == "l" || ing.Unit == "dl" || ing.Unit == "cl" || ing.Unit == "ml" {
		ConvertUnit(&ing, "l") // convert unit to g
	} else if ing.Unit == "pc" {
		// no conversion needed for pc
	} else if ing.Unit == "tablespoon" || ing.Unit == "teaspoon" {
		// check nutrition for it in API. No conversion needed
		err := GetNutrients(&ing, w)
		if err != nil {

		}
	}

	ing.Calories = temping.Calories * ing.Quantity //calculates calories based on ingredients quantity
	ing.Weight = temping.Weight * ing.Quantity     //calculates weight based on ingredients quantity

	// Calc nutrition :
	ing.Nutrients.Carbohydrate.Quantity *= ing.Quantity
	ing.Nutrients.Energy.Quantity *= ing.Quantity
	ing.Nutrients.Fat.Quantity *= ing.Quantity
	ing.Nutrients.Protein.Quantity *= ing.Quantity
	ing.Nutrients.Sugar.Quantity *= ing.Quantity

	return ing
}

// ConvertUnit converts units for ingredients, and changes their quantity respectively.
func ConvertUnit(ing *Ingredient, unitConvertTo string) {

	if ing.Unit == "kg" && unitConvertTo == "g" {
		ing.Quantity *= 1000
		ing.Unit = unitConvertTo
	}

	if ing.Unit == "g" && unitConvertTo == "kg" {
		ing.Quantity /= 1000
		ing.Unit = unitConvertTo
	}

	if unitConvertTo == "l" {

		switch ing.Unit {
		case "dl":
			ing.Quantity /= 10
		case "cl":
			ing.Quantity /= 100
		case "ml":
			ing.Quantity /= 1000
		}
		ing.Unit = unitConvertTo
	}
	if unitConvertTo == "dl" {

		switch ing.Unit {
		case "l":
			ing.Quantity *= 10
		case "cl":
			ing.Quantity /= 10
		case "ml":
			ing.Quantity /= 100
		}
		ing.Unit = unitConvertTo
	}
	if unitConvertTo == "cl" {

		switch ing.Unit {
		case "dl":
			ing.Quantity *= 10
		case "l":
			ing.Quantity *= 100
		case "ml":
			ing.Quantity /= 10
		}
		ing.Unit = unitConvertTo
	}
	if unitConvertTo == "ml" {

		switch ing.Unit {
		case "cl":
			ing.Quantity *= 10
		case "dl":
			ing.Quantity *= 100
		case "l":
			ing.Quantity *= 1000
		}
		ing.Unit = unitConvertTo
	}

}

// InitAPICredentials func opens up local file and reads the application id and key from that file
func InitAPICredentials() error {
	//  Opens local file which contains application id and key
	file, err := os.Open("appIdAndKey.txt")
	if err != nil {
		fmt.Println("Error: Unable to open file " + err.Error())
	}
	defer file.Close()
	//  Scans the lines of the file
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	App_id = scanner.Text()
	scanner.Scan()
	App_key = scanner.Text()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error: Unable to read the application ID and key from file " + err.Error())
	}
	return nil
}

// UnitCheck func checks the unit measurements of two ingredients and checks if they are of the same type solid/liquid
func UnitCheck(firstIngredient string, secondIngredient string) bool {
	if strings.Contains(firstIngredient, "l") {
		if strings.Contains(secondIngredient, "l") {
			return true
		}
	}
	if strings.Contains(firstIngredient, "g") {
		if strings.Contains(secondIngredient, "g") {
			return true
		}
	}
	if strings.Contains(firstIngredient, "spoon") {
		// table/teaspoon can be registered as liquid or solid
		return true
	}
	return false
}
