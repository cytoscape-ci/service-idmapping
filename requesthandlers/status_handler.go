package handlers

import (
	"net/http"
	"encoding/json"
)


type Status struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Description string `json:"description"`
}


func StatusHandler(w http.ResponseWriter, r *http.Request) {

	serviceStatus := Status{Name:"ID Mapping Service", Version:"v1", Description:"ID Mapper."}

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(serviceStatus)
	} else {
		http.Error(w, "Request method must be GET.", 405)
	}
}
