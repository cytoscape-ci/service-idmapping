package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"errors"
	"strconv"
	"github.com/cytoscape-ci/service-idmapping/idmapper"
)


func LabelGeneratorHandler(w http.ResponseWriter, r *http.Request) {

	method := r.Method
	switch method {
	case POST:
		generateLabel(w, r)
	default:
		unsupported(w, r)
	}
}


func generateLabel(w http.ResponseWriter, r *http.Request) {
	var e map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		code := 422
		res := getErrorMsg(code, "Could not decode your input data.  You should pass name-to-@id map.", err)
		http.Error(w, res, code)
	} else {
		getLabels(w, e)
	}
}


func getLabels(w http.ResponseWriter, e map[string]interface{}) {

	idlist := e["ids"]

	dummyMap := make(map[string]bool)

	var ids []string

	switch idlist.(type) {

	case []interface{} :
		for _, val := range idlist.([]interface{}) {
			ids = append(ids, val.(string))
			dummyMap[val.(string)] = true
		}
	default:
		panic("NOT array")
	}

	sp := e["species"].(string)

	filter := []string{"Symbol"}

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


	labels := createLabelMappings(result, dummyMap, sp)

	if err := json.NewEncoder(w).Encode(labels); err != nil {
		code := 500
		res := getErrorMsg(code, "Could not parse your mapping result.  This may be a bug in this service.", err)
		http.Error(w, res, code)
	} else {
		log.Println("Success: Mapping finished for ", strconv.Itoa(len(ids)), " IDs.")
	}
}


func createLabelMappings(res idmapper.MappingResult, ids map[string]bool, species string) (table map[string]string) {

	matches := res.Matched

	table = make(map[string]string)

	for _, entry := range matches {
		if entry.Species != species {
			continue
		}

		_, exists := ids[entry.In]

		if exists {
			match := entry.Matches
			hit := match["Symbol"]
			table[entry.In] = hit.(string)
		}
	}

	return table
}