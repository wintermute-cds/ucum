package special

import "UCUM_Golang/ucum"

type HoldingHandler struct{
	Code string
	Units string
	Value *ucum.Decimal
}

func ( c * HoldingHandler)GetCode() string{
	return c.Code
}

func ( c * HoldingHandler)GetUnits() string{
	return c.Units
}

func ( c * HoldingHandler)GetValue() *ucum.Decimal{
	return c.Value
}

func NewHoldingHandler(code, units string, value *ucum.Decimal)*HoldingHandler{
	result := &HoldingHandler{}
	result.Code = code
	result.Units = units
	result.Value = value
	return result
}

