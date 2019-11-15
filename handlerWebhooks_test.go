package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWebhooksPost(t *testing.T) {
	// Test Post method for endpoint /cravings/webhooks/ ******************'
	webH := Webhook{Event: "testevent", URL: "www.testurl.com"} // create a webhook with event and url to send as body
	req, _ := json.Marshal(webH)
	reqTest := bytes.NewReader(req)                                   // convert struct to *Reader
	r, err := http.NewRequest("POST", "/cravings/webhooks/", reqTest) // creats request with body

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponseRecorder

	handler := http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing webhooks POST method")
}

func TestHandlerWebhooksGetA(t *testing.T) {
	// Test Get method for endpoint /cravings/webhooks/ ***************************
	r, err := http.NewRequest("GET", "/cravings/webhooks/", nil) //creates request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder() // creates ResponseRecorder

	handler := http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing webhooks GET method for all webhooks")
}

func TestHandlerWebhooksGetO(t *testing.T) {
	// Test Get method for endpoint /cravings/webhooks/ID **************
	w := httptest.NewRecorder() // creates ResponseRecorder

	wh, err := DBReadAllWebhooks(w) // reads all webhooks from database

	if err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest("GET", "/cravings/webhooks/"+wh[1].ID, nil) // creats request

	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // creates ResponseRecorder

	handler := http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing webhooks GET method for one webhook")
}

func TestHandlerWebhooksDelete(t *testing.T) {
	// Test Delete method for endpoint /cravings/webhooks/ ****************
	w := httptest.NewRecorder() // creates ResponseRecorder

	wh, err := DBReadAllWebhooks(w) // reads all webhooks from database

	if err != nil {
		t.Error(err)
	}

	webH := Webhook{Event: "testevent", URL: "www.testurl.com"} // same webhooks as the one i created in POST test

	var temp string

	fmt.Println("webH: ", webH.Event)

	for i := range wh { // loops throue all webhooks
		fmt.Println("event: ", wh[i].Event)
		if wh[i].Event == webH.Event { // check if webhook is same ass the test webhook we made earlyer
			temp = wh[i].ID // sets temp to be the same as the id for temp webhook
			fmt.Println("tempStruct: ", wh[i])
		}
	}

	fmt.Println("temp: ", temp)

	tempstruct := Webhook{ID: temp} // creates temp struct to send with request
	req, _ := json.Marshal(tempstruct)
	reqTest := bytes.NewReader(req)                                     // convert struct to *Reader
	r, err := http.NewRequest("DELETE", "/cravings/webhooks/", reqTest) // creates requests

	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder() // creates ResponseRecorder

	handler := http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testing webhooks DELETE method")
}
