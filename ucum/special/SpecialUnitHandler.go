package special

import "UCUM_Golang/ucum"

type SpecialUnitHandlerer interface{
	GetCode() string
	GetUnits() string
	GetValue() *ucum.Decimal
}
