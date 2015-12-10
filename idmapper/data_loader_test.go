package idmapper


import (
	"testing"
	"fmt"
)

func TestLoad(t *testing.T) {

	fmt.Println("Test start---------------")

	resourceFile := "../data/idmapping.tsv"
	Load(resourceFile)

}