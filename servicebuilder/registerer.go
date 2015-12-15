package servicebuilder

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"
	"log"
	"time"
)

const (
	BodyType = "application/json"

	defAgentUrl = "http://localhost:8080/registration"

	interval = 10
	retryMax = 10
)


type Instance struct {
	Location string `json:"location"`
	Capacity interface{}    `json:"capacity"`
}

type Registration struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	Instances []Instance `json:"instances"`
}


// Register single service
func RegisterService(agentUrl string, reg *Registration) error {

	log.Println("Registration start...")

	if agentUrl == "" {
		agentUrl = defAgentUrl
	}


	// Can post multiple services.
	var regs []*Registration
	regs = append(regs, reg)

	retryCount := 0

	for retryCount < retryMax {
		result, err := RegisterServices(agentUrl, regs)

		if err == nil {
			log.Println("Registered: ", result)
			break
		}

		log.Println("Retry in 10 sec...")
		time.Sleep(interval * time.Second)
	}

	return nil
}


// Register multiple services at once
func RegisterServices(agentUrl string, regs []*Registration) (result string, regError error) {

	if agentUrl == "" {
		agentUrl = defAgentUrl
	}

	// Encode JSON
	regJson, err := json.Marshal(regs)
	if err != nil {
		return "", errors.New("Invalid service registration info.")
	}

	return register(agentUrl, regJson)
}


func register(agentUrl string, regJson []byte) (result string, regError error) {

	log.Println("Registering service to the Agent...")

	// POST this service to submit agent
	res, err := http.Post(agentUrl, BodyType, bytes.NewReader(regJson))

	if err == nil {
		defer res.Body.Close()

		resBody, err := ioutil.ReadAll(res.Body)

		if err == nil {
			return string(resBody), nil
		} else {
			return string(resBody), errors.New("Could not read contents of response.")
		}
	} else {
		return "", errors.New("Could not complete registration.")
	}
}
