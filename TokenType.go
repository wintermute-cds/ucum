package ucum

type TokenType int

const (
	_ TokenType = iota
	NONE
	NUMBER
	SYMBOL
	SOLIDUS
	PERIOD
	OPEN
	CLOSE
	ANNOTATION
)
