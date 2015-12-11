package idmapper

type Mapper interface {

	Map(values []string) MappingResult

}
