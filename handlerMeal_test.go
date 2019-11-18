package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandelerMeal(t *testing.T) {
	var tempstruct []Ingredient // creat a tempstruct for testing

	Post := "POST"
	Get := "GET"
	URL := "/cravings/meal/"

	temp1 := Ingredient{Name: "milk", Unit: "l", Quantity: 2}      // create ingredient for test
	temp2 := Ingredient{Name: "olive oil", Unit: "l", Quantity: 1} // create ingredient for test

	tempstruct = append(tempstruct, temp1)
	tempstruct = append(tempstruct, temp2)

	i := []Ingredient{}

	// test POST method for Meal
	fmt.Println("testing handlerMeal POST method")

	resp := ALLMethodMeal(Post, URL, tempstruct, t)

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}

	// test GET method for Meal
	fmt.Println("testing handlerMeal GET method")

	resp = ALLMethodMeal(Get, URL+"?ingredients=milk|2|l_olive oil|1|l", i, t)

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}
}

func ALLMethodMeal(m string, url string, s []Ingredient, t *testing.T) *http.Response {
	r, err := http.NewRequest(m, url, nil) // creates request with body

	if err != nil {
		t.Error(err)
	}

	if len(s) > 0 {

		fmt.Println("len > 0")

		test := s
		req, _ := json.Marshal(test)
		reqTest := bytes.NewReader(req) // convert struct to *Reader

		r, err = http.NewRequest(m, url, reqTest) // creates request with body

		if err != nil {
			t.Error(err)
		}
	}

	w := httptest.NewRecorder() // creates ResponsRecorder

	handler := http.HandlerFunc(HandlerMeal) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	return resp
}
