package cravings

import (
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
