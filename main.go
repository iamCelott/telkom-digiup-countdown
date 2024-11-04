package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CountdownResponse struct {
	Hour int `json:"hour"`
	Minute int `json:"minute"`
	Second  int    `json:"second"`
	Message string `json:"message,omitempty"`
}
	
func handleIndex(w http.ResponseWriter, r * http.Request) {
	http.ServeFile(w,r, "views/index.html")
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), 17, 30, 0, 0, now.Location())
	
	endTime := time.Date(2024, time.November, 4, 17, 30, 0, 0, now.Location())

	response := CountdownResponse{}

	if now.After(endTime) {
		response.Message = "The target day has passed."
	} else {
		if now.After(targetTime) {
			nextDay := now.Add(24 * time.Hour)
			targetTime = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 13, 0, 0, 0, nextDay.Location())
		}

		duration := targetTime.Sub(now)

		if duration > 0 {
			response.Hour = int(duration.Hours()) % 24
			response.Minute = int(duration.Minutes()) % 60
			response.Second = int(duration.Seconds()) % 60
		} else {
			response.Hour = 0
			response.Minute = 0
			response.Second = 0
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main () {
	mux := http.NewServeMux()
	var port = "8000"

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/time", handleTime)
	
	fmt.Printf("Server Started at Port %s\n", port)
	http.ListenAndServe(":" + port, mux)
}