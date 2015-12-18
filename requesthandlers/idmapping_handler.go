package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cytoscape-ci/service-go/idmapper"
	"log"
	"errors"
	"strconv"
)


const (
	GET = "GET"
	POST = "POST"

	ids = "ids"
	idTypes = "idTypes"
)

type Message struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Error string `json:"error,omitempty"`
}

var idMapper idmapper.IdMapper


func init() {
	// Initialize mapping table using resource file.
	resourceFileDir := "./data/"
	idMapper = idmapper.NewIdMapper(resourceFileDir)
	log.Println("Mapping table loaded.")
}


func IdMappingHandler(w http.ResponseWriter, r *http.Request) {

	method := r.Method
	switch method {
	case POST:
		post(w, r)
	default:
		unsupported(w, r)
	}
}


func post(w http.ResponseWriter, r *http.Request) {
	var e map[string][]string

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		code := 422
		res := getErrorMsg(code, "Could not decode your input data.  You should provide 'ids' and optional 'idTypes' list of strings.", err)
		http.Error(w, res, code)
	} else {
		mapping(w, e)
	}
}

func mapping(w http.ResponseWriter, e map[string][]string) {
	ids, exists := e[ids]
	filter, _ := e[idTypes]

	if !exists {
		code := 400
		res := getErrorMsg(code, "Could not decode your input data.", errors.New("ids field missing."))
		http.Error(w, res, code)
		return
	}

	inputSize := len(ids)

	// Call actual mapper service
	result := idMapper.Map(ids, filter)

	unmatchedSize := len(result.Unmatched)

	// No match found.  Return 404.
	if inputSize == unmatchedSize {
		code := 404
		res := getErrorMsg(code, "Could not find any match.", errors.New("There is no matching in the database."))
		http.Error(w, res, code)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		code := 500
		res := getErrorMsg(code, "Could not parse your mapping result.  This may be a bug in this service.", err)
		http.Error(w, res, code)
	} else {
		log.Println("Success: Mapping finished for ", strconv.Itoa(len(ids)), " IDs.")
	}
}


func unsupported(w http.ResponseWriter, r *http.Request) {
	code := 405
	res := getErrorMsg(code, "You need to POST your data to use this service.",
		errors.New("Invalid HTTP method used: " + r.Method))
	log.Println("Unsupported method call from: ", r.RemoteAddr)
	http.Error(w, res, code)
}


func getErrorMsg(code int, msg string, err error) (jsonMsg string) {

	message := Message{
		Code: code,
		Message: msg,
	}

	if err != nil {
		message.Error = err.Error()
	}

	result, err := json.Marshal(message)

	if err !=nil {
		// TODO: What should I return?
	}

	return string(result)
}
