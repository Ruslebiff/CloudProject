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

	// Reads token from text file

	var testToken string
	file, err := os.Open("testToken.txt") // opens text file
	if err != nil {
		fmt.Println("Can't open file: ", err)
	}

	defer file.Close() // close file at the end

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		testToken = scanner.Text()
	}

	fmt.Println(testToken)

	// Testing methon POST for Ingredient ***********************************************
	i := testIngredient{Token: testToken, Name: "turmeric", Unit: "g"}
	req, _ := json.Marshal(i)
	reqTest := bytes.NewReader(req)
	r, err := http.NewRequest("POST", "/cravings/food/ingredient", reqTest)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing handlerFood POST method ingredient")

	// Testing methon POST for Recipe ***********************************************
	ingredient1 := Ingredient{Name: "milk", Quantity: 5, Unit: "ml"}
	ingredient2 := Ingredient{Name: "salt", Quantity: 2, Unit: "tablespoon"}
	ingredient3 := Ingredient{Name: "flour", Quantity: 1, Unit: "kg"}
	var testI []Ingredient
	testI = append(testI, ingredient1)
	testI = append(testI, ingredient2)
	testI = append(testI, ingredient3)

	re := testRecipe{Token: testToken, RecipeName: "TestRecipePOST", Ingredients: testI}
	req, _ = json.Marshal(re)
	reqTest = bytes.NewReader(req)
	r, err = http.NewRequest("POST", "/cravings/food/recipe", reqTest)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing handlerFood POST method recipe")

	// Testing method GET for all ingredients  ***********************************************
	r, err = http.NewRequest("GET", "/cravings/food/ingredient/", nil)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood GET method for all ingredient")

	// Testing method GET for one ingredient  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/ingredient/turmeric", nil)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood GET method for one ingredient")

	// Testing method GET for all recipes  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/recipe/", nil)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood GET method for all recipe")

	// Testing method GET for one recipe  ***********************************************

	r, err = http.NewRequest("GET", "/cravings/food/recipe/TestRecipePOST", nil)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood GET method for one recipe")

	// Testing method DELETE for ingredient ***********************************************
	test := testIngredient{Token: testToken, Name: "turmeric"}
	req, _ = json.Marshal(test)
	reqTest = bytes.NewReader(req)

	r, err = http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood DELETE method ingredient")

	// Testing method DELETE for ingredient **********************************************
	test2 := testRecipe{Token: testToken, RecipeName: "TestRecipePOST"}
	req, _ = json.Marshal(test2)
	reqTest = bytes.NewReader(req)

	r, err = http.NewRequest("DELETE", "/cravings/food/recipe", reqTest)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerFood)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing handlerFood DELETE method recipe")

}

func TestGetAllRecipes(t *testing.T) {

	r, err := http.NewRequest("GET", "/cravings/food/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRecipe := GetAllRecipes(w, r)
	if len(testRecipe) == 0 {
		t.Error("Cant read recipes from database")
	}
}

func TestGetAllIngredients(t *testing.T) {

	r, err := http.NewRequest("GET", "/cravings/food/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRecipe := GetAllIngredients(w, r)
	if len(testRecipe) == 0 {
		t.Error("Cant read ingredients from database")
	}

}
