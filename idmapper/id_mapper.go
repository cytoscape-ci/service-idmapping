package idmapper


const(
	defaultSpecies = "homo sapiens"
)

type IdMapper struct {

	Table *ConversionTable
	Species string
}

func NewIdMapper(resourceFile string) (mapper IdMapper) {
	mappingTable := Load(resourceFile)
	return IdMapper{Table: mappingTable, Species:defaultSpecies}
}

func (mapper IdMapper) Map(values []string) MappingResult {

	speciesName := mapper.Species
	if speciesName == "" {
		speciesName = defaultSpecies
	}

	var matches []MappingEntry

	unmatched := make([]string, 0)

	for _, id := range values {
		found := false

		for _, table :=range mapper.Table.MappingTable {
			match, exists := table[id]

			if exists {
				matches = append(matches, match)
				found = true
				break
			}
		}

		if !found  {
			unmatched = append(unmatched, id)
		}
	}

	result := MappingResult{Matched:matches, Unmatched:unmatched}

	return result
}
