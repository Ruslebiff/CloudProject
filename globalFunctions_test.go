package cravings

import (
	"fmt"
	"testing"
)

// func TestDoRequest(t *testing.T) { // needs to bee fixed!!!!!!!!
// 	TestURL := "http://www.google.com"
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
		t.Error("somthing went wrong")
	}

}

func TestConvertUnit(t *testing.T) {

	testIngredient := Ingredient{Name: "TestIngrdient", Quantity: 1000, Unit: "g"}
	testUnitKG := "kg"
	testUnitG := "g"

	ConvertUnit(&testIngredient, testUnitKG)
	fmt.Println("testIngredient: ", testIngredient)
	if testIngredient.Quantity != 1 {
		t.Error("quanity did not get converted")
	}
	ConvertUnit(&testIngredient, testUnitG)
	fmt.Println("testIngredient: ", testIngredient)
	if testIngredient.Quantity != 1000 {
		t.Error("quanity did not get converted")
	}

	testIngredient2 := Ingredient{Name: "TestIngredient2", Quantity: 1000, Unit: "ml"}
	testUnitL := "l"
	testUnitDl := "dl"
	testUnitCl := "cl"
	testUnitMl := "ml"

	ConvertUnit(&testIngredient2, testUnitCl) // test convert from ml to cl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 100 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitDl) // test convert from cl to dl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 10 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitL) // test convert from dl to l
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitMl) // test convert from l to ml
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1000 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitDl) // test convert from ml to dl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 10 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitCl) // test convert from dl to cl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 100 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitL) // test convert from cl to l
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitCl) // test convert from l to cl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 100 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitMl) // test convert from cl to ml
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1000 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitL) // test convert from ml to l
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitDl) // test convert from l to dl
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 10 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient2, testUnitMl) // test convert from dl to ml
	fmt.Println("testingredient2: ", testIngredient2)
	if testIngredient2.Quantity != 1000 {
		t.Error("quanity did not get converted")
	}

}

func TestInitAPICredentials(t *testing.T) {
	err := InitAPICredentials()
	if err != nil {
		t.Error(err)
	}
}
