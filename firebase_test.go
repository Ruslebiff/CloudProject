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
	"time"
)

func TestDBInit(t *testing.T) {
	err := DBInit() // check that the database initialises
	if err != nil {
		t.Error(err) // failed to initialise
	}

	fmt.Println("testing DBInit")
}

func TestDBIngredient(t *testing.T) {
	// testing of Ingredients functions for database **************
	w := httptest.NewRecorder() // creates ResponsRecoder

	i := Ingredient{Name: "TestIngredient"}

	err := DBSaveIngredient(&i, w) //test saving Ingredient
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	var ingredients []Ingredient
	ingredients, err = DBReadAllIngredients(w) // test read all ingredients

	if err != nil {
		t.Error(err)
	}

	fmt.Println("Ingredients: ", len(ingredients))

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	i2, err := DBReadIngredientByName(i.Name, w) // test read ingredients by name
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	fmt.Println("test ingredient", i2)
	err = DBDelete(i2.ID, IngredientCollection, w) // test delete and delete test ingredient

	if err != nil {
		t.Error(err)
	}

	fmt.Println("testing DBIngredient")
}

func TestDBRecipe(t *testing.T) {
	// testing of Recipe functions for database ****************************
	w := httptest.NewRecorder() // creates ResponsRecoder

	r := Recipe{RecipeName: "TestRecipe"}

	err := DBSaveRecipe(&r, w) // test saving recipe
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	var Recipes []Recipe
	Recipes, err = DBReadAllRecipes(w) // test reading all recipes

	if err != nil {
		t.Error(err)
	}

	fmt.Println("Recipe: ", len(Recipes))

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	r2, err := DBReadRecipeByName(r.RecipeName, w) // test reading recipe by name
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	fmt.Println("test recipe ", r2)
	err = DBDelete(r2.ID, RecipeCollection, w) // Delets test recipe

	if err != nil {
		t.Error(err)
	}

	fmt.Println("testing DBRecipe")
}

func TestDBWebhook(t *testing.T) {
	// testing of Webhooks functions for database *****************************
	w := httptest.NewRecorder() // creates ResponsRecoder

	testWh := Webhook{Event: "Test", URL: "www.Test.no", Time: time.Now()}

	err := DBSaveWebhook(&testWh, w) // test saving webhook

	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	var wh []Webhook
	wh, err = DBReadAllWebhooks(w) // test reading all webhooks

	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	for i := range wh {
		if wh[i].Event == testWh.Event {
			fmt.Println("test webhooks ", wh[i].ID)
			err = DBDelete(wh[i].ID, WebhooksCollection, w) // delets test webhook

			if err != nil {
				t.Error(err)
			}
		}
	}

	w = httptest.NewRecorder() // resets ResponsRecoder for next func test

	err = DBDelete("", WebhooksCollection, w) // test deleting somthing that dont exist, error is supposed to be sendt

	if err == nil {
		t.Error(err)
	}

	fmt.Println("testing DBWebhook")
}

func TestDBCheckAuthorization(t *testing.T) {
	file, err := os.Open("testToken.txt") // opens text file

	if err != nil {
		fmt.Println("Can't open file: ", err)
	}

	defer file.Close() // close file at the end

	scanner := bufio.NewScanner(file)

	var testToken string
	for scanner.Scan() { // loop throue lengt of text file
		testToken = scanner.Text() // sets testToken to the token read from file
	}
	fmt.Println("text: ", testToken)

	testStruct := Token{AuthToken: testToken} // creats test struct that will be sent as json body for GET request
	request, _ := json.Marshal(testStruct)
	requestTest := bytes.NewReader(request)
	r, err := http.NewRequest("GET", "/cravings/food/", requestTest) // create request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // create ResponseRecorder

	testBool, _ := DBCheckAuthorization(w, r) //test for vallid authorization with a vallid token

	if testBool == false {
		t.Error("Token was not vallid")
	}

	fmt.Println("test DBCheckAuthorization")
}
