package cravings

import (
	"encoding/json"
	"net/http"
	"time"
)

type Status struct {
	Edemam   int     `json:"edemam"`
	Database int     `json:"database"`
	Uptime   float64 `json:"uptime"`
	Version  string  `json:"version"`
}

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	var S Status

	resp, err := http.Get("https://api.edamam.com/api/nutrition-details")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	S.Edemam = resp.StatusCode

	resp, err = http.Get("https://firebase.google.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	S.Database = resp.StatusCode

	elapse := time.Since(StartTime)
	S.Uptime = elapse.Seconds()

	S.Version = "v1"

	http.Header.Add(w.Header(), "Content-Type", "application/json") // makes the print look good

	json.NewEncoder(w).Encode(S)
}
