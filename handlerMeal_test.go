package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerMeal(t *testing.T) {

	// test GET method for Meal

	r, err := http.NewRequest("GET", "/cravings/meal/", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerMeal)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing handlerMeal GET method")
}
