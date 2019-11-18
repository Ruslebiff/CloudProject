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

	temp1 := Ingredient{Name: "milk", Unit: "l", Quantity: 2}      // create ingredient for test
	temp2 := Ingredient{Name: "olive oil", Unit: "l", Quantity: 1} // create ingredient for test

	tempstruct = append(tempstruct, temp1)
	tempstruct = append(tempstruct, temp2)

	i := []Ingredient{}

	// test POST method for Meal
	fmt.Println("testing handlerMeal POST method")

	resp := ALLMethodMeal(Post, "/cravings/meal/", tempstruct, t)

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}

	// test GET method for Meal
	fmt.Println("testing handlerMeal GET method")

	resp = ALLMethodMeal(Get, "/cravings/meal/?ingredients=milk|2|l_olive oil|1|l", i, t)

	if resp.StatusCode != http.StatusOK { // check that everything is ok
		t.Error(resp.StatusCode)
	}

}

// func TestHandlerMealPost(t *testing.T) {
// 	// test POST method for Meal
// 	var tempstruct []Ingredient // creat a tempstruct for testing

// 	temp1 := Ingredient{Name: "milk", Unit: "l", Quantity: 2}      // create ingredient for test
// 	temp2 := Ingredient{Name: "olive oil", Unit: "l", Quantity: 1} // create ingredient for test

// 	tempstruct = append(tempstruct, temp1)
// 	tempstruct = append(tempstruct, temp2)

// 	req, _ := json.Marshal(tempstruct)
// 	reqTest := bytes.NewReader(req)                               // convert struct over to a *reader
// 	r, err := http.NewRequest("POST", "/cravings/meal/", reqTest) // creats a Request with a struct

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // create ResponsRecorder

// 	handler := http.HandlerFunc(HandlerMeal) // test handlerMeal
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that everything is ok
// 		t.Error(resp.StatusCode)
// 	}

// 	fmt.Println("testing handlerMeal POST method")
// }

// func TestHandlerMealGet(t *testing.T) {
// 	// test GET method for Meal
// 	r, err := http.NewRequest("GET", "/cravings/meal/?ingredients=milk|2|l_olive oil|1|l", nil) // create request

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	w := httptest.NewRecorder() // create ResponsRecorder for all test

// 	handler := http.HandlerFunc(HandlerMeal) // test HandlerMeal
// 	handler.ServeHTTP(w, r)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK { // check that everything is ok
// 		t.Error(resp.StatusCode)
// 	}

// 	fmt.Println("testing handlerMeal GET method")
// }

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
