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

func TempToken() (s string) {
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

	return testToken
}

func TestHandlerFoodPostI(t *testing.T) {
	// Testing methon POST for Ingredient ***********************************************
	testToken := TempToken()

	if testToken == "" {
		t.Error("Token was not read from file")
	}

	i := TestIngredient{Token: testToken, Name: "turmeric", Unit: "g"} // creates test ingredient
	req, _ := json.Marshal(i)
	reqTest := bytes.NewReader(req)                                         // convert struct to *Reader
	r, err := http.NewRequest("POST", "/cravings/food/ingredient", reqTest) // creates request with body

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood POST method ingredient")
}

func TestHandlerFoodPostR(t *testing.T) {
	// Testing methon POST for Recipe ***********************************************
	testToken := TempToken()

	if testToken == "" {
		t.Error("Token was not read from file")
	}

	ingredient1 := Ingredient{Name: "milk", Quantity: 5, Unit: "ml"}
	ingredient2 := Ingredient{Name: "salt", Quantity: 2, Unit: "tablespoon"}
	ingredient3 := Ingredient{Name: "flour", Quantity: 1, Unit: "kg"}
	var testI []Ingredient
	testI = append(testI, ingredient1)
	testI = append(testI, ingredient2)
	testI = append(testI, ingredient3)

	re := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST", Ingredients: testI} // creates test recipe
	req, _ := json.Marshal(re)
	reqTest := bytes.NewReader(req)                                     // convert struct to *Reader
	r, err := http.NewRequest("POST", "/cravings/food/recipe", reqTest) // creates request with body

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood POST method recipe")
}

func TestHandlerFoodGetAI(t *testing.T) {
	// Testing method GET for all ingredients  ***********************************************
	r, err := http.NewRequest("GET", "/cravings/food/ingredient/", nil) // creates request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for all ingredient")
}

func TestHandlerFoodGetOI(t *testing.T) {
	// Testing method GET for one ingredient  ***********************************************
	r, err := http.NewRequest("GET", "/cravings/food/ingredient/turmeric", nil) // creates request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for one ingredient")
}

func TestHandlerFoodGetAR(t *testing.T) {
	// Testing method GET for all recipes  ***********************************************
	r, err := http.NewRequest("GET", "/cravings/food/recipe/", nil) // creates request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for all recipe")
}

func TestHandlerFoodGetOR(t *testing.T) {
	// Testing method GET for one recipe  ***********************************************
	r, err := http.NewRequest("GET", "/cravings/food/recipe/TestRecipePOST", nil) // creates request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood GET method for one recipe")
}

func TestHandlerFoodDeleteI(t *testing.T) {
	// Testing method DELETE for ingredient ***********************************************
	testToken := TempToken()

	if testToken == "" {
		t.Error("Token was not read from file")
	}

	test := TestIngredient{Token: testToken, Name: "turmeric"}
	req, _ := json.Marshal(test)
	reqTest := bytes.NewReader(req) // convert struct to *Reader

	r, err := http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest) // creates request with body

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing handlerFood DELETE method ingredient")
}

func TestHandlerFoodDeleteR(t *testing.T) {
	// Testing method DELETE for ingredient **********************************************
	testToken := TempToken()

	if testToken == "" {
		t.Error("Token was not read from file")
	}

	test2 := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST"}
	req, _ := json.Marshal(test2)
	reqTest := bytes.NewReader(req) // convert struct to *Reader

	r, err := http.NewRequest("DELETE", "/cravings/food/recipe", reqTest) // creates request with body

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

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

	if len(testRecipe) == 0 { // check that it dident return an empty slice
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

	if len(testRecipe) == 0 { // check that it dident return an empty slice
		t.Error("Cant read ingredients from database")
	}
	if err != nil {
		t.Error(err)
	}
}
