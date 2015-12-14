package main


import (
	"log"
	"fmt"
	"log/syslog"
	"os"
	builder "github.com/cytoscape-ci/service-go/servicebuilder"
	"net"
	"strconv"
)


const (
	serviceName = "go-service"
	defPort = 3000

	defAgentUrl = "http://localhost:8080/register"
	defServ = "127.0.0.1"
)

var syslogWriter *syslog.Writer


func init() {
	// Initialize logging
	fmt.Println("Initializing API Server...")

	// Init logger to use syslog
	logger, err := syslog.New(
		syslog.LOG_NOTICE|syslog.LOG_USER, serviceName)
	if err != nil {
		panic(err)
	}

	syslogWriter = logger
	log.SetOutput(syslogWriter)

	log.Println("* Logging start using syslog...")
}


func main() {

	myLoc := getAddress()
	myUrl := myLoc + ":" + strconv.Itoa(defPort)
	log.Println("Service API Location:", myUrl)

	// Register service
	// TODO: Use config file instead.

	instance := builder.Instance{Capacity:4, Location:myUrl}
	reg := &builder.Registration {
		Service:     "idmapping",
		Version: "v1",
		Instances: []builder.Instance{instance},
	}

	err := builder.RegisterService("http://192.168.99.100:8080/registration", reg)
	if err != nil {
		log.Println(err)
		log.Println("Could not register service.  Running in stand-alone mode.")
	} else {
		log.Println("Service registered to Agent")
	}


	// Start API server
	err = builder.StartServer(defPort)

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
