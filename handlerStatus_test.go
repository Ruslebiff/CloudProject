package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerStatus(t *testing.T) {

	r, err := http.NewRequest("GET", "/cravings/status/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerStatus)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("test handlerStatus")
}
