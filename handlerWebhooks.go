package cravings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HandlerWebhooks Handler fun for webhooks
func HandlerWebhooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		Wh := Webhook{}

		err := json.NewDecoder(r.Body).Decode(&Wh) //decode to webhook
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		Wh.Event = strings.ToLower(Wh.Event)
		Wh.Time = time.Now() // sets time stamp

		err = DBSaveWebhook(&Wh) // saves webhook to firebase
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Webhooks " + Wh.URL + " has been regstrerd")

	case http.MethodGet:
		var webhooks []Webhook //Webhook DB
		webhooks, err := DBReadAllWebhooks()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		http.Header.Add(w.Header(), "Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(webhooks) // encode all webhooks
		if err != nil {
			http.Error(w, "Some thing went wrong"+err.Error(), http.StatusInternalServerError)
		}

	case http.MethodDelete:
		Wh := Webhook{}

		err := json.NewDecoder(r.Body).Decode(&Wh) //decode to webhook
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = DBDelete(Wh.ID, WebhooksCollection) //Deletes webhook from id
		if err != nil {
			fmt.Println("Error: ", err)
		}

	default:
		http.Error(w, "Method is invalid "+r.Method, http.StatusBadRequest)
	}

}
