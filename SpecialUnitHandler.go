package ucum


import "github.com/wintermute-cds/ucum/decimal"

type SpecialUnitHandlerer interface {
	GetCode() string
	GetUnits() string
	GetValue() decimal.Decimal
}
