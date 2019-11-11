package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWebhooks(t *testing.T) {

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

	//r, err = http.NewRequest("GET")
}
