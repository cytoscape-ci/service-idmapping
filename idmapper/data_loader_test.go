package idmapper


import (
	"testing"
	"fmt"
)

func TestLoad(t *testing.T) {

	fmt.Println("Test start---------------")

	resourceFile := "../data/idmapping.tsv"
	table := Load(resourceFile)

	numRecords := len(table.Entrez2Symbol)

	if numRecords != 190475 {
		t.Fail()
	}
}