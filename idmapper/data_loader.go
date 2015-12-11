package idmapper

import (
	"fmt"
	"encoding/csv"
	"log"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	listSeparator = "; "
)

type ConversionTable struct {
	MappingTable map[string]map[string]MappingEntry
}


func Load(mappingFile string) *ConversionTable {
	fmt.Println("Loading data into memory...")

	targetColumns := []string{"UniProtKB-AC", "UniProtKB-ID", "GeneID", "Ensembl", "Symbol"}

	f, err := os.Open(mappingFile)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	r.Comma = rune('\t')

	var baseTable [][]string

	i := 1

	var cols []string
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if i==1 {
			// This is a header line
			cols = record
			i++
			continue
		}

		buildBaseTable(record, &baseTable)
		i++
	}


	mappingTable := createMap(targetColumns, cols, &baseTable)

	conv := ConversionTable{ MappingTable:mappingTable }
//	checkMemory()


	fmt.Println("Done!")
	baseTable = nil

	return &conv
}

func buildBaseTable(rec []string, table *[][]string) {
	*table = append(*table, rec)
}


func createMap(columns []string, allColumnNames []string, table *[][]string) map[string]map[string]MappingEntry {

	mappingTable := make(map[string]map[string]MappingEntry)

	indicies := getKeyIndicies(columns, allColumnNames)

	for i, idx := range indicies {

		key2rec := make(map[string]MappingEntry)

		for _, rec := range *table {
			key := rec[idx]
			columnType := allColumnNames[idx]

			if key != "" {
				key2rec[key] = createEntry(key, columnType, rec, allColumnNames)
			}
		}

		mappingTable[columns[i]] = key2rec
	}
	return mappingTable
}

func createEntry(id string, columnType string, rec []string, columnNames []string) (entry MappingEntry) {

	mapping := make(map[string]interface{})

	for idx, val := range rec {
		if val != "" && val != id && val != "-" {
			idType := columnNames[idx]

			if strings.Contains(val, listSeparator) {
				mapping[idType] = createListEntry(val)
			} else {
				mapping[idType] = val
			}

		}
	}

	return MappingEntry{In:id, InType:columnType, Matches:mapping}
}

func createListEntry(in string) []string {

	result := make([]string, 0)

	parts := strings.Split(in, listSeparator)
	for _, val := range parts {
		if val != "" && val != "-" {
			result = append(result, val)
		}
	}
	return result
}

func getKeyIndicies(keys []string, allColumns []string) []int {
	indicies := make([]int, len(keys))

	for idx, val := range keys {

		for i, columnName := range allColumns {
			if columnName == val {
				indicies[idx] = i
				break
			}
		}
	}
	return indicies
}


func checkMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Println("Allocated: ", mem.Alloc / 1024)
	log.Println("Total: ", mem.TotalAlloc / 1024)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
}
