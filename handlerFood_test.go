package cravings

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testIngredient struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
}

type testRecipe struct {
	Token       string       `json:"token"`
	RecipeName  string       `json:"recipeName"`
	Ingredients []Ingredient `json:"ingredients"`
}

func TestHandlerFood(t *testing.T) {

	w := httptest.NewRecorder() // creates ResponsRecorder for all tests

	// Reads token from text file

	var testToken string
	file, err := os.Open("testToken.txt") // opens text file
	if err != nil {
		fmt.Println("Can't open file: ", err)
	}

	defer file.Close() // close file at the end

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		testToken = scanner.Text() // sets testToken to be Token read from file
	}

	fmt.Println(testToken)

	// Testing methon POST for Ingredient ***********************************************

	i := testIngredient{Token: testToken, Name: "turmeric", Unit: "g"} // creates test ingredient
	req, _ := json.Marshal(i)
	reqTest := bytes.NewReader(req)                                         // convert struct to *Reader
	r, err := http.NewRequest("POST", "/cravings/food/ingredient", reqTest) // creates request with body
	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testing handlerFood POST method ingredient")

	// Testing methon POST for Recipe ***********************************************
	ingredient1 := Ingredient{Name: "milk", Quantity: 5, Unit: "ml"}
	ingredient2 := Ingredient{Name: "salt", Quantity: 2, Unit: "tablespoon"}
	ingredient3 := Ingredient{Name: "flour", Quantity: 1, Unit: "kg"}
	var testI []Ingredient
	testI = append(testI, ingredient1)
	testI = append(testI, ingredient2)
	testI = append(testI, ingredient3)

	re := testRecipe{Token: testToken, RecipeName: "TestRecipePOST", Ingredients: testI} // creates test recipe
	req, _ = json.Marshal(re)
	reqTest = bytes.NewReader(req)                                     // convert struct to *Reader
	r, err = http.NewRequest("POST", "/cravings/food/recipe", reqTest) // creates request with body
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testing handlerFood POST method recipe")

	// Testing method GET for all ingredients  ***********************************************
	r, err = http.NewRequest("GET", "/cravings/food/ingredient/", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for all ingredient")

	// Testing method GET for one ingredient  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/ingredient/turmeric", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for one ingredient")

	// Testing method GET for all recipes  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/recipe/", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for all recipe")

	// Testing method GET for one recipe  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/recipe/TestRecipePOST", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for one recipe")

	// Testing method DELETE for ingredient ***********************************************
	test := testIngredient{Token: testToken, Name: "turmeric"}
	req, _ = json.Marshal(test)
	reqTest = bytes.NewReader(req) // convert struct to *Reader

	r, err = http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest) // creates request with body
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood DELETE method ingredient")

	// Testing method DELETE for ingredient **********************************************
	test2 := testRecipe{Token: testToken, RecipeName: "TestRecipePOST"}
	req, _ = json.Marshal(test2)
	reqTest = bytes.NewReader(req) // convert struct to *Reader

	r, err = http.NewRequest("DELETE", "/cravings/food/recipe", reqTest) // creates request with body
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood DELETE method recipe")

}

func TestGetAllRecipes(t *testing.T) {

	r, err := http.NewRequest("GET", "/cravings/food/", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	testRecipe, err := GetAllRecipes(w, r) // test func with w and r
	if len(testRecipe) == 0 {              // check that it dident return an empty slice
		t.Error("Cant read recipes from database")
	}
	if err != nil {
		t.Error(err)
	}

}

func TestGetAllIngredients(t *testing.T) {

	r, err := http.NewRequest("GET", "/cravings/food/", nil) // creates request
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	testRecipe, err := GetAllIngredients(w, r) // test func with w and r
	if len(testRecipe) == 0 {                  // check that it dident return an empty slice
		t.Error("Cant read ingredients from database")
	}
	if err != nil {
		t.Error(err)
	}

}
