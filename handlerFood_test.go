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

	// Testing methon POST ***********************************************
	i := testIngredient{Token: testToken, Name: "turmeric", Unit: "g"}
	requestI, _ := json.Marshal(i)
	requestTest := bytes.NewReader(requestI)
	r, err := http.NewRequest("POST", "/cravings/food/ingredient", requestTest)
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	fmt.Println("testeing handlerFood POST method")

	//Testing method POST ***********************************************
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
