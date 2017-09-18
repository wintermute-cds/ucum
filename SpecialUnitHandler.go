package ucum

type SpecialUnitHandlerer interface {
	GetCode() string
	GetUnits() string
	GetValue() Decimal
}
