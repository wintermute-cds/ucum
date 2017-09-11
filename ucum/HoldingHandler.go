package ucum

type HoldingHandler struct{
	Code string
	Units string
	Value *Decimal
}

func ( c * HoldingHandler)GetCode() string{
	return c.Code
}

func ( c * HoldingHandler)GetUnits() string{
	return c.Units
}

func ( c * HoldingHandler)GetValue() *Decimal{
	return c.Value
}

func NewHoldingHandler(code, units string, value *Decimal)*HoldingHandler{
	result := &HoldingHandler{}
	result.Code = code
	result.Units = units
	result.Value = value
	return result
}

