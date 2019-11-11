package cravings

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
