package main

import (
	"cravings"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var StartTime = time.Now() // sets start time

func main() {

	cravings.StartTime = StartTime // sends start time to ctavins.StartTime
	err := cravings.DBInit()
	if err != nil {
		fmt.Println("Failed to initialise database")
	} else {
		fmt.Println("Database init OK")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", cravings.HandlerNil)
	http.HandleFunc("/register/", cravings.HandlerRegister) // runs handler function
	http.HandleFunc("/status/", cravings.HandlerStatus)     // runs handler function
	http.HandleFunc("/meal/", cravings.HandlerMeal)         // runs handler function
	http.HandleFunc("/webhooks/", cravings.HandlerWebhooks) // runs handler function
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
