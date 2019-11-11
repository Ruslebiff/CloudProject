package cravings

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFirebase(t *testing.T) {

	err := DBInit() // check that the database initalises
	if err != nil {
		t.Error(err) // failed to initalise
	}

	// testing of Ingredients functions for database **************
	i := Ingredient{Name: "TestIngredient"}

	err = DBSaveIngredient(&i) //test saving Ingredient
	if err != nil {
		t.Error(err)
	}

	var ingredients []Ingredient
	ingredients, err = DBReadAllIngredients() // test read all ingredients
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Ingredients: ", len(ingredients))

	i2, err := DBReadIngredientByName(i.Name) // test read ingredients by name
	if err != nil {
		t.Error(err)
	}

	fmt.Println("test ingredient", i2)
	err = DBDelete(i2.ID, IngredientCollection) // test delete and delete test ingredient
	if err != nil {
		t.Error(err)
	}

	// testing of Recipe functions for database ****************************

	r := Recipe{RecipeName: "TestRecipe"}

	err = DBSaveRecipe(&r) // test saving recipe
	if err != nil {
		t.Error(err)
	}

	var Recipes []Recipe
	Recipes, err = DBReadAllRecipes() // test reading all recipes
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Recipe: ", len(Recipes))

	r2, err := DBReadRecipeByName(r.RecipeName) // test reading recipe by name
	if err != nil {
		t.Error(err)
	}

	fmt.Println("test recipe ", r2)
	err = DBDelete(r2.ID, RecipeCollection) // Delets test recipe
	if err != nil {
		t.Error(err)
	}

	// testing of Webhooks functions for database *****************************
	w := Webhook{Event: "Test", URL: "www.Test.no", Time: time.Now()}

	err = DBSaveWebhook(&w) // test saving webhook
	if err != nil {
		t.Error(err)
	}

	var wh []Webhook
	wh, err = DBReadAllWebhooks() // test reading all webhooks
	if err != nil {
		t.Error(err)
	}

	for i := range wh {
		if wh[i].Event == w.Event {
			fmt.Println("test webhooks ", wh[i].ID)
			err = DBDelete(wh[i].ID, WebhooksCollection) // delets test webhook
			if err != nil {
				t.Error(err)
			}
		}
	}

}

func TestDBCheckAuthorization(t *testing.T) {

	file, err := os.Open("testToken.txt") // opens text file
	if err != nil {
		fmt.Println("Can't open file: ", err)
	}

	defer file.Close() // close file at the end

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println("text: ", scanner.Text())
		testBool := DBCheckAuthorization(scanner.Text()) //test for vallid authorization with a vallid token
		if testBool == false {
			t.Error("Token was not vallid")
		}
	}

	testBool := DBCheckAuthorization("") // test for unvallid authorization with a unvallig token
	if testBool == true {
		t.Error("Somthing went vrong, unvallid token is not supose to return true")
	}

}
