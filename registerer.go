package main
import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"os"
	"net"
)

type Registration struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Location string `json:"location"`
	Version  string `json:"version"`
}

func registerService() {

	reg := &Registration{
		Name:     "test-service",
		Capacity: 4,
		Location: getAddress(),
		Version: "v1",
	}

	var regs []*Registration

	regs = append(regs, reg)

	regJson, err := json.Marshal(regs)

	if err == nil {
		res, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewReader(regJson))
		if err == nil {
			defer res.Body.Close()
			contents, err := ioutil.ReadAll(res.Body)
			if err == nil {
				fmt.Println(string(contents))
			} else {
				panic("Could not read contents of response.")
			}
		} else {
			panic("Could not complete request.")
		}
	} else {
		panic("Could not read struct into json.")
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
	return "127.0.0.1"
}
