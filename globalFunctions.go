package cravings

import "net/http"

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
