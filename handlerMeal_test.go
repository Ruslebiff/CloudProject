package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerMeal(t *testing.T) {

	w := httptest.NewRecorder() // create ResponsRecorder for all test

	// test POST method for Meal

	var tempstruct []Ingredient // creat a tempstruct for testing

	temp1 := Ingredient{Name: "milk", Unit: "l", Quantity: 2}      // create ingredient for test
	temp2 := Ingredient{Name: "olive oil", Unit: "l", Quantity: 1} // create ingredient for test

	tempstruct = append(tempstruct, temp1)
	tempstruct = append(tempstruct, temp2)

	req, _ := json.Marshal(tempstruct)
	reqTest := bytes.NewReader(req)                               // convert struct over to a *reader
	r, err := http.NewRequest("POST", "/cravings/meal/", reqTest) // creats a Request with a struct
	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(HandlerMeal) // test handlerMeal
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testing handlerMeal POST method")

	// test GET method for Meal

	r, err = http.NewRequest("GET", "/cravings/meal/?ingredients=milk|2|l_olive oil|1|l", nil) // create request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerMeal) // test HandlerMeal
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing handlerMeal GET method")
}
