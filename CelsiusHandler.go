package ucum


import "github.com/wintermute-cds/ucum/decimal"

type CelsiusHandler struct {
}

func (c *CelsiusHandler) GetCode() string {
	return "Cel"
}

func (c *CelsiusHandler) GetUnits() string {
	return "K"
}

func (c *CelsiusHandler) GetValue() decimal.Decimal {
	d, _ := decimal.NewFromString("1")
	return d
}
