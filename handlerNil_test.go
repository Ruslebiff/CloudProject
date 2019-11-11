package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerNil(t *testing.T) {

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerNil)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode == http.StatusOK {
		t.Error(resp.StatusCode)
	}
	fmt.Println("test handlerNil")
}
