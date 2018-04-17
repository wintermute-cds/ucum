package ucum

type FahrenheitHandler struct {
}

func (c *FahrenheitHandler) GetCode() string {
	return "[degF]"
}

func (c *FahrenheitHandler) GetUnits() string {
	return "K"
}

func (c *FahrenheitHandler) GetValue() *Decimal {
	d5, _ := NewDecimal("5")
	d9, _ := NewDecimal("9")
	d := d5.Divide(d9)
	return d
}
