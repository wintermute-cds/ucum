package ucum

import (
	"strings"
	"unicode"
	"fmt"
	"strconv"
)

type Decimal struct {
	Decimal int
	Digits string
	Negative bool
	Precision int
	Scientific bool
}

func NewDecimal(value string, precision int)(*Decimal, error){
	d := &Decimal{}
	value = strings.ToLower(value)
	var err error
	if strings.Contains(value, "e"){
		err = d.setValueScientific(value)
	}else{
		err = d.setValueDecimal(value)
	}
	if err !=nil {
		return nil, err
	}
	d.Precision = precision
	return d,nil
}

func (d *Decimal)setValueDecimal(value string)error{
	var err error
	d.Scientific = false
	dec := -1
	d.Negative = strings.Index(value, "-") == 0
	if d.Negative {
		value = value[1:]
	}
	for {
		if !(strings.Index(value, "0") == 0 && len(value) > 1) {
			break
		}
	}
	for i,c := range value{
		if c == '.' && dec == -1 {
			dec = i
		}else if !unicode.IsDigit(c){
			return fmt.Errorf("'"+value+"'  is not a valid decimal")

		}
	}
	if dec == -1 {
		d.Precision = len(value)
		d.Decimal = len(value)
		d.Digits = value
	}else if dec == len(value) - 1 {
		return fmt.Errorf("'"+value+"'  is not a valid decimal")
	}else {
		d.Decimal = dec
		if d.allZeros(value, 1){
			d.Precision = len(value) - 1
		} else {
			d.Precision, err = d.countSignificants(value)
			if err!=nil {
				return err
			}
		}
		d.Digits, err = d.delete(value, d.Decimal, 1)
		if err!=nil {
			return err
		}
		if d.allZeros(d.Digits, 0){
			d.Precision++
		}else {
			tmp := d.Digits
			for _,c := range tmp{
				if c != '0' {
					break
				}
				d.Digits = d.Digits[1:]
				d.Decimal--
			}
		}
	}
	return nil
}

func (d *Decimal)allZeros(s string, start int)bool{
	for i,c := range s {
		if i >= start {
			if c != '0' {
				return false
			}
		}
	}
	return true
}

func (d *Decimal)countSignificants(value string)(int, error){
	i := strings.Index(value, ".")
	var err error
	if i > -1 {
		value, err = d.delete(value, i, 1)
		if err!=nil {
			return 0, err
		}
	}
	tmp := value
	for _,c := range tmp{
		if c != '0' {
			break
		}
		value = value[1:]
	}
	return len(value), nil
}

func (d *Decimal)delete(value string, offset, length int)(string, error){
	if length > len(value) {
		return "", fmt.Errorf("Length cannot be greater then the length of the value string")
	}
	if offset > len(value) {
		return "", fmt.Errorf("Offset cannot be greater then the length of the value string")
	}
	if offset == 0 {
		return value[length:], nil
	}else{
		return value[0:offset]+value[offset+length:], nil
	}
}

func (d *Decimal)setValueScientific(value string) (err error) {
	i := strings.Index(value, "e")
	s := value[0:i]
	e := value[i+1:]
	if s == "" || s == "-" || !IsDecimal(s){
		return fmt.Errorf("'"+value+"' is not a valid decimal (numeric)")
	}
	if e == "" || e == "-" || !IsInteger(e){
		return fmt.Errorf("'"+value+"' is not a valid decimal (exponent)")
	}
	d.setValueDecimal(s)
	d.Scientific = true
	//exponent
	if e[0] == '-' {
		i = 1
	} else {
		i = 0
	}
	for j,c := range e {
		if j >= i {
			if !unicode.IsDigit(c){
				return fmt.Errorf("'"+value+"' is not a valid decimal")
			}
		}
	}
	i,err = strconv.Atoi(e)
	if err !=nil{
		return err
	}
	d.Decimal = d.Decimal + i
	return nil
}

func (d *Decimal)stringMultiply(c rune, i int)string{
	return PadLeft("", c, i)
}

func (d *Decimal)insert(ins, value string, offset int)(string, error){
	if offset > len(value) {
		return "", fmt.Errorf("Offset cannot be greater then the length of the value string")
	}
	if offset == 0 {
		return ins + value,nil
	}else{
		return value[0:offset] + ins + value[offset:], nil
	}
}

func (d *Decimal)String()string{
	return d.AsDecimal()
}

func (d *Decimal)copy()*Decimal{
	r := &Decimal{}
	r.Precision = d.Precision
	r. Scientific = d.Scientific
	r.Negative = d.Negative
	r.Digits = d.Digits
	r.Decimal = d.Decimal
	return r
}

func Zero()(*Decimal){
	d := &Decimal{}
	d.setValueDecimal(strconv.Itoa(0))
	return d
}

func One()(*Decimal){
	d := &Decimal{}
	d.setValueDecimal(strconv.Itoa(1))
	return d
}

func (d *Decimal)IsZero()bool{
	return d.allZeros(d.Digits,0)
}

func (d *Decimal)IsOne()bool{
	one := One()
	return d.ComparesTo(one) == 0
}

func (d *Decimal)Equals(other *Decimal)bool{
	return d.ComparesTo(other) == 0
}

func (d *Decimal)ComparesTo(other *Decimal)int{
	if other == nil {
		return 0
	}

	if d.Negative && !other.Negative {
		return -1
	}else if !d.Negative && other.Negative {
		return 1
	}else {
		max := other.Decimal
		if d.Decimal >= other.Decimal {
			max = d.Decimal
		}
		s1 := d.stringMultiply('0', max - d.Decimal + 1) + d.Digits
		s2 := d.stringMultiply('0', max - other.Decimal + 1) + other.Digits
		if len(s1) < len(s2) {
			s1 = s1 + d.stringMultiply( '0', len(s2) - len(s1))
		}else if len(s2) < len(s1){
			s2 = s2 + d.stringMultiply('0', len(s1)-len(s2))
		}
		result := strings.Compare(s1, s2)
		if d.Negative {
			result = -result
		}
		return result
	}
}

func (d *Decimal)IsWholeNUmber()(bool, error){
	s, err := d.AsDecimal()
	if err != nil {
		return false, err
	}
	b := !strings.Contains(s, ".")
	return b, nil
}

func (d *Decimal)AsDecimal()(string, error){
	var err error
	result := d.Digits
	if d.Decimal != len(d.Digits){
		if d.Decimal < 0 {
			result = "0." + d.stringMultiply('0', 0 - d.Decimal) + d.Digits
		}else if d.Decimal < len(result){
			if d.Decimal == 0{
				result = "0." + result
			} else {
				result, err = d.insert(".", result, d.Decimal)
				if err != nil {
					return "", err
				}
			}
		}else{
			result = result + d.stringMultiply('0', d.Decimal - len(result) )
		}
	}
	if d.Negative && !d.allZeros(result, 0){
		result = "-" + result
	}
	return result, nil
}

func (d *Decimal)AsInteger()(int, error){
	b, err := d.IsWholeNUmber()
	if err!= nil {
		return 0,err
	}
	if !b {
		return 0, fmt.Errorf("Unable to represent "+d.String()+" as an integer")
	}
	dec := new(Decimal)
	dec.setValueDecimal(strconv.Itoa(-0x80000000))
	if d.ComparesTo(dec)<0 {
		return 0, fmt.Errorf("Unable to represent "+dec.String()+" as a signed 8 byte integer")
	}
	dec.setValueDecimal(strconv.Itoa(0x7fffffff))
	if d.ComparesTo(dec)>0 {
		return 0, fmt.Errorf("Unable to represent "+dec.String()+" as a signed 8 byte integer")
	}
	de, err := d.AsDecimal()
	if err!= nil {
		return 0,err
	}
	return strconv.Atoi(de)
}
