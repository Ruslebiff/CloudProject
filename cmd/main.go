package main

import (
	"CloudProject"
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
	http.HandleFunc("/register", CloudProject.HandlerRegister) // runs handler function
	http.HandleFunc("/meal", CloudProject.HandlerMeal)         // runs handler function
	http.HandleFunc("/status", CloudProject.HandlerStatus)     // runs handler function
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
