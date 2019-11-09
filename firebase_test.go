package cravings

import (
	"fmt"
	"testing"
	"time"
)

func TestFirebase(t *testing.T) {

	err := DBInit() // check that the database initalises
	if err != nil {
		t.Error(err) // feild to initalise

	}

	// testing of Ingredients functions for database **************
	i := Ingredient{Name: "TestIngredient"}

	err = DBSaveIngredient(&i)
	if err != nil {
		t.Error(err)
	}

	var ingredients []Ingredient
	ingredients, err = DBReadAllIngredients()
	if err != nil {
		t.Error(err)
	}

	_, err = DBReadIngredientByID(ingredients[1].ID)
	if err != nil {
		t.Error(err)
	}

	i2, err := DBReadIngredientByName(i.Name)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("test ingredient", i2)
	err = DBDelete(i2.ID, IngredientCollection)
	if err != nil {
		t.Error(err)
	}

	// testing of Recipe functions for database ****************************

	r := Recipe{RecipeName: "TestRecipe"}

	err = DBSaveRecipe(&r)
	if err != nil {
		t.Error(err)
	}

	var Recipes []Recipe
	Recipes, err = DBReadAllRecipes()
	if err != nil {
		t.Error(err)
	}

	_, err = DBReadRecipeByID(Recipes[1].ID)
	if err != nil {
		t.Error(err)
	}

	r2, err := DBReadRecipeByName(r.RecipeName)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("test recipe ", r2)
	err = DBDelete(r2.ID, RecipeCollection)
	if err != nil {
		t.Error(err)
	}

	// testing of Webhooks functions for database *****************************
	w := Webhook{Event: "Test", URL: "www.Test.no", Time: time.Now()}

	err = DBSaveWebhook(&w)
	if err != nil {
		t.Error(err)
	}

	var wh []Webhook
	wh, err = DBReadAllWebhooks()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("test webhooks ", wh[1].ID)
	err = DBDelete(wh[1].ID, WebhooksCollection)
	if err != nil {
		t.Error(err)
	}

}
