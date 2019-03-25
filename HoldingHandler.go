package ucum


import "github.com/wintermute-cds/ucum/decimal"

type HoldingHandler struct {
	Code  string
	Units string
	Value decimal.Decimal
}

func (c *HoldingHandler) GetCode() string {
	return c.Code
}

func (c *HoldingHandler) GetUnits() string {
	return c.Units
}

func (c *HoldingHandler) GetValue() decimal.Decimal {
	return c.Value
}

func NewHoldingHandler(code, units string, value decimal.Decimal) *HoldingHandler {
	result := &HoldingHandler{}
	result.Code = code
	result.Units = units
	result.Value = value
	return result
}
