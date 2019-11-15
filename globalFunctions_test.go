package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRequest(t *testing.T) {
	TestURL := "http://www.google.com"
	TestClient := http.DefaultClient

	test, err := DoRequest(TestURL, TestClient) // test func with a test url

	if test.StatusCode != http.StatusOK {
		t.Error(test)
	}

	if err != nil {
		t.Error(err)
	}

	defer test.Body.Close()

	fmt.Println("TestDoRequest")
}

func TestDoRequestFail(t *testing.T) {

}

func TestQueryGet(t *testing.T) {
	r, err := http.NewRequest("GET", "/cravings/food/", nil) //creat a request without anny body

	if err != nil {
		t.Error(err)
	}

	Test := "app_id"

	test := QueryGet(Test, "", r) // test to read app_id and expecting it to return an empty string

	if test != "" {
		t.Error("not found")
	}

	fmt.Println("TestQueryGet")
}

func TestCallURL(t *testing.T) {
	w := httptest.NewRecorder()                     // creates ResponsRecoder
	TestRecipe := Recipe{RecipeName: "TestCallURl"} // create a struct with a name

	err := CallURL(RecipeCollection, TestRecipe, w) // check that we can callUrl

	if err != nil {
		t.Error(err)
	}
}

func TestReadIngredients(t *testing.T) {
	w := httptest.NewRecorder() // create ResponseRecorder

	var testIngredient []string

	a1 := "cheese"
	a2 := "milk|70"
	a3 := "cardamom|500|g"

	testIngredient = append(testIngredient, a1)
	testIngredient = append(testIngredient, a2)
	testIngredient = append(testIngredient, a3)

	fmt.Println("testIngredient", testIngredient)

	test := ReadIngredients(testIngredient, w) // test to read ingredients

	if len(test) == 0 { // check that it dont return an empty slice
		t.Error("somthing went wrong")
	}
}

func TestConvertUnit(t *testing.T) {
	testIngredient := Ingredient{Name: "TestIngrdient", Quantity: 1000, Unit: "g"}
	testUnitKG := "kg"
	testUnitG := "g"

	ConvertUnit(&testIngredient, testUnitKG) // test convert from g to kg
	fmt.Println("testIngredient: ", testIngredient)

	if testIngredient.Quantity != 1 {
		t.Error("quanity did not get converted")
	}

	ConvertUnit(&testIngredient, testUnitG) // test convert from kg to g

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
