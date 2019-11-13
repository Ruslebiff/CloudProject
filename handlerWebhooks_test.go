package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWebhooks(t *testing.T) {

	// Test Post method for endpoint /cravings/webhooks/ ******************'
	webH := Webhook{Event: "testevent", URL: "www.testurl.com"} // create a webhook with event and url to send as body
	req, _ := json.Marshal(webH)
	reqTest := bytes.NewReader(req)                                   // convert struct to *Reader
	r, err := http.NewRequest("POST", "/cravings/webhooks/", reqTest) // creats request with body
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()                  // creates ResponseRecorder
	handler := http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks POST method")

	// Test Get method for endpoint /cravings/webhooks/ ***************************
	r, err = http.NewRequest("GET", "/cravings/webhooks/", nil) //creates request
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()                  // creates ResponsRecorder
	handler = http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks GET method for all webhooks")

	// Test Get method for endpoint /cravings/webhooks/ID **************

	w = httptest.NewRecorder() // creates ResponsRecorder

	wh, err := DBReadAllWebhooks(w) // reads all webhooks from database
	if err != nil {
		t.Error(err)
	}

	r, err = http.NewRequest("GET", "/cravings/webhooks/"+wh[1].ID, nil) // creats request
	if err != nil {
		t.Error(err)
	}

	handler = http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing webhooks GET method for one webhook")

	// Test Delete method for endpoint /cravings/webhooks/ ****************

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
	req, _ = json.Marshal(tempstruct)
	reqTest = bytes.NewReader(req)                                     // convert struct to *Reader
	r, err = http.NewRequest("DELETE", "/cravings/webhooks/", reqTest) // creates requests
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()                  // creates ResponsRecorder
	handler = http.HandlerFunc(HandlerWebhooks) // test handlerWebhooks
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks DELETE method")

}
