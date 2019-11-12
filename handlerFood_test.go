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
	requestI, _ := json.Marshal(i)
	requestTestIngredient := bytes.NewReader(requestI)
	ri, err := http.NewRequest("POST", "/cravings/food/ingredient", requestTestIngredient)
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Error(err)
	}
	wi := httptest.NewRecorder()
	handlerIngredient := http.HandlerFunc(HandlerFood)
	handlerIngredient.ServeHTTP(wi, ri)

	respIngredient := wi.Result()

	if respIngredient.StatusCode != http.StatusOK {
		t.Error(respIngredient.StatusCode)
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
	requestR, _ := json.Marshal(re)
	requestTestRecipe := bytes.NewReader(requestR)
	rr, err := http.NewRequest("POST", "/cravings/food/recipe", requestTestRecipe)
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Error(err)
	}
	wr := httptest.NewRecorder()
	handlerRecipe := http.HandlerFunc(HandlerFood)
	handlerRecipe.ServeHTTP(wr, rr)

	respRecipe := wr.Result()

	if respRecipe.StatusCode != http.StatusOK {
		t.Error(respRecipe.StatusCode)
	}
	fmt.Println("testeing handlerFood POST method recipe")

	// Testing method GET ***********************************************

	// Testing method DELETE

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
