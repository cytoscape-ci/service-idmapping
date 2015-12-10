package main

import (
	"log"
	"fmt"
	"net/http"
	mapper "github.com/cytoscape-ci/service-go/idmapper"
	req "github.com/cytoscape-ci/service-go/requesthandlers"
	"github.com/rs/cors"
)


func init() {
	fmt.Println("Initializing API Server...")
}


func main() {

	table := mapper.Load("./data/idmapping.tsv")

	fmt.Println("Table size: ", len(table.Entrez2Symbol))

	mux := http.NewServeMux()
	mux.HandleFunc("/", req.StatusHandler)
	mux.HandleFunc("/idmapping", req.IdMappingHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("\n========== Serving API on port 3000 =========")

	for _, val := range table.Entrez2Symbol {
		fmt.Println("Line = ", val)
	}

	err := http.ListenAndServe(":3000", handler)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
