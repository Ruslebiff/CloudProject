package cravings

import (
	"fmt"
	"net/http"
	"strings"
)

// Function which registers either an ingredient or a recipe
func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[2]
	fmt.Println(endpoint)
}

//	parts :=  Siste part er enteen ingredient eller recipe
// 	Kjør switch og kall på respektiv handler
