package cravings

import (
	"fmt"
	"net/http"
)

func HandlerNil(w http.ResponseWriter, r *http.Request) { //standard default response
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)

	// ********** Informatian about endpoints *************

	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Endpoints available: ")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "************************************************************************")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Register endpoint: /cravings/register/")
	fmt.Fprintln(w, "Here you can registrate a recipe")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "************************************************************************")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Status endpoint: /cravings/status/")
	fmt.Fprintln(w, "Here you can see that status for the website")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "************************************************************************")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Meal endpoint: /cravings/meal/")
	fmt.Fprintln(w, "Here you can get recipe to make, out from what ingridients you have")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "************************************************************************")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Webhooks endpoint: /cravings/webhooks/")
	fmt.Fprintln(w, "Here you can get information about webhooks for this website")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "************************************************************************")

}
