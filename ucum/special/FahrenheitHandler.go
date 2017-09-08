package special

import "UCUM_Golang/ucum"

type FahrenheitHandler struct{
}

func ( c * FahrenheitHandler)GetCode() string{
	return "[degF]"
}

func ( c * FahrenheitHandler)GetUnits() string{
	return "K"
}

func ( c * FahrenheitHandler)GetValue() *ucum.Decimal{
	d := ucum.NewDecimal("5")
	return d.Divide(ucum.NewDecimal("9"))
}

