package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerNil(t *testing.T) {

	r, err := http.NewRequest("GET", "/", nil) // creats request
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()             // create ResponsRcorder
	handler := http.HandlerFunc(HandlerNil) // test handlerNil
	handler.ServeHTTP(w, r)                 // sends with request and respons

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK { // check if handler worked ass it should
		t.Error(resp.StatusCode)
	}
	fmt.Println("test handlerNil")
}
