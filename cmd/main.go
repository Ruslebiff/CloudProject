package main

import (
	"assignment2-cloud/assignment2"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", assignment2.HandlerNil)
	http.HandleFunc("/register") // runs handler function
	http.HandleFunc("/meal")     // runs handler function
	http.HandleFunc("/status")   // runs handler function
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
