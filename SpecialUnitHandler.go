package ucum


import "github.com/bertverhees/ucum/decimal"

type SpecialUnitHandlerer interface {
	GetCode() string
	GetUnits() string
	GetValue() decimal.Decimal
}
