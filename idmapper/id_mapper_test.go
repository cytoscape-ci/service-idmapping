package idmapper

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
)

var mapper IdMapper

func TestMain(m *testing.M) {
	resourceFile := "../data/"
	mapper = NewIdMapper(resourceFile)
	code := m.Run()
	os.Exit(code)
}

func TestMapping(t *testing.T) {
	fmt.Println("ID Mapper test start...")

	idSet1 := []string{"ZCRB1", "474225", "foobar"}
	idSet2 := []string{"BRCA1", "BRCA2", "TP53", "rad5"}

	result1 := mapper.Map(idSet1)
	result2 := mapper.Map(idSet2)

	json.NewEncoder(os.Stdout).Encode(result1)
	json.NewEncoder(os.Stdout).Encode(result2)

	numMatch1 := len(result1.Matched)
	numUnmatch1 := len(result1.Unmatched)

	numMatch2 := len(result2.Matched)
	numUnmatch2 := len(result2.Unmatched)


	testCount(t, 3, numMatch1)
	testCount(t, 1, numUnmatch1)
	testCount(t, 7, numMatch2)
	testCount(t, 0, numUnmatch2)
}


func testCount(t *testing.T, expected int, actual int) {
	if actual != expected {
		t.Error("Expected: ", expected, " actual: ", actual)
		t.Fail()
	} else {
		t.Log("Success")
	}
}
