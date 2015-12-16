package idmapper


import (
	"testing"
	"fmt"
)

func TestLoad(t *testing.T) {

	fmt.Println("Test start---------------")

	resourceFile := "../data/"
	tableMap := Load(resourceFile)

	mapper := tableMap["human"].MappingTable
	numKeys := len(mapper)


	for key, _ := range mapper {
		t.Log("Key# = ", key)
	}

	if numKeys != 4 {
		t.Log("Num keys = ", numKeys)
		t.Fail()
	}

	i := 0
	for key, value :=range mapper["Symbol"] {
		fmt.Println(key, " = ", value)
		if i > 100 {
			break
		}
		i++
	}
}