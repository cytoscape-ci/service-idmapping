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

type ConversionTable struct {
	Entrez2Symbol map[string]string
}

func Load(mappingFile string) *ConversionTable {
	fmt.Println("Loading data into memory...")
	checkMemory()

	f, err := os.Open(mappingFile)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	r.Comma = rune('\t')


	entrez := make(map[string]string)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		processLine(record, entrez)
	}

	conv := ConversionTable{Entrez2Symbol:entrez}
	checkMemory()


	fmt.Println("Done!")
	return &conv
}

func processLine(line []string, table map[string]string) {
//	fmt.Println(len(line))

	table[line[0]] = strings.Join(line, ", ")
}

func checkMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Println("Allocated: ", mem.Alloc / 1024)
	log.Println("Total: ", mem.TotalAlloc / 1024)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
}
