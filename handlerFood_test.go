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

	return testToken
}

func TestHandlerFood(t *testing.T) {

	Post := "POST"
	Get := "GET"
	Delete := "DELETE"

	testToken := TempToken()

	if testToken == "" {
		t.Error("Token was not read from file")
	}

	in := TestIngredient{Token: testToken, Name: "turmeric", Unit: "g"} // creates test ingredient
	i := TestIngredient{}

	ingredient1 := Ingredient{Name: "milk", Quantity: 5, Unit: "ml"}
	ingredient2 := Ingredient{Name: "salt", Quantity: 2, Unit: "tablespoon"}
	ingredient3 := Ingredient{Name: "flour", Quantity: 1, Unit: "kg"}

	var testI []Ingredient

	testI = append(testI, ingredient1)
	testI = append(testI, ingredient2)
	testI = append(testI, ingredient3)

	re := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST", Ingredients: testI} // creates test recipe

	r := TestRecipe{}

	testIngred := TestIngredient{Token: testToken}
	testIngred2 := TestIngredient{Token: testToken, Name: "Something"}

	testRec := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST"}
	testRec1 := TestRecipe{Token: testToken}
	testRec2 := TestRecipe{Token: testToken, RecipeName: "something"}

	// Testing method POST for Ingredient
	fmt.Println("testing handlerFood POST method ingredient")

	resp := ALLMethodIngredient(Post, "/cravings/food/ingredient", in, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method POST for Recipe
	fmt.Println("testing handlerFood POST method recipe")

	resp = ALLMethodRecipe(Post, "/cravings/food/recipe", re, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test to get a POST request without a body expeting error
	fmt.Println("testing handlerFood POST method for error")

	resp = ALLMethodRecipe(Post, "/cravings/food/recipe", r, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method GET for all ingredients
	fmt.Println("testing handlerFood GET method for all ingredient")

	resp = ALLMethodIngredient(Get, "/cravings/food/ingredient/", i, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method GET for one ingredient
	fmt.Println("testing handlerFood GET method for one ingredient")

	resp = ALLMethodIngredient(Get, "/cravings/food/ingredient/turmeric", i, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method GET for all recipes
	fmt.Println("testing handlerFood GET method for all recipe")

	resp = ALLMethodRecipe(Get, "/cravings/food/recipe/", r, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method GET for one recipe
	fmt.Println("testing handlerFood GET method for one recipe")

	resp = ALLMethodRecipe(Get, "/cravings/food/recipe/TestRecipePOST", r, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test Geting somthing that dosent exists
	fmt.Println("testing handlerFood GET method for error in recipe")

	resp = ALLMethodRecipe(Get, "/cravings/food/recipe/Somthing", r, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test Geting somthing that dosent exists
	fmt.Println("testing handlerFood GET method for error in ingredient")

	resp = ALLMethodIngredient(Get, "/cravings/food/ingredient/Somthing", i, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for ingredient
	fmt.Println("testing handlerFood DELETE method ingredient")

	resp = ALLMethodIngredient(Delete, "/cravings/food/ingredient", in, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for ingredient with only token expeting error
	fmt.Println("testing handlerFood DELETE method ingredient with only token")

	resp = ALLMethodIngredient(Delete, "/cravings/food/ingredient", testIngred, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for ingredient with a ingredient that dont exsist expeting error
	fmt.Println("testing handlerFood DELETE method ingredient with wrong name")

	resp = ALLMethodIngredient(Delete, "/cravings/food/ingredient", testIngred2, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for recipe
	fmt.Println("testing handlerFood DELETE method recipe")

	resp = ALLMethodRecipe(Delete, "/cravings/food/recipe", testRec, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for recipe with only token expeting error
	fmt.Println("testing handlerFood DELETE method recipe only token")

	resp = ALLMethodRecipe(Delete, "/cravings/food/recipe", testRec1, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing method DELETE for recipe with a name that dont exsist
	fmt.Println("testing handlerFood DELETE method recipe with wrong name")

	resp = ALLMethodRecipe(Delete, "/cravings/food/recipe", testRec2, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Testing to delete without body, expecting to get error
	fmt.Println("Testing handlerFood DELETE method without body")

	resp = ALLMethodRecipe(Delete, "/cravings/food/recipe", r, t)

	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
		t.Error(resp.StatusCode)
	}

}

// func TestHandlerFoodPostI(t *testing.T) {
// 	// Testing method POST for Ingredient
// 	fmt.Println("testing handlerFood POST method ingredient")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	i := TestIngredient{Token: testToken, Name: "turmeric", Unit: "g"} // creates test ingredient
// 	req, _ := json.Marshal(i)
// 	reqTest := bytes.NewReader(req)                                         // convert struct to *Reader
// 	r, err := http.NewRequest("POST", "/cravings/food/ingredient", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodPostR(t *testing.T) {
// 	// Testing method POST for Recipe
// 	fmt.Println("testing handlerFood POST method recipe")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	ingredient1 := Ingredient{Name: "milk", Quantity: 5, Unit: "ml"}
// 	ingredient2 := Ingredient{Name: "salt", Quantity: 2, Unit: "tablespoon"}
// 	ingredient3 := Ingredient{Name: "flour", Quantity: 1, Unit: "kg"}

// 	var testI []Ingredient

// 	testI = append(testI, ingredient1)
// 	testI = append(testI, ingredient2)
// 	testI = append(testI, ingredient3)

// 	re := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST", Ingredients: testI} // creates test recipe
// 	req, _ := json.Marshal(re)
// 	reqTest := bytes.NewReader(req)                                     // convert struct to *Reader
// 	r, err := http.NewRequest("POST", "/cravings/food/recipe", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodPostError(t *testing.T) {
// 	// Test to get a POST request without a body expeting error
// 	fmt.Println("testing handlerFood POST method for error")

// 	r, err := http.NewRequest("POST", "/cravings/food/recipe", nil) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}

// }

// func TestHandlerFoodGetAI(t *testing.T) {
// 	// Testing method GET for all ingredients
// 	fmt.Println("testing handlerFood GET method for all ingredient")

// 	r, err := http.NewRequest("GET", "/cravings/food/ingredient/", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodGetOI(t *testing.T) {
// 	// Testing method GET for one ingredient
// 	fmt.Println("testing handlerFood GET method for one ingredient")

// 	r, err := http.NewRequest("GET", "/cravings/food/ingredient/turmeric", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodGetAR(t *testing.T) {
// 	// Testing method GET for all recipes
// 	fmt.Println("testing handlerFood GET method for all recipe")

// 	r, err := http.NewRequest("GET", "/cravings/food/recipe/", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodGetOR(t *testing.T) {
// 	// Testing method GET for one recipe
// 	fmt.Println("testing handlerFood GET method for one recipe")

// 	r, err := http.NewRequest("GET", "/cravings/food/recipe/TestRecipePOST", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodGetRError(t *testing.T) {
// 	// Test Geting somthing that dosent exists
// 	fmt.Println("testing handlerFood GET method for error in recipe")

// 	r, err := http.NewRequest("GET", "/cravings/food/recipe/Somthing", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodGetIError(t *testing.T) {
// 	// Test Geting somthing that dosent exists
// 	fmt.Println("testing handlerFood GET method for error in ingredient")

// 	r, err := http.NewRequest("GET", "/cravings/food/ingredient/Somthing", nil) // creates request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodDeleteI(t *testing.T) {
// 	// Testing method DELETE for ingredient
// 	fmt.Println("testing handlerFood DELETE method ingredient")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test := TestIngredient{Token: testToken, Name: "turmeric"}
// 	req, _ := json.Marshal(test)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodDeleteIErr(t *testing.T) {
// 	// Testing method DELETE for ingredient with only token expeting error
// 	fmt.Println("testing handlerFood DELETE method ingredient with only token")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test := TestIngredient{Token: testToken}
// 	req, _ := json.Marshal(test)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}

// }

// func TestHandlerFoodDeleteIErr2(t *testing.T) {
// 	// Testing method DELETE for ingredient with a ingredient that dont exsist expeting error
// 	fmt.Println("testing handlerFood DELETE method ingredient with wrong name")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test := TestIngredient{Token: testToken, Name: "Somthing"}
// 	req, _ := json.Marshal(test)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/ingredient", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}

// }

// func TestHandlerFoodDeleteR(t *testing.T) {
// 	// Testing method DELETE for recipe
// 	fmt.Println("testing handlerFood DELETE method recipe")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test2 := TestRecipe{Token: testToken, RecipeName: "TestRecipePOST"}
// 	req, _ := json.Marshal(test2)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/recipe", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodDeleteRErr(t *testing.T) {
// 	// Testing method DELETE for recipe with only token expeting error
// 	fmt.Println("testing handlerFood DELETE method recipe only token")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test2 := TestRecipe{Token: testToken}
// 	req, _ := json.Marshal(test2)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/recipe", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodDeleteRErr2(t *testing.T) {
// 	// Testing method DELETE for recipe with a name that dont exsist
// 	fmt.Println("testing handlerFood DELETE method recipe with wrong name")

// 	testToken := TempToken()

// 	if testToken == "" {
// 		t.Error("Token was not read from file")
// 	}

// 	test2 := TestRecipe{Token: testToken, RecipeName: "something"}
// 	req, _ := json.Marshal(test2)
// 	reqTest := bytes.NewReader(req) // convert struct to *Reader

// 	r, err := http.NewRequest("DELETE", "/cravings/food/recipe", reqTest) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

// func TestHandlerFoodDeleteError(t *testing.T) {
// 	// Testing to delete without body, expecting to get error
// 	fmt.Println("Testing handlerFood DELETE method without body")

// 	r, err := http.NewRequest("DELETE", "/cravings/food/recipe", nil) // creates request with body

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // creates ResponsRecorder

// 	handler := http.HandlerFunc(HandlerFood) // test handlerFood
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusBadRequest { // check that test went ok
// 		t.Error(resp.StatusCode)
// 	}
// }

func ALLMethodIngredient(m string, url string, s TestIngredient, t *testing.T) *http.Response {

	r, err := http.NewRequest(m, url, nil) // creates request with body

	if err != nil {
		t.Error(err)
	}

	if len(s.Token) > 0 {

		test := s
		req, _ := json.Marshal(test)
		reqTest := bytes.NewReader(req) // convert struct to *Reader

		r, err = http.NewRequest(m, url, reqTest) // creates request with body

		if err != nil {
			t.Error(err)
		}
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	return resp

}

func ALLMethodRecipe(m string, url string, s TestRecipe, t *testing.T) *http.Response {

	r, err := http.NewRequest(m, url, nil) // creates request with body

	if err != nil {
		t.Error(err)
	}

	if len(s.Token) > 0 {

		test := s
		req, _ := json.Marshal(test)
		reqTest := bytes.NewReader(req) // convert struct to *Reader

		r, err = http.NewRequest(m, url, reqTest) // creates request with body

		if err != nil {
			t.Error(err)
		}
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerFood) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	return resp

}
