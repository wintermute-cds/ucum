package ucum

import "strconv"

func IsDecimal(value string)bool{
	if value == ""{
		return false
	}
	_,err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return true
}

func IsInteger(value string)bool{
	if value == ""{
		return false
	}
	_,err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return true
}

func PadLeft(src string, c rune, l int)string{
	s := ""
	for i := 0; i < l - len(src); i++ {
		s = s + string(c)
	}
}

func max(a, b int)int{
	if a>=b {
		return a
	}else {
		return b
	}
}

func min(a, b int)int{
	if a<=b {
		return a
	}else {
		return b
	}
}
