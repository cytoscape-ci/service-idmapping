package idmapper

type Mapper interface {

	Map(values []string, filter []string) MappingResult

}
