package main


import (
	"log"
	"fmt"
	"os"
	builder "github.com/cytoscape-ci/service-go/servicebuilder"
	"net"
	"strconv"
	"flag"
)


const (
	serviceName = "go-service"
	defPort = 3000

	defAgentUrl = "http://localhost:8080/register"
	defServ = "127.0.0.1"
)


var reg *builder.Registration
var agent *string


func init() {
	// Initialize logging
	fmt.Println("Initializing API Server...")

	// Init logger to use syslog
//	logger, err := syslog.New(
//		syslog.LOG_NOTICE|syslog.LOG_USER, serviceName)
//	if err != nil {
//		panic(err)
//	}
//
//	log.SetOutput(logger)
//	log.Println("* Logging start using syslog...")

	// Parse parameters
	reg, agent = buildReg()
}


func buildReg() (reg *builder.Registration, agent *string) {
	// This is required!
	name := flag.String("service", "", "Service name")
	cap := flag.Int("cap", 4, "Number of instances")
	ver := flag.String("version", "v1", "API version")
	agentUrl := flag.String("agent", "http://192.168.99.100:8080/registration", "Submit Agent Location")
	loc := flag.String("location", "192.168.99.100", "This API server Location")

	flag.Parse()

	if *name == "" {
		log.Panic("Missing service endpoint name: You must provide it as -service param.")
		os.Exit(1)
	}

	var myLoc string
	if *loc == "" {
		myLoc = getAddress()
	} else {
		myLoc = *loc
	}

	myUrl := myLoc + ":" + strconv.Itoa(defPort)
	log.Println("Service API Location:", myUrl)

	instance := builder.Instance{Capacity: *cap, Location:myUrl}
	reg = &builder.Registration {
		Service: *name,
		Version: *ver,
		Instances: []builder.Instance{instance},
	}

	return reg, agentUrl
}


func main() {
	go builder.RegisterService(*agent, reg)

//	if err != nil {
//		log.Println(err)
//		log.Println("Could not register service.  Running in stand-alone mode.")
//	} else {
//		log.Println("Service registered to Agent")
//	}

	// Start API server
	err := builder.StartServer(defPort)

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
	return defServ
}
