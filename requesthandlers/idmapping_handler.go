package handlers

import (
	"net/http"
	"encoding/json"
)


type Message struct {
	Message string `json:"message"`
}


func IdMappingHandler(w http.ResponseWriter, r *http.Request) {

	msg := Message{Message:"Use POST for ID Mapping."}

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(msg)
	} else if r.Method == "POST" {
		var e interface{}

		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			if err := json.NewEncoder(w).Encode(e); err != nil {
				http.Error(w, err.Error(), 500)
			}
		}

	} else {
		http.Error(w, "ID Mapping request method must be POST.", 405)
	}
}
