package main

import (
	"log"
	"os"
	builder "github.com/cytoscape-ci/service-go/servicebuilder"
	"net"
	"strconv"
	"flag"
)



func buildReg() (reg *builder.Registration, agent *string, port *int) {
	// This is required!
	name := flag.String("id", "", "Service name")
	cap := flag.Int("cap", 4, "Number of instances")
	ver := flag.String("ver", "v1", "API version")
	agentUrl := flag.String("agent", "http://192.168.99.100:8080/registration", "Submit Agent Location")
	ip := flag.String("ip", "", "This API server IP Address")
	port = flag.Int("port", 3000, "This API server IP Address")

	flag.Parse()

	if *name == "" {
		log.Panic("Missing service endpoint name: You must provide it with '-service' param.")
		os.Exit(1)
	}

	var myLoc string
	if *ip == "" {
		myLoc = getAddress()
	} else {
		myLoc = *ip
	}

	myUrl := myLoc + ":" + strconv.Itoa(*port)
	log.Println("Service API Location:", myUrl)

	instance := builder.Instance{Capacity: *cap, Location:myUrl}
	reg = &builder.Registration{
		Service: *name,
		Version: *ver,
		Instances: []builder.Instance{instance},
	}

	return reg, agentUrl, port
}


func main() {
	// Parse parameters
	reg, agentUrl, port := buildReg()

	// Asynchronously register this service
	go builder.RegisterService(*agentUrl, reg)

	// Start API server
	err := builder.StartServer(*port)

	if err != nil {
		log.Fatal("Could not start API server: ", err.Error())
		os.Exit(1)
	}
}


func getAddress() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	panic("Could not find service IP address")
}
