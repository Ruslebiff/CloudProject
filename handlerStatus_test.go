package cravings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerStatus(t *testing.T) {
	r, err := http.NewRequest("GET", "/cravings/status/", nil) // creats request

	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()                // create ResponsRcorder
	handler := http.HandlerFunc(HandlerStatus) // test handlerNil
	handler.ServeHTTP(w, r)                    // sends with request and respons
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // check if handler worked ass it should
		t.Error(resp.StatusCode)
	}

	fmt.Println("test handlerStatus")
}
