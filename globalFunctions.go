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

//ReadIngredients splits up the ingredient name from the quantity from the URL
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

// CalcNutrition calculates nutritional info for given ingredient
func CalcNutrition(ing Ingredient, unit string, quantity float64) Ingredient { //maybe only ingredient as parameter

	temping, err := DBReadIngredientByName(ing.Name)
	if err != nil {
		fmt.Println("Cound not read ingredient by name")
	}

	ing.ID = temping.ID               // add ID to ing since it's a copy
	ing.Nutrients = temping.Nutrients // reset nutrients to nutrients for 1g or 1l
	if ing.Unit == "kg" || ing.Unit == "g" {
		ConvertUnit(&ing, "g") // convert unit to g
	} else if ing.Unit == "l" || ing.Unit == "dl" || ing.Unit == "cl" || ing.Unit == "ml" {
		ConvertUnit(&ing, "l") // convert unit to g
	}

	// Calc nutrition:
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
		ing.Unit = unitConvertTo
		switch ing.Unit {
		case "dl":
			ing.Quantity /= 10
		case "cl":
			ing.Quantity /= 100
		case "ml":
			ing.Quantity /= 1000
		}
	}
	if unitConvertTo == "dl" {
		ing.Unit = unitConvertTo
		switch ing.Unit {
		case "l":
			ing.Quantity *= 10
		case "cl":
			ing.Quantity /= 10
		case "ml":
			ing.Quantity /= 100
		}
	}
	if unitConvertTo == "cl" {
		ing.Unit = unitConvertTo
		switch ing.Unit {
		case "dl":
			ing.Quantity *= 10
		case "l":
			ing.Quantity *= 100
		case "ml":
			ing.Quantity /= 10
		}
	}
	if unitConvertTo == "ml" {
		ing.Unit = unitConvertTo
		switch ing.Unit {
		case "cl":
			ing.Quantity *= 10
		case "dl":
			ing.Quantity *= 100
		case "l":
			ing.Quantity *= 1000
		}
	}

}

//  InitAPICredentials func opens up local file and reads the application id and key from that file
func InitAPICredentials() error {
	//  Opens local file which contains application id and key
	file, err := os.Open("appIdAndKey.txt")
	if err != nil {
		fmt.Println("Error: Unable to open file")
	}
	defer file.Close()
	//  Scans the lines of the file
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	App_id = scanner.Text()
	scanner.Scan()
	App_key = scanner.Text()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error: Unable to read the application ID and key from file ")
	}
	return nil
}
