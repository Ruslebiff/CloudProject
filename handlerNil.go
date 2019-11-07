package cravings

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

// HandlerNil is the default handler for bad URL requests
func HandlerNil(w http.ResponseWriter, r *http.Request) { //standard default response
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)

	// ********** Informatian about endpoints *************

	file, err := os.Open("nil.text")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}

}
