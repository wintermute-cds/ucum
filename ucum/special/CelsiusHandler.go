package special

import "UCUM_Golang/ucum"

type CelsiusHandler struct{
}

func ( c * CelsiusHandler)GetCode() string{
	return "Cel"
}

func ( c * CelsiusHandler)GetUnits() string{
	return "K"
}

func ( c * CelsiusHandler)GetValue() *ucum.Decimal{
	return ucum.One()
}

