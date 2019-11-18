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

	Post := "POST"
	Get := "GET"
	Delete := "DELETE"
	URL := "/cravings/webhooks/"

	webH := Webhook{Event: "testevent", URL: "www.testurl.com"} // create a webhook with event and url to send as body
	s := Webhook{}

	// Test Post method for endpoint /cravings/webhooks/
	fmt.Println("testeing webhooks POST method")

	resp := ALLMethodWebhook(Post, URL, webH, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test Get method for endpoint /cravings/webhooks/
	fmt.Println("testing webhooks GET method for all webhooks")

	resp = ALLMethodWebhook(Get, URL, s, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test Get method for endpoint /cravings/webhooks/ID
	fmt.Println("testing webhooks GET method for one webhook")

	w := httptest.NewRecorder() // creates ResponseRecorder

	wh, err := DBReadAllWebhooks(w) // reads all webhooks from database

	if err != nil {
		t.Error(err)
	}

	resp = ALLMethodWebhook(Get, URL+wh[1].ID, s, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}

	// Test Delete method for endpoint /cravings/webhooks/
	fmt.Println("testing webhooks DELETE method")

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

	resp = ALLMethodWebhook(Delete, URL, tempstruct, t)

	if resp.StatusCode != http.StatusOK { // check that test went ok
		t.Error(resp.StatusCode)
	}
}

func ALLMethodWebhook(m string, url string, s Webhook, t *testing.T) *http.Response {

	r, err := http.NewRequest(m, url, nil) // creates request with body

	if err != nil {
		t.Error(err)
	}

	if len(s.Event) > 0 || len(s.ID) > 0 {

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

	handler := http.HandlerFunc(HandlerWebhooks) // test handlerFood
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	return resp

}
