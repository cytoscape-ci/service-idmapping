package idmapper

import (
	"testing"
	"fmt"
	"os"
	"encoding/json"
)

var mapper IdMapper

func TestMain(m *testing.M) {
	resourceFile := "../data/idmapping.tsv"
	mapper = NewIdMapper(resourceFile)

	code := m.Run()
	os.Exit(code)
}

//func TestMapping(t *testing.T) {
//
//	fmt.Println("Id mapping Test start---------------")
//
//
//	// This is how to set species.
//	mapper.Species = "homo sapiens"
//
//	idSet1 := []string{"ZCRB1", "474225", "foobar"}
//	idSet2 := []string{"BRCA1", "BRCA2", "TP53"}
//
//	result1 := mapper.Map(idSet1)
//	result2 := mapper.Map(idSet2)
//
//	json.NewEncoder(os.Stdout).Encode(result1)
//	json.NewEncoder(os.Stdout).Encode(result2)
//
//	t.Log("Res1 Matched = ", result1.Matched)
//	t.Log("Res1 Unmatched = ", result1.Unmatched)
//
//	t.Log("Res2 Matched = ", result2.Matched)
//	t.Log("Res1 Unmatched = ", result2.Unmatched)
//}

func BenchmarkMapper(b *testing.B) {

	fmt.Println("\n----------- Performance bench start---------------")

	idSet1 := []string{"ZCRB1", "474225", "foobar"}

	for i := 0; i < b.N; i++ {
		mapper.Map(idSet1)

		result1 := mapper.Map(idSet1)
//		json.NewEncoder("").Encode(result1)
		json.Marshal(result1)
	}
}
