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
	return s
}

func MaxInt(a, b int)int{
	if a>=b {
		return a
	}else {
		return b
	}
}

func MinInt(a, b int)int{
	if a<=b {
		return a
	}else {
		return b
	}
}

func IsAsciiChar(ch rune) bool {
	return ch >= ' ' && ch <= '~';
}
