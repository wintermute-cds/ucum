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
			d.Precision = d.countSignificants(value)
		}
		d.Digits = d.delete(value, d.Decimal, 1)
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

func (d *Decimal)countSignificants(value string)(int){
	i := strings.Index(value, ".")
	if i > -1 {
		value = d.delete(value, i, 1)
	}
	tmp := value
	for _,c := range tmp{
		if c != '0' {
			break
		}
		value = value[1:]
	}
	return len(value)
}

func (d *Decimal)delete(value string, offset, length int)(string){
	if length > len(value) {
		panic("Length cannot be greater then the length of the value string")
	}
	if offset > len(value) {
		panic("Offset cannot be greater then the length of the value string")
	}
	if offset == 0 {
		return value[length:]
	}else{
		return value[0:offset]+value[offset+length:]
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

func (d *Decimal)insert(ins, value string, offset int)(string){
	if offset > len(value) {
		panic("Offset cannot be greater then the length of the value string")
	}
	if offset == 0 {
		return ins + value
	}else{
		return value[0:offset] + ins + value[offset:]
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
		max := max(d.Decimal, other.Decimal)
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

func (d *Decimal)IsWholeNUmber()(bool){
	s := d.AsDecimal()
	b := !strings.Contains(s, ".")
	return b
}

func (d *Decimal)AsDecimal()(string){
	result := d.Digits
	if d.Decimal != len(d.Digits){
		if d.Decimal < 0 {
			result = "0." + d.stringMultiply('0', 0 - d.Decimal) + d.Digits
		}else if d.Decimal < len(result){
			if d.Decimal == 0{
				result = "0." + result
			} else {
				result = d.insert(".", result, d.Decimal)
			}
		}else{
			result = result + d.stringMultiply('0', d.Decimal - len(result) )
		}
	}
	if d.Negative && !d.allZeros(result, 0){
		result = "-" + result
	}
	return result
}

func (d *Decimal)AsInteger()(int, error){
	b := d.IsWholeNUmber()
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
	de := d.AsDecimal()
	return strconv.Atoi(de)
}

func (d *Decimal)AsScientific()(string){
	result := d.Digits
	zero := d.allZeros(result, 0)
	if zero {
		if d.Precision < 2 {
			result = "0e0"
		} else {
			result = "0." + d.stringMultiply('0', d.Precision-1) + "e0"
		}
	} else {
		if len(d.Digits) > 1 {
			result = d.insert( ".", result, 1)
		}
		result = result + "e" + strconv.Itoa( d.Decimal - 1)
	}
	if d.Negative && !zero {
		result = "-" + result
	}
	return result
}

func (d *Decimal)Trunc()*Decimal{
	if d.Decimal < 0 {
		return Zero()
	}

	result := d.copy()
	if len(result.Digits) >= result.Decimal {
		result.Digits = result.Digits[0:result.Decimal]
	}
	if result.Digits == "" {
		result.Digits = "0"
		result.Decimal = 1
		result.Negative = false
	}
	return result
}

func (d *Decimal)Add(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	result := d
	if d.Negative == other.Negative {
		result = d.doAdd(other)
		result.Negative = d.Negative
	}else if d.Negative {
		result = other.doSubtract(d)
	}else {
		result = d.doSubtract(other)
	}
	return result
}

func (d *Decimal)Subtract(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	result := d
	if d.Negative == !other.Negative {
		result = d.doAdd(other)
		result.Negative = true
	} else if !d.Negative && other.Negative {
		result = d.doAdd(other)
	} else if d.Negative && other.Negative {
		result = d.doSubtract(other)
		result.Negative = !result.Negative
	} else {
		result = other.doSubtract(d)
		result.Negative = !result.Negative
	}
	return result
}

func (d *Decimal)doAdd(other *Decimal)*Decimal{
	max := max(d.Decimal, other.Decimal)
	s1 := d.stringMultiply('0', max - d.Decimal + 1) + d.Digits
	s2 := d.stringMultiply('0', max - other.Decimal + 1) + other.Digits
	if len(s1) < len(s2) {
		s1 = s1 + d.stringMultiply( '0', len(s2) - len(s1))
	}else if len(s2) < len(s1){
		s2 = s2 + d.stringMultiply('0', len(s1)-len(s2))
	}

	s3 := d.stringAddition(s1, s2)

	if s3[0] == '1' {
		max++
	} else {
		s3 = d.delete(s3, 0, 1)
	}

	if max != len(s3) {
		if max < 0 {
			panic("Unhandled")
		}else if max < len(s3) {
			s3 = d.insert(".", s3, max)
		}else {
			panic("Unhandled")
		}
	}
	result := &Decimal{}
	result.setValueDecimal(s3)
	result.Scientific = d.Scientific || other.Scientific
	// todo: the problem with this is you have to figure out the absolute precision and take the lower of the two, not the relative one
	if d.Decimal < other.Decimal {
		result.Precision = d.Precision
	}else if other.Decimal < d.Decimal {
		result.Precision = other.Precision
	}else{
		result.Precision = min(d.Precision, other.Precision)
	}
	return result
}

func dig(c rune)int32{
	return c - '0'
}

func cdig(i int32)rune{
	return i + '0'
}

func (d *Decimal)doSubtract(other *Decimal)*Decimal{
	max := max(d.Decimal, other.Decimal)
	s1 := d.stringMultiply('0', max - d.Decimal + 1) + d.Digits
	s2 := d.stringMultiply('0', max - other.Decimal + 1) + other.Digits
	if len(s1) < len(s2) {
		s1 = s1 + d.stringMultiply( '0', len(s2) - len(s1))
	}else if len(s2) < len(s1){
		s2 = s2 + d.stringMultiply('0', len(s1)-len(s2))
	}

	s3 := ""

	neg := strings.Compare(s1, s2)<0
	if neg {
		s3 = s2
		s2 = s1
		s1 = s3
	}

	s3 = d.stringSubtraction(s1, s2)

	if s3[0] == '1' {
		max++
	} else {
		s3 = d.delete(s3, 0, 1)
	}

	if max != len(s3) {
		if max < 0 {
			panic("Unhandled")
		}else if max < len(s3) {
			s3 = d.insert(".", s3, max)
		}else {
			panic("Unhandled")
		}
	}
	result := &Decimal{}
	result.setValueDecimal(s3)
	result.Negative = neg
	result.Scientific = d.Scientific || other.Scientific
	// todo: the problem with this is you have to figure out the absolute precision and take the lower of the two, not the relative one
	if d.Decimal < other.Decimal {
		result.Precision = d.Precision
	}else if other.Decimal < d.Decimal {
		result.Precision = other.Precision
	}else{
		result.Precision = min(d.Precision, other.Precision)
	}
	return result
}

func (d *Decimal)stringAddition(s1, s2 string)string{
	if len(s1)!=len(s2){
		panic("string length assertion failed")
	}
	result := make([]int32,len(s2))
	for i := 0; i < len(s2);i++ {
		result[i]= '0'
	}
	c := int32(0)
	for i := len(s1) - 1; i>=0; i-- {
		t := c + dig(rune(s1[i])) + dig(rune(s2[i]))
		result[i] = cdig(t % 10)
		c = t / 10
	}
	if c!=0 {
		panic("c should be 0")
	}
	s := ""
	for i:=0; i<len(result);i++ {
		s = s + string(result[i])
	}
	return s
}

func (d *Decimal)stringSubtraction(s1, s2 string)string{
	if len(s1)!=len(s2){
		panic("string length assertion failed")
	}
	result := make([]int32,len(s2))
	for i := 0; i < len(s2);i++ {
		result[i]= '0'
	}
	c := int32(0)
	for i := len(s1) - 1; i>=0; i-- {
		t := c + dig(rune(s1[i])) + dig(rune(s2[i]))
		if t < 0 {
			t = t + 10
			if i==0{
				panic("internal logic error")
			}else{
				s1 = d.replaceChar( s1, i-1, cdig(dig(rune(s1[i-1])) - 1))
			}
			result[i] = cdig(t)
		}
	}
	if c!=0 {
		panic("c should be 0")
	}
	s := ""
	for i:=0; i<len(result);i++ {
		s = s + string(result[i])
	}
	return s
}

func (d *Decimal)replaceChar(s string, offset int, c rune)string{
	if offset == 0 {
		s = string(c) + s[1:]
	}else{
		s = s[0:offset] + string(c) + s[offset+1:]
	}
	return s
}

func (d *Decimal)multiply(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	if d.isZero()||other.isZero() {
		return Zero()
	}
	maxi := max(d.Decimal, other.Decimal)
	s1 := d.stringMultiply('0', maxi - d.Decimal + 1) + d.Digits
	s2 := d.stringMultiply('0', maxi - other.Decimal + 1) + other.Digits
	if len(s1) < len(s2) {
		s1 = s1 + d.stringMultiply( '0', len(s2) - len(s1))
	}else if len(s2) < len(s1){
		s2 = s2 + d.stringMultiply('0', len(s1)-len(s2))
	}
	s3 := ""
	if  strings.Compare(s2, s1)>0 {
		s3 = s2
		s2 = s1
		s1 = s3
	}
	s := make([]string,len(s2))
	tr := int32(0)
	for i := len(s2) - 1; i>=0;i-- {
		s[i] = d.stringMultiply('0', len(s2) - (i+1))
		c := int32(0)
		for j := len(s1) - 1; j >= 0; j-- {
			tr = c + dig(rune(s1[j])) + dig(rune(s2[j]))
			s[i] = d.insert(string(cdig( tr % 10)), s[i], 0)
			c = tr / 10
		}
		for {
			if c <= 0 {
				break
			}
			s[i] = d.insert(string(cdig( tr % 10)), s[i], 0)
			c = tr / 10
		}
	}

	t := 0
	for _,sv := range s {
		t = max(t, len(sv))
	}
	for i := 0; i < len(s); i++ {
		s[i] = d.stringMultiply('0', t-len(s1)) + s[i]
	}
	res := ""
	c := int32(0)
	for i := t -1 ; i >=0; i-- {
		for j := 0; j < len(s); j++ {
			c = c + dig(rune(s[j][i]))
		}
		res = d.insert(string(cdig( c % 10)), res, 0)
		c = c / 10
	}
	if c > 0 {
		panic("internal logic error")
	}

	dec := len(res) - ((len(s1)-(maxi+1))*2)

	for {
		if res != "" && res != "0" && res[0] == '0' {
			res = res[1:]
			dec--
		}
	}

	prec := 0
	if d.IsWholeNUmber() && other.IsWholeNUmber() {
		prec = max(max(len(d.Digits), len(other.Digits)), min(d.Precision, other.Precision))
	}else if d.IsWholeNUmber(){
		prec = other.Precision
	}else if other.IsWholeNUmber() {
		prec = d.Precision
	}else {
		prec = min (d.Precision, other.Precision)
	}
	res = d.delete(res, len(res)-1, 1)

	for {
		if !(len(res)>prec && res[len(res)]=='0') {
			break
		}
	}

	result := &Decimal{}
	result.Precision = prec
	result.Decimal = dec
	result.Negative = d.Negative != other.Negative
	result.Scientific = d.Scientific || other.Scientific
	return result
}


