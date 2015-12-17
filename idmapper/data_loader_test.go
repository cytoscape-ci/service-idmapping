package idmapper

import (
	"testing"
)

func TestLoad(t *testing.T) {

	t.Log("Data Loading test start")

	resourceFile := "../data/"
	tableMap := Load(resourceFile)

	// Check number of tables (one / species)
	numMap := len(tableMap)

	exp := 4 // Human, Fly, Yeast, and Mouse.
	if numMap != exp {
		t.Error("Expected: ", exp, " actual: ", numMap)
		t.Fail()
	} else {
		t.Log("Success species table")
	}


	mapper := tableMap["human"].MappingTable
	numKeys := len(mapper)

	exp = len(TargetColumns)
	if numKeys != exp {
		t.Error("Expected: ", exp, " actual: ", numKeys)
		t.Fail()
	} else {
		t.Log("Success key column count for human")
	}
}