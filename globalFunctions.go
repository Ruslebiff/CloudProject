package cravings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func doRequest(url string, c *http.Client, w http.ResponseWriter) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}

	resp, err := c.Do(req)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}

	return resp
}

//QueryGet func to read  variable for link
func QueryGet(s string, w http.ResponseWriter, r *http.Request) string {

	test := r.URL.Query().Get(s) // gets app key or app id
	if test == "" {              // check if it is empty
		fmt.Fprintln(w, s+" is missing")
	}
	return test

}

// CallURL post webhooks to webhooks.site
func CallURL(event string, s interface{}) {

	webhooks, err := DBReadAllWebhooks()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for i := range webhooks {
		if webhooks[i].Event == event {
			var request = s

			requestBody, err := json.Marshal(request)
			if err != nil {
				fmt.Println("Can not encode: " + err.Error())
			}

			fmt.Println("Attempting invoation of URL " + webhooks[i].URL + "...")

			resp, err := http.Post(webhooks[i].URL, "json", bytes.NewReader([]byte(requestBody)))
			if err != nil {
				fmt.Println("Error in HTTP request: " + err.Error())
			}

			response, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Something vent wrong: " + err.Error())
			}

			fmt.Println("Webhook body: " + string(response))

		}

	}

}
