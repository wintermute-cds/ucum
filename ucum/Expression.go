package ucum

import (
	"bytes"
	"strconv"
)

type ExpressionComposer struct{
}

func Compose(item interface{}, canonicalValue bool)string{
	if item == nil {
		return "1"
	}
	var buffer *bytes.Buffer
	ec := &ExpressionComposer{}
	if _,instanceof := item.(*Term); instanceof {
		ec.composeTerm(buffer, item.(*Term))
	}else if _,instanceof := item.(*Canonical); instanceof {
		ec.composeCanonical(buffer, item.(*Canonical), canonicalValue)
	}else{
		panic("Can only compose expression from Term or Canonical")
	}
	return buffer.String()
}

func (e *ExpressionComposer)composeTerm(buffer *bytes.Buffer, term *Term){
	if term.Comp!=nil{
		e.composeComp(buffer, term.Comp.(Component))
	}
	if term.Op > 0 {
		e.composeOp(buffer, term.Op)
	}
	if term.Term!=nil{
		e.composeTerm(buffer, term)
	}
}
func (e *ExpressionComposer)composeComp(buffer *bytes.Buffer, comp Componenter){
	if _,instanceof := comp.(*Factor); instanceof {
		e.composeFactor(buffer, comp.(*Factor))
	}else if _,instanceof := comp.(*Symbol); instanceof {
		e.composeSymbol(buffer, comp.(*Symbol))
	}else if _,instanceof := comp.(*Term); instanceof {
		buffer.WriteString("(")
		e.composeTerm(buffer, comp.(*Term))
		buffer.WriteString(")")
	}else {
		buffer.WriteString("?")
	}
}
func (e *ExpressionComposer)composeSymbol(buffer *bytes.Buffer, symbol *Symbol){
	if symbol.Prefix!=nil{
		buffer.WriteString(symbol.Prefix.Code)
	}
	buffer.WriteString(symbol.Unit.GetCode())
	if symbol.Exponent!=1{
		buffer.WriteString(strconv.Itoa(symbol.Exponent))
	}
}
func (e *ExpressionComposer)composeFactor(buffer *bytes.Buffer, factor *Factor){
	buffer.WriteString(strconv.Itoa(factor.Value))
}
func (e *ExpressionComposer)composeOp(buffer *bytes.Buffer, op Operator){
	if op==DIVISION{
		buffer.WriteString("/")
	}else{
		buffer.WriteString(".")
	}
}
func (e *ExpressionComposer)composeCanonical(buffer *bytes.Buffer, can *Canonical, canonicalValue bool){
	if canonicalValue{
		buffer.WriteString(can.Value.AsDecimal())
	}
	first := true
	for _,c := range can.Units {
		if first {
			first = false
		}else{
			buffer.WriteString(".")
		}
		buffer.WriteString(c.Base().Code)
		if c.Exponent != 1 {
			buffer.WriteString(strconv.Itoa(c.Exponent))
		}
	}
}



type ExpressionParser struct {

}
