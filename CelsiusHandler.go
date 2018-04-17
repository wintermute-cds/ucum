package ucum

type CelsiusHandler struct {
}

func (c *CelsiusHandler) GetCode() string {
	return "Cel"
}

func (c *CelsiusHandler) GetUnits() string {
	return "K"
}

func (c *CelsiusHandler) GetValue() *Decimal {
	d, _ := NewDecimal("1")
	return d
}
