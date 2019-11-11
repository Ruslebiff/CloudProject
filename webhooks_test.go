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

}
