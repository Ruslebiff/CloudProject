package cravings

import (
	"fmt"
	"net/http"
)

func HandlerNil(w http.ResponseWriter, r *http.Request) { //standar default respons
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}
