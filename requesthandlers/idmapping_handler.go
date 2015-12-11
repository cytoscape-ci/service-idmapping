package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cytoscape-ci/service-go/idmapper"
	"fmt"
)


type Message struct {
	Message string `json:"message"`
}

var idMapper idmapper.IdMapper

func init() {
	resourceFile := "./data/idmapping.tsv"
	idMapper = idmapper.NewIdMapper(resourceFile)
	fmt.Println("############# Loaded.")
}


func IdMappingHandler(w http.ResponseWriter, r *http.Request) {

	msg := Message{Message:"Use POST for ID Mapping."}

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(msg)
	} else if r.Method == "POST" {

		var e map[string][]string

		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), 500)
		} else {

			ids := e["ids"]
			result := idMapper.Map(ids)


			if err := json.NewEncoder(w).Encode(result); err != nil {
				http.Error(w, err.Error(), 500)
			}
		}

	} else {
		http.Error(w, "ID Mapping request method must be POST.", 405)
	}
}
