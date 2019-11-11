package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWebhooks(t *testing.T) {

	// Test Get method for endpoint /cravings/webhooks/ ***************************
	r, err := http.NewRequest("GET", "/cravings/webhooks/", nil)
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
	fmt.Println("testeing webhooks GET method")

	// Test Get method for endpoint /cravings/webhooks/ID **************

	wh, err := DBReadAllWebhooks()
	if err != nil {
		t.Error(err)
	}

	r2, err := http.NewRequest("GET", "/cravings/webhooks/"+wh[1].ID, nil)
	if err != nil {
		t.Error(err)
	}

	w2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(HandlerWebhooks)
	handler2.ServeHTTP(w2, r2)

	resp2 := w2.Result()

	if resp2.StatusCode != http.StatusOK {
		t.Error(resp2.StatusCode)
	}

	// Test Delete method for endpoint /cravings/webhooks/ ****************

	// testWebhook := Webhook{Event: "TestWebhookDelete", URL: "www.Test.no", Time: time.Now()}
	// testWebhook2 := Webhook{}

	// err = DBSaveWebhook(&testWebhook) // test saving webhook
	// if err != nil {
	// 	t.Error(err)
	// }

	// test, err := DBReadAllWebhooks()
	// if err != nil {
	// 	t.Error(err)
	// }

	// for i := range test {
	// 	if test[i].Event == testWebhook.Event {
	// 		testWebhook2 = test[i]
	// 	}
	// }

	// fmt.Println("testWebhook2: ", testWebhook2)

	// r3, err := http.NewRequest("DELETE", "/cravings/webhooks/", nil)
	// if err != nil {
	// 	t.Error(err)
	// }

	// w3 := httptest.NewRecorder()
	// handler3 := http.HandlerFunc(HandlerWebhooks)
	// handler3.ServeHTTP(w3, r3)

	// resp3 := w3.Result()

	// if resp3.StatusCode != http.StatusOK {
	// 	t.Error(resp3.StatusCode)
	// }

}
