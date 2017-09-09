package ucum

type ConceptKind int

const(
	_ ConceptKind = iota
	PREFIX
	BASEUNIT
	UNIT
)