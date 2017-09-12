/*******************************************************************************
 * Crown Copyright (c) 2006 - 2014, Copyright (c) 2006 - 2014 Kestral Computing & Health Intersections P/L.
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Kestral Computing P/L - initial implementation (pascal)
 *    Health Intersections P/L - port to Java
 *******************************************************************************/

/**
    Precision aware Decimal implementation. Any size number with any number of significant digits is supported.

    Note that operations are precision aware operations. Note that whole numbers are assumed to have
    unlimited precision. For example:
      2 x 2 = 4
      2.0 x 2.0 = 4.0
      2.00 x 2.0 = 4.0
    and
     10 / 3 = 3.33333333333333333333333333333333333333333333333
     10.0 / 3 = 3.33
     10.00 / 3 = 3.333
     10.00 / 3.0 = 3.3
     10 / 3.0 = 3.3

    Addition
      2 + 0.001 = 2.001
      2.0 + 0.001 = 2.0

    Note that the string representation is precision limited, but the internal representation
    is not.


  * This class is defined to work around the limitations of Java Big Decimal
 *
 * @author Grahame
 *
 */

 /**
 Decimal implementation for Golang for use in UCUM Service
  */

package ucum

import (
	"strings"
	"unicode"
	"fmt"
	"strconv"
)

const MAX_INT = 1<<31 - 1
const MIN_INT = -1 << 31


type Decimal struct {
	Decimal int
	Digits string
	Negative bool
	Precision int
	Scientific bool
}

func NewDecimal(value string)(*Decimal, error){
	var err error
	d := &Decimal{}
	value = strings.ToLower(value)
	if strings.Contains(value, "e"){
		err = d.setValueScientific(value)
	}else{
		err = d.setValueDecimal(value)
	}
	return d, err
}

func NewDecimalInt(value int)(*Decimal, error){
	var err error
	d := &Decimal{}
	err = d.setValueDecimal(strconv.Itoa(value))
	return d, err
}
/**
	 * There are a few circumstances where a simple value is known to be correct to a high
	 * precision. For instance, the unit prefix milli is not ~0.001, it is precisely 0.001
	 * to whatever precision you want to specify. This constructor allows you to specify
	 * an alternative precision than the one implied by the stated string
 */
func NewDecimalAndPrecision(value string, precision int)(*Decimal, error){
	var err error
	d := &Decimal{}
	value = strings.ToLower(value)
	if strings.Contains(value, "e"){
		err = d.setValueScientific(value)
	}else{
		err = d.setValueDecimal(value)
	}
	d.Precision = precision
	return d, err
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
		value = value[1:]
	}
	for i,c := range value{
		if c == '.' && dec == -1 {
			dec = i
		}else if c < '0' || c > '9' {
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

func (d *Decimal)countSignificants(value string)(int, error){
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
	return len(value), nil
}

func (d *Decimal)delete(value string, offset, length int)(string){
	if length + offset > len(value) {
		return value
	}
	if offset > len(value) {
		return value
	}
	if offset == 0 {
		return value[length:]
	}else{
		return value[0:offset]+value[offset+length:]
	}
}

func (d *Decimal)setValueScientific(value string)error {
	var err error
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
		max := MaxInt(d.Decimal, other.Decimal)
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

func (d *Decimal) IsWholeNumber()(bool){
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
	b := d.IsWholeNumber()
	if !b {
		return 0, fmt.Errorf("Unable to represent "+d.String()+" as an integer")
	}
	dec, err := NewDecimalInt(MIN_INT)
	if err!=nil {
		return 0, err
	}
	if d.ComparesTo(dec)<0 {
		return 0, fmt.Errorf("Unable to represent "+dec.String()+" as a signed 8 byte integer")
	}
	dec, err = NewDecimalInt(MAX_INT)
	if err!=nil {
		return 0, err
	}
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
	max := MaxInt(d.Decimal, other.Decimal)
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
		result.Precision = MinInt(d.Precision, other.Precision)
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
	max := MaxInt(d.Decimal, other.Decimal)
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
		result.Precision = MinInt(d.Precision, other.Precision)
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
		t := c + dig(rune(s1[i])) - dig(rune(s2[i]))
		if t < 0 {
			t = t + 10
			if i==0{
				panic("internal logic error")
			}else{
				s1 = d.replaceChar( s1, i-1, cdig(dig(rune(s1[i-1])) - 1))
			}
		}
		result[i] = cdig(t)
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

func (d *Decimal)Multiply(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	if d.IsZero()||other.IsZero() {
		return Zero()
	}
	maxi := MaxInt(d.Decimal, other.Decimal)
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
		t = MaxInt(t, len(sv))
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
	if d.IsWholeNumber() && other.IsWholeNumber() {
		prec = MaxInt(MaxInt(len(d.Digits), len(other.Digits)), MinInt(d.Precision, other.Precision))
	}else if d.IsWholeNumber(){
		prec = other.Precision
	}else if other.IsWholeNumber() {
		prec = d.Precision
	}else {
		prec = MinInt(d.Precision, other.Precision)
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

func (d *Decimal)Divide(other *Decimal)*Decimal {
	if other == nil {
		return nil
	}
	if d.IsZero() {
		return Zero()
	}
	if other.IsZero() {
		panic("Attempt to divide "+d.String()+" by zero")
	}

	s := "0" + other.Digits
	m := MaxInt(len(d.Digits),len(other.Digits)) + 40 //MaxInt loops we do
	tens := make([]string, 10)
	tens[0] = d.stringAddition(d.stringMultiply('0',len(s)), s)
	for i := 1; i<10; i++ {
		tens[i] = d.stringAddition(tens[i-1], s)
	}
	v := d.Digits
	r := ""
	l := 0
	di := (len(d.Digits) - d.Decimal + 1) - (len(other.Digits) - other.Decimal + 1)

	for {
		if len(v) >= len(tens[0]) {
			break
		}
		v = v + "0"
		di++
	}

	w := ""
	vi := 0
	if strings.Compare(v[0:len(other.Digits)], other.Digits) < 0 {
		if len(v) == len(tens[0]) {
			v = v + "0"
			di++
		}
		w = v[0:len(other.Digits)+1]
		vi = len(w)
	} else {
		w = "0" + v[0:len(other.Digits)]
		vi = len(w) + 1
	}

	handled := false
	proc := false

	for {
		if !(!(handled && ((l > m) || (( vi >= len(v)) && ((w == "" || d.allZeros(w, 0))))))) {
			break
		}
		l++
		handled = true
		proc = false
		for i := 8; i >= 0; i-- {
			if strings.Compare(tens[i], w) <= 0 {
				proc = true
				r = r + string(cdig(int32(i+1)))
				w = d.trimLeadingZeros(d.stringSubtraction(w, tens[i]))
				if !(handled && ((l > m) || (( vi >= len(v)) && ((w == "" || d.allZeros(w, 0)))))) {
					if vi < len(v) {
						w = w + string(v[vi])
						vi++
						handled = false
					} else {
						w = w + "0"
						di++
					}
					for {
						if !(len(w) < len(tens[0])) {
							break
						}
						w = "0" + w
					}
				}
				break
			}
		}
		if !proc {
			if w[0]=='0' {
				panic("w should not start with 0")
			}
			w = d.delete(w, 0, 1)
			r = r + "0"
			if !(handled && ((l > m) || (( vi >= len(v)) && ((w == "" || d.allZeros(w, 0)))))) {
				if vi < len(v) {
					w = w + string(v[vi])
					vi++
					handled = false
				} else {
					w = w + "0"
					di++
				}
				for {
					if !(len(w) < len(tens[0])) {
						break
					}
					w = "0" + w
				}
			}
		}
	}

	prec := 0

	if d.IsWholeNumber() && other.IsWholeNumber() && l < m {
		for i := 0; i < di; i++ {
			if strings.HasSuffix(r, "0"){
				r = d.delete(r, len(r) - 1, 1)
				di--
			}
		}
		prec = 100
	}else{
		if d.IsWholeNumber() && other.IsWholeNumber() {
			prec = MaxInt(len(d.Digits), len(other.Digits))
		}else if d.IsWholeNumber(){
			prec = MaxInt(other.Precision, len(r) - di)
		}else if other.IsWholeNumber(){
			prec = MaxInt(d.Precision, len(r) - di)
		}else{
			prec = MaxInt(MinInt(d.Precision, other.Precision), len(r) - di)
		}
		for{
			if !(len(r) > prec){
				break
			}
			up := strings.HasSuffix(r, "5")
			r = d.delete(r, len(r), 1)
			if up {
				i := len(r) - 1
				for{
					if !(up && i > 0){
						break
					}
					up = r[i] == '9'
					if up {
						r = r[0:i] + "0" + r[i+1:]
					} else {
						r = r[0:i] + string(cdig(dig(rune(r[i]))+1)) + r[i+1:]
					}
					i--
				}
				if up {
					r = "1" + r
					di++
				}else {
					r = r
				}
			}
			di--
		}
	}
	result := &Decimal{}
	result.setValueDecimal(r)
	result.Decimal = len(r) - di
	result.Negative = d.Negative != other.Negative
	result.Precision = prec
	result.Scientific = d.Scientific || other.Scientific
	return result
}

func (d *Decimal)trimLeadingZeros(s string)string{
	if s == ""{
		return ""
	}
	for{
		if strings.HasPrefix(s, "0"){
			s = s[1:]
		}else{
			break
		}
	}
	return s
}

func (d *Decimal)DivInt(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	t := d.Divide(other)
	return t.Trunc()
}

func (d *Decimal)Modulo(other *Decimal)*Decimal{
	if other == nil {
		return nil
	}
	t := d.DivInt(other)
	t2 := t.Multiply(other)
	return d.Subtract(t2)
}

func (d *Decimal)Equal(value, maxDifference *Decimal)bool{
	diff := d.Subtract(value).absolute()
	return diff.ComparesTo(maxDifference) <= 0
}

func (d *Decimal)absolute()*Decimal{
	di := d.copy()
	di.Negative = false
	return di
}



