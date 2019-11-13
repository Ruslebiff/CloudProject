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

		err = DBSaveWebhook(&Wh, w) // saves webhook to firebase
		if err != nil {
			fmt.Fprintln(w, "Error saving webhook ", http.StatusInternalServerError)
		}

		fmt.Fprintln(w, "Webhooks "+Wh.URL+" has been regstrerd")

	case http.MethodGet:
		var webhooks []Webhook //Webhook DB
		parts := strings.Split(r.URL.Path, "/")

		http.Header.Add(w.Header(), "Content-Type", "application/json")

		webhooks, err := DBReadAllWebhooks(w) // reads all webhooks from database
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if parts[3] != "" { //check if an id is choosen
			for i := range webhooks { // loop true webhooks
				if webhooks[i].ID == parts[3] { // check if choosen id is in webhooks
					err = json.NewEncoder(w).Encode(webhooks[i]) // encode choosen webhook
					if err != nil {
						http.Error(w, "Something went wrong"+err.Error(), http.StatusInternalServerError)
					}

				}
			}

		} else {
			err = json.NewEncoder(w).Encode(webhooks) // encode all webhooks
			if err != nil {
				http.Error(w, "Something went wrong"+err.Error(), http.StatusInternalServerError)
			}
		}

	case http.MethodDelete:
		Wh := Webhook{}

		err := json.NewDecoder(r.Body).Decode(&Wh) //decode to webhook
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = DBDelete(Wh.ID, WebhooksCollection, w) //Deletes webhook from id
		if err != nil {
			fmt.Println("Error: ", err)
		}

	default:
		http.Error(w, "Method is invalid "+r.Method, http.StatusBadRequest)
	}

}
