package ucum

import (
	"UCUM_Golang/ucum/special"
	"reflect"
)

type Converter struct{
	Model *UcumModel
	Handlers *special.Registry
}

func NewConverter(model *UcumModel, handlers *special.Registry)*Converter{
	r := &Converter{}
	r.Model = model
	r.Handlers = handlers
	return r
}

func (c *Converter)Convert(term Term)*Canonical{
	return c.normaliseTerm(" ", term)
}

func (c *Converter)normaliseTerm(indent string, term Term)*Canonical{
	result,_ := NewCanonical(One())
	div := false
	t := &term
	for{
		if t == nil {
	  		break
		}
		if _,instanceof := t.Comp.(Term); instanceof {
			temp := c.normaliseTerm( indent + " ", t.Comp.(Term))
			if div {
				result.DivideValueDecimal(temp.Value)
				for _,c := range temp.Units {
					c.Exponent = 0 - c.Exponent
				}
			}else{
				result.MultiplyValueDecimal(temp.Value)
			}
			result.Units = append(result.Units, temp.Units...)
		}else if _,instanceof := t.Comp.(Factor); instanceof{
			if div {
				result.DivideValueInt(t.Comp.(Factor).Value)
			}else{
				result.MultiplyValueInt(t.Comp.(Factor).Value)
			}
		}else if _,instanceof := t.Comp.(Symbol); instanceof {
			o := t.Comp.(Symbol)
			temp := c.normaliseSymbol(indent, o)
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
	return result
}

func (c *Converter)normaliseSymbol(indent string, sym Symbol)*Canonical {
	result,_ :=  NewCanonical(One())
	if _,instanceof := sym.Unit.(BaseUnit); instanceof {
		cf,_ := NewCanonicalUnit(&sym.Unit.(BaseUnit), sym.Exponent)
		result.Units = append(result.Units, cf)
	}else{
		can := c.expandDefinedUnit(indent, sym.Unit.(DefinedUnit))
		for _,c := range can.Units {
			c.Exponent = c.Exponent * sym.Exponent
		}
		result.Units = append(result.Units, can.Units...)
		if sym.Exponent > 0 {
			for i := 0; i < sym.Exponent; i++ {
				result.MultiplyValueDecimal(can.Value)
			}
		}else{
			for i := 0; i > sym.Exponent; i-- {
				result.DivideValueDecimal(can.Value)
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
	}
	return result
}

func (c *Converter)expandDefinedUnit(indent string, unit DefinedUnit)*Canonical{
	u := unit.Value.Unit
	if unit.IsSpecial{
		if !c.Handlers.Exists(unit.Code){
			panic("Not handled yet (special unit)")
		}else{
			u = c.Handlers.Get(unit.Code).GetUnits()
		}
	}
	t := NewExpressionParser(model).Parse(u)
	result := c.normaliseTerm(indent + " ", t)
	result.MultiplyValueDecimal(unit.Value.Value)
	return result
}

