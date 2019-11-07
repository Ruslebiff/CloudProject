package cravings

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

// HandlerNil is the default handler for bad URL requests
func HandlerNil(w http.ResponseWriter, r *http.Request) { //standard default response
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)

	// ********** Informatian about endpoints *************

	file, err := os.Open("nil.text") // opens text file
	if err != nil {
		fmt.Println("Can't open file: ", err)
	}

	defer file.Close() // close file at the end

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // loops true file
		fmt.Fprintln(w, scanner.Text()) // print out one and one line
	}

}
