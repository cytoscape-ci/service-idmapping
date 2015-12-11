package main

import (
	"log"
	"fmt"
	"net/http"
	req "github.com/cytoscape-ci/service-go/requesthandlers"
	"github.com/rs/cors"
)


func init() {
	fmt.Println("Initializing API Server...")
}


func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", req.StatusHandler)
	mux.HandleFunc("/idmapping", req.IdMappingHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("\n========== Serving API on port 3000 =========")


	err := http.ListenAndServe(":3000", handler)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
