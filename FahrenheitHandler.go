package ucum



import "github.com/bertverhees/ucum/decimal"

type FahrenheitHandler struct {
}

func (c *FahrenheitHandler) GetCode() string {
	return "[degF]"
}

func (c *FahrenheitHandler) GetUnits() string {
	return "K"
}

func (c *FahrenheitHandler) GetValue() decimal.Decimal {
	d5, _ := decimal.NewFromString("5")
	d9, _ := decimal.NewFromString("9")
	d := d5.Div(d9)
	return d
}
