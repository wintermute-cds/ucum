package ucum

import (
	"bytes"
	"strconv"
	"fmt"
	"strings"
)

// COMPOSER==================================================================================================

type ExpressionComposer struct{
}

func ComposeExpression(item interface{}, canonicalValue bool)string{
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

// FORMALSTRUCTTURECOMPOSER==================================================================================================

type FormalStructureComposer struct{
}

func ComposeFormalStructure(term *Term)string{
	var buffer *bytes.Buffer
	ec := &FormalStructureComposer{}
	ec.composeTerm(buffer, term)
	return buffer.String()
}

func (e *FormalStructureComposer)composeTerm(buffer *bytes.Buffer, term *Term){
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
func (e *FormalStructureComposer)composeComp(buffer *bytes.Buffer, comp Componenter){
	if _,instanceof := comp.(*Factor); instanceof {
		e.composeFactor(buffer, comp.(*Factor))
	}else if _,instanceof := comp.(*Symbol); instanceof {
		e.composeSymbol(buffer, comp.(*Symbol))
	}else if _,instanceof := comp.(*Term); instanceof {
		e.composeTerm(buffer, comp.(*Term))
	}else {
		buffer.WriteString("?")
	}
}
func (e *FormalStructureComposer)composeSymbol(buffer *bytes.Buffer, symbol *Symbol){
	buffer.WriteString("(")
	if symbol.Prefix!=nil{
		buffer.WriteString(symbol.Prefix.Names[0])
	}
	buffer.WriteString(symbol.Unit.GetNames()[0])
	if symbol.Exponent!=1{
		buffer.WriteString("^")
		buffer.WriteString(strconv.Itoa(symbol.Exponent))
	}
	buffer.WriteString(")")
}
func (e *FormalStructureComposer)composeFactor(buffer *bytes.Buffer, factor *Factor){
	buffer.WriteString(strconv.Itoa(factor.Value))
}
func (e *FormalStructureComposer)composeOp(buffer *bytes.Buffer, op Operator){
	if op==DIVISION{
		buffer.WriteString(" / ")
	}else{
		buffer.WriteString(" * ")
	}
}

// PARSER==================================================================================================

type ExpressionParser struct {
	Model *UcumModel
}

func NewExpressionParser(model *UcumModel)*ExpressionParser{
	e := &ExpressionParser{}
	e.Model = model
	return e
}

func (p *ExpressionParser)Parse(code string)(*Term, error){
	l := NewLexer(code)
	l.Consume()
	res, err := p.parseTerm(l, true)
	if err!=nil{
		return nil, err
	}
	if !l.Finished(){
		return nil, fmt.Errorf("Expression was not parsed completely. Syntax Error?")
	}
	return res, nil
}

func (p *ExpressionParser)parseTerm(l *Lexer, first bool)(*Term, error){
	var err error
	res := &Term{}
	if first && l.TokenType == NONE{
		res.Comp = NewFactor(1)
	}else if l.TokenType == SOLIDUS {
		res.Op = DIVISION
		l.Consume()
		res.Term, err = p.parseTerm(l, false)
		if err!=nil{
			return nil, err
		}
	}else{
		if l.TokenType == ANNOTATION {
			res.Comp = NewFactor(1)
			l.Consume()
		}else{
			res.Comp, err = p.parseComp(l)
			if err!=nil{
				return nil, err
			}
		}
		if l.TokenType!=NONE && l.TokenType!=CLOSE {
			if l.TokenType==SOLIDUS{
				res.Op = DIVISION
				l.Consume()
			}else if l.TokenType == PERIOD{
				res.Op = MULTIPLICATION
				l.Consume()
			}else if l.TokenType== ANNOTATION {
				res.Op = MULTIPLICATION
			}else{
				return nil, fmt.Errorf("Error processing unit '"+l.Source+"': "+"Expected '/' or '.'"+"' at position "+strconv.Itoa(l.Start))
			}
			res.Term,err = p.parseTerm(l, false)
		}
	}
	return res, nil
}

func (p *ExpressionParser)parseComp(l *Lexer)(Componenter, error){
	if l.TokenType == NUMBER {
		fact := NewFactor(l.TokenAsInt())
		l.Consume()
		return fact, nil
	}else if l.TokenType == SYMBOL {
		return p.parseSymbol(l)
	}else if l.TokenType == NONE {
		return nil, fmt.Errorf("Error processing unit '"+l.Source+"': "+"unexpected end of expression looking for a symbol or a number"+"' at position "+strconv.Itoa(l.Start))
	}else if l.TokenType == OPEN {
		l.Consume()
		res, err := p.parseTerm(l, true)
		if err!=nil{
			return nil, err
		}
		if l.TokenType == CLOSE {
			l.Consume()
		} else{
			return nil, fmt.Errorf("Error processing unit '"+l.Source+"': "+"Unexpected Token Type '" + l.TokenType.String() + "' looking for a close bracket"+"' at position "+strconv.Itoa(l.Start))
		}
		return res, nil
	}else{
		return nil, fmt.Errorf("Error processing unit '"+l.Source+"': "+"unexpected token looking for a symbol or a number"+"' at position "+strconv.Itoa(l.Start))
	}
	return nil, nil
}

func (p *ExpressionParser)parseSymbol(l *Lexer)(Componenter, error){
	symbol := &Symbol{}
	sym := l.Token

	// now, can we pick a prefix that leaves behind a metric unit?
	var selected *Prefix
	var unit Uniter
	for _,prefix := range p.Model.Prefixes {
		if strings.HasPrefix(sym, prefix.Code){
			unit = p.Model.GetUnit(sym[len(prefix.Code):])
			if unit != nil && (unit.GetKind() == BASEUNIT || unit.(DefinedUnit).Metric){
				selected = prefix
				break
			}
		}
	}
	if selected != nil {
		symbol.Prefix = selected
		symbol.Unit = unit
	}else{
		unit = p.Model.GetUnit(sym)
		if unit!=nil{
			symbol.Unit = unit
		}else if sym!="1"{
			return nil, fmt.Errorf("Error processing unit '"+l.Source+"': "+"The unit '" + sym + "' is unknown"+"' at position "+strconv.Itoa(l.Start))
		}
	}

	l.Consume()
	if l.TokenType == NUMBER {
		symbol.Exponent = l.TokenAsInt()
		l.Consume()
	}else{
		symbol.Exponent = 1
	}
	return symbol, nil
}

// LEXER==================================================================================================

const NO_CHAR = 0

type Lexer struct{
	Source string
	Index int
	Token string
	TokenType TokenType
	Start int
}

func NewLexer(source string)*Lexer{
	l := &Lexer{}
	l.Source = source
	l.Index = 0
	return l
}

func (l *Lexer)Consume()error{
	l.Token = ""
	l.TokenType = NONE
	l.Start = l.Index
	if l.Index < len(l.Source){
		ch := l.nextChar()
		annotation, err := l.checkAnnotation(ch)
		checkNumber, err := l.checkNumber(ch)
		checkNumberOrSymbol, err := l.checkNumberOrSymbol(ch)
		if err != nil {
			return err
		}
		if !(
			l.checkSingleChar(ch, '/', SOLIDUS) || l.checkSingleChar(ch, '.', PERIOD)||
				l.checkSingleChar(ch, '(', OPEN)|| l.checkSingleChar(ch, ')', CLOSE)|| annotation||
					checkNumber || checkNumberOrSymbol ){
			return fmt.Errorf("Error processing unit '"+l.Source+"': unexpected character '"+string(ch)+"' at position "+strconv.Itoa(l.Start))
		}
	}
	return nil
}

func (l *Lexer)nextChar()rune{
	var r rune
	if l.Index < len(l.Source){
		r = []rune(l.Source)[l.Index]
	}else {
		r = NO_CHAR
	}
	l.Index++
	return r
}

func (l *Lexer)checkSingleChar(ch rune, test rune, tokenType TokenType)bool{
	if ch==test {
		l.Token = string(ch)
		l.TokenType = tokenType
		return true
	}
	return false
}

func (l *Lexer)checkAnnotation(ch rune)(bool, error){
	if ch == '{'{
		b := ""
		for{
			if ch=='}'{
				break
			}
			ch = l.nextChar()
			if !IsAsciiChar(ch){
				return false, fmt.Errorf("Error processing unit'"+l.Source+"': Annotation contains non-ascii characters")
			}
			if ch == NO_CHAR{
				return false, fmt.Errorf("Error processing unit'"+l.Source+"': unterminated annotation")
			}
			b = b + string(ch)
		}
		l.Token = b
		l.TokenType = ANNOTATION
		return true, nil
	}
	return false, nil
}

func (l *Lexer)checkNumber(ch rune)(bool, error){
	if ch=='+' || ch == '-' {
		l.Token = string(ch)
		ch = l.peekChar()
		for{
			if !( ch>='0' && ch <= '9'){
				break
			}
			l.Token = l.Token + string(ch)
			l.Index++
			ch = l.peekChar()
		}
		if len(l.Token)==1{
			return false, fmt.Errorf("Error processing unit'"+l.Source+"': unexpected character '"+string(ch)+"' at position "+strconv.Itoa(l.Start)+": a + or - must be followed by at least one digit")
		}
		l.TokenType = NUMBER
		return true, nil
	}
	return false, nil
}

func (l *Lexer)checkNumberOrSymbol(ch rune)(bool, error){
	var err error
	isSymbol := false
	isInBrackets := false
	if l.isValidSymbolChar(ch, true, false){
		l.Token = string(ch)
		isSymbol = !( ch>='0' && ch <= '9')
		isInBrackets, err = l.checkBrackets(ch, isInBrackets)
		if err!=nil{
			return false, err
		}
		ch = l.peekChar()
		isInBrackets, err = l.checkBrackets(ch, isInBrackets)
		if err!=nil{
			return false, err
		}
		for{
			if !(l.isValidSymbolChar(ch, !isSymbol || isInBrackets, isInBrackets)){
				break
			}
			l.Token = l.Token + string(ch)
			isSymbol = isSymbol || (ch!=NO_CHAR && !( ch>='0' && ch <= '9'))
			l.Index++
			ch = l.peekChar()
			isInBrackets, err = l.checkBrackets(ch, isInBrackets)
			if err!=nil{
				return false, err
			}
		}
		if isSymbol{
			l.TokenType = SYMBOL
		}else{
			l.TokenType = NUMBER
		}
		return true, nil
	}
	return false, nil
}

func (l *Lexer)checkBrackets(ch rune, isInBrackets bool)(bool, error) {
	if ch == '[' {
		if isInBrackets {
			return false, fmt.Errorf("Error processing unit '"+l.Source+"': "+"Nested ["+"' at position "+strconv.Itoa(l.Start))
		}
	} else {
		return true, nil
	}
	if ch == ']' {
		if !isInBrackets {
			return false, fmt.Errorf("Error processing unit '"+l.Source+"': "+"] without ["+"' at position "+strconv.Itoa(l.Start))
		}
	}else {
		return false, nil
	}
	return isInBrackets, nil
}

func (l *Lexer)isValidSymbolChar(ch rune, allowDigits, isInBrackets bool) bool{
	return 	allowDigits && ch >= '0' && ch <= '9' || ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '[' ||
		ch == ']' || ch == '%' || ch == '*' || ch == '^' || ch == '\'' || ch == '"' || ch == '_' || isInBrackets && ch == '.';
}

func (l *Lexer)peekChar()rune{
	var r rune
	if l.Index < len(l.Source){
		r = []rune(l.Source)[l.Index]
	}else {
		r = NO_CHAR
	}
	return r
}

func (l *Lexer)Finished()bool{
	return l.Index == len(l.Source)
}

func (l *Lexer)TokenAsInt()int{
	if l.Token[0]=='+' {
		return int([]rune(l.Token)[1])
	} else {
		return int([]rune(l.Token)[0])
	}
}