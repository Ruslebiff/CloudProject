package main

import (
	"cravings"
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

	http.HandleFunc("/", cravings.HandlerNil)
	http.HandleFunc("/register", cravings.HandlerRegister) // runs handler function
	http.HandleFunc("/meal", cravings.HandlerMeal)         // runs handler function
	http.HandleFunc("/status", cravings.HandlerStatus)     // runs handler function
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
