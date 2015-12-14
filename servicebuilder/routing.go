package servicebuilder

import (
	"net/http"
	"log"
	req "github.com/cytoscape-ci/service-go/requesthandlers"
	"github.com/rs/cors"
	"strconv"
)


func StartServer(portNumber int) (err error) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", req.StatusHandler)
	mux.HandleFunc("/map", req.IdMappingHandler)

	handler := cors.Default().Handler(mux)

	log.Println("Serving API on port ", portNumber)

	portNumStr := strconv.Itoa(3000)

	err = http.ListenAndServe(":" + portNumStr, handler)

	if err != nil {
		log.Fatal("Could not start API server: ", err)
	}

	return nil
}
