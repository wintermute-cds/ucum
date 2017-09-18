package ucum

import (
	"reflect"
	"fmt"
)

type Converter struct{
	Model *UcumModel
	Handlers *Registry
}

func NewConverter(model *UcumModel, handlers *Registry)*Converter{
	r := &Converter{}
	r.Model = model
	r.Handlers = handlers
	return r
}

func (c *Converter)Convert(term *Term)(*Canonical, error){
	return c.normaliseTerm(" ", term)
}

func (c *Converter)normaliseTerm(indent string, term *Term)(*Canonical, error){
	d,_ := NewDecimal("1")
	result,_ := NewCanonical(d)
	div := false
	t := term
	for{
		if t == nil {
			break
		}
		if _,instanceof := t.Comp.(*Term); instanceof {
			temp, err := c.normaliseTerm( indent + " ", t.Comp.(*Term))
			if err!=nil{
				return nil, err
			}
			if div {
				result.DivideValueDecimal(temp.Value)
				for _,c := range temp.Units {
					c.Exponent = 0 - c.Exponent
				}
			}else{
				result.MultiplyValueDecimal(temp.Value)
			}
			result.Units = append(result.Units, temp.Units...)
		}else if _,instanceof := t.Comp.(*Factor); instanceof{
			if div {
				result.DivideValueInt(t.Comp.(*Factor).Value)
			}else{
				result.MultiplyValueInt(t.Comp.(*Factor).Value)
			}
		}else if _,instanceof := t.Comp.(*Symbol); instanceof {
			o := t.Comp.(*Symbol)
			temp, err := c.normaliseSymbol(indent, o)
			if err!=nil{
				return nil, err
			}
			if div {
				result.DivideValueDecimal(temp.Value)
				for _,c := range temp.Units {
					c.Exponent = 0 - c.Exponent
				}
			}else{
				result.MultiplyValueDecimal(temp.Value)
			}
			result.Units = append(result.Units, temp.Units...)
		}
		div = t.Op == DIVISION
		t = t.Term
	}
	for i := len(result.Units) - 1; i>=0; i-- {
		sf := result.Units[i]
		for j := i - 1; j>=0; j-- {
			st := result.Units[j]
			if reflect.DeepEqual(st.Base(), sf.Base()){
				st.Exponent = sf.Exponent + st.Exponent
				result.RemoveFromUnits( i)
				break
			}
		}
	}
	for i := len(result.Units) - 1; i>=0; i-- {
		if result.Units[i].Exponent == 0 {
			result.RemoveFromUnits( i)
		}

	}
	result.SortUnits()
	return result, nil
}

func (c *Converter)normaliseSymbol(indent string, sym *Symbol)(*Canonical, error) {
	d,_ := NewDecimal("1")
	result,_ :=  NewCanonical(d)
	bu,instanceof := sym.Unit.(*BaseUnit);
	if  instanceof {
		cf,_ := NewCanonicalUnit(bu, sym.Exponent)
		result.Units = append(result.Units, cf)
	}else {
		can, err := c.expandDefinedUnit(indent, sym.Unit.(*DefinedUnit))
		if err != nil {
			return nil, err
		}
		for _, c := range can.Units {
			c.Exponent = c.Exponent * sym.Exponent
		}
		result.Units = append(result.Units, can.Units...)
		if sym.Exponent > 0 {
			for i := 0; i < sym.Exponent; i++ {
				result.MultiplyValueDecimal(can.Value)
			}
		} else {
			for i := 0; i > sym.Exponent; i-- {
				result.DivideValueDecimal(can.Value)
			}
		}
	}
	if sym.Prefix != nil {
		if sym.Exponent > 0 {
			for i := 0; i < sym.Exponent; i++ {
				result.MultiplyValueDecimal(sym.Prefix.Value)
			}
		}else{
			for i := 0; i > sym.Exponent; i-- {
				result.DivideValueDecimal(sym.Prefix.Value)
			}
		}
	}
	return result, nil
}

func (c *Converter)expandDefinedUnit(indent string, unit *DefinedUnit)(*Canonical,error){
	u := unit.Value.Unit
	if unit.IsSpecial{
		if !c.Handlers.Exists(unit.Code){
			return nil, fmt.Errorf("Not handled yet (special unit)")
		}else{
			u = c.Handlers.Get(unit.Code).GetUnits()
		}
	}
	t, err := NewExpressionParser(c.Model).Parse(u)
	if err!=nil{
		return nil, err
	}
	result, err := c.normaliseTerm(indent + " ", t)
	if err!=nil{
		return nil, err
	}
	result.MultiplyValueDecimal(unit.Value.Value)
	return result,nil
}

