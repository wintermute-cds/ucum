package ucum

type ConceptKind int

const (
	PREFIX ConceptKind = iota + 1
	BASEUNIT
	UNIT
)
