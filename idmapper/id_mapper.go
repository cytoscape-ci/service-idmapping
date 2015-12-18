package idmapper


import "strings"

const(
	defaultSpecies = "homo sapiens"
)

type IdMapper struct {
	Tables map[string]*ConversionTable
	Species string
}

func NewIdMapper(resourceFile string) (mapper IdMapper) {
	mappingTables := Load(resourceFile)
	return IdMapper{Tables: mappingTables}
}

func (mapper IdMapper) Map(values []string, filter []string) MappingResult {

//	speciesName := mapper.Species
//	if speciesName == "" {
//		speciesName = defaultSpecies
//	}

	var matches []MappingEntry

	unmatched := make([]string, 0)

	for _, id := range values {
		found := false

		// For each species, find match (ignore case, always upper)
		for species, value := range mapper.Tables {
			tbl := value

			for _, table := range tbl.MappingTable {
				match, exists := table[strings.ToUpper(id)]

				if exists {
					match.Species = species
					match.In = id

					if len(filter) > 0 {
						matches = append(matches, *filterResult(match, filter))
					} else {
						matches = append(matches, *match)
					}

					found = true
					break
				}
			}
		}

		if !found  {
			unmatched = append(unmatched, id)
		}
	}
	result := MappingResult{Matched:matches, Unmatched:unmatched}

	return result
}

func filterResult(entry *MappingEntry, filter []string) (*MappingEntry) {


	filtered := MappingEntry{Species: entry.Species, In: entry.In, InType: entry.InType}
	matches := make(map[string]interface{})

	for _, idType := range filter {
		val := entry.Matches[idType]
		if val != nil {
			matches[idType] = val
		}
	}

	filtered.Matches = matches

	return &filtered
}
