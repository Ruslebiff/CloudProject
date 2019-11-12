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
	webH := Webhook{Event: "testevent", URL: "www.testurl.com"}
	req, _ := json.Marshal(webH)
	reqTest := bytes.NewReader(req)
	r, err := http.NewRequest("POST", "/cravings/webhooks/", reqTest)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerWebhooks)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks POST method")

	// Test Get method for endpoint /cravings/webhooks/ ***************************
	r, err = http.NewRequest("GET", "/cravings/webhooks/", nil)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerWebhooks)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks GET method for all webhooks")

	// Test Get method for endpoint /cravings/webhooks/ID **************

	wh, err := DBReadAllWebhooks()
	if err != nil {
		t.Error(err)
	}

	r, err = http.NewRequest("GET", "/cravings/webhooks/"+wh[1].ID, nil)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerWebhooks)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}

	fmt.Println("testeing webhooks GET method for one webhook")

	// Test Delete method for endpoint /cravings/webhooks/ ****************

	var temp string
	fmt.Println("webH: ", webH.Event)
	for i := range wh {
		fmt.Println("event: ", wh[i].Event)
		if wh[i].Event == webH.Event {
			temp = wh[i].ID
			fmt.Println("tempStruct: ", wh[i])
		}
	}

	fmt.Println("temp: ", temp)

	tempstruct := Webhook{ID: temp}
	req, _ = json.Marshal(tempstruct)
	reqTest = bytes.NewReader(req)
	r, err = http.NewRequest("DELETE", "/cravings/webhooks/", reqTest)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(HandlerWebhooks)
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("testeing webhooks DELETE method")

}
