package cravings

import "testing"

func TestFirebase(t *testing.T) {
	db := FirestoreDatabase{}

	err := DBInit() // check that the database initalises
	if err != nil {
		t.Error(err) // feild to initalise

	}

	err = DBDelete(&Recipe{}.ID, "")
	if err != nil {
		t.Error(err)
	}

	err = DBDelete(&Ingredient{}.ID, "")
	if err != nil {
		t.Error(err)
	}

	var r []Recipe
	r, err := DBReadAllRecipes()
	if err != nil {
		t.Error(err)
	}

	var i []Ingredient
	i, err := DBReadAllIngredients()
	if err != nil {
		t.Error(err)
	}

}
