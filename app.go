package main

import (
	"log"
	"os"
	"strconv"
	"flag"
	elsa "github.com/cytoscape-ci/elsa-client"
	handlers "github.com/cytoscape-ci/service-idmapping/requesthandlers"
)



func main() {
	reg := elsa.NewRegistrationFromCommandline()

	elsaLocation := flag.String("agent", "http://192.168.99.100:8080/registration", "Agent URL")
	flag.Parse()

	var servicePort int

	portFlag := flag.Lookup("port")

	if portFlag == nil {
		servicePort = elsa.DefPort
	} else {
		var portErr error
		servicePort, portErr = strconv.Atoi(portFlag.Value.String())
		if portErr != nil {
			log.Fatal("Could not start API server: ", portErr.Error())
			os.Exit(1)
		}
	}

	// Asynchronously register this service to Submit Agent
	go elsa.RegisterService(*elsaLocation, reg, elsa.RetrySetting{})

	// Start API server
	serverErr := handlers.StartServer(servicePort)

	if serverErr!= nil {
		log.Fatal("Could not start API server: ", serverErr.Error())
		os.Exit(1)
	}
}
