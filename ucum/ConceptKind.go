package ucum

type ConceptKind int

const(
	PREFIX ConceptKind = iota
	BASEUNIT
	UNIT
)