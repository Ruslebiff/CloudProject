package cravings

import (
	"fmt"
	"testing"
)

// func TestDoRequest(t *testing.T) { // needs to bee fixed!!!!!!!!
// 	TestURL := "www.google.com"
// 	TestClient := http.DefaultClient

// 	_ = func(w http.ResponseWriter) {
// 		test := DoRequest(TestURL, TestClient, w)
// 		if test.StatusCode != http.StatusOK {
// 			t.Error(test)
// 		}
// 	}
// 	fmt.Println("TestDoRequest")

// }

// func TestQueryGet(t *testing.T) { // needs to bee fixed!!!!!!!!

// 	Test := "app_id"

// 	_ = func(w http.ResponseWriter, r *http.Request) {
// 		test := QueryGet(Test, w, r)
// 		if test == "" {
// 			t.Error("not found")
// 		}
// 	}

// 	fmt.Println("TestQueryGet")

// }

func TestCallURL(t *testing.T) {
	TestRecipe := Recipe{RecipeName: "TestCallURl"}

	err := CallURL(RecipeCollection, TestRecipe)
	if err != nil {
		t.Error(err)
	}
}

func TestReadIngredients(t *testing.T) {

	var testIngredient []string
	a1 := "cheese"
	a2 := "milk|70"
	a3 := "cardamom|500|g"

	testIngredient = append(testIngredient, a1)
	testIngredient = append(testIngredient, a2)
	testIngredient = append(testIngredient, a3)
	fmt.Println("testIngredient", testIngredient)

	test := ReadIngredients(testIngredient)
	if len(test) == 0 {
		t.Error("somthing vent wrong")
	}

}

func TestConvertUnit(t *testing.T) {

}
