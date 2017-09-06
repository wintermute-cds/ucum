package ucum

import (
	"strings"
	"strconv"
)

// Concept
type Concepter interface {
	GetDescription() string
	String()
}

type Concept struct{
	Code string
	CodeUC string
	Kind ConceptKind
	Names []string
	PrintSymbol string
}

func NewConcept(kind ConceptKind, code string, codeUC string)(*Concept,error){
	c := &Concept{
		Kind : kind,
		Code: code,
		CodeUC: codeUC,
	}
	return c, nil
}

func (c *Concept)GetDescription()string {
	description := strings.ToLower(c.Kind.String()) + " " + c.Code + " ('" + c.Names[0] + "')"
	return description
}

func (c *Concept)String()string{
	return c.Code + " = " + c.GetDescription()
}

//Unit
type Uniter interface{
	Concepter
}

type Unit struct{
	Concept
	Property string
}

func NewUnit(kind ConceptKind, code string, codeUC string)(*Unit,error){
	u := &Unit{}
	u.Kind = kind
	u.Code = code
	u.CodeUC = codeUC
	return u, nil
}

func (u *Unit)GetDescription()string {
	return strings.ToLower(u.Kind.String()) + " " + u.Code + " ('" + u.Names[0] + "')" + " (" + u.Property + ")"
}

func (u *Unit)String()string{
	return u.Code + " = " + u.GetDescription()
}

//BaseUnit
type BaseUnit struct{
	Unit
	Dim rune
}

func NewBaseUnit(kind ConceptKind, code string, codeUC string)(*BaseUnit,error){
	b := &BaseUnit{}
	b.Kind = BASEUNIT
	b.Code = code
	b.CodeUC = codeUC
	return b, nil
}

//DefinedUnit
type DefinedUnit struct{
	Unit
	Class string
	IsSpecial bool
	Metric bool
	Value Value
}

func NewDefinedUnit(kind ConceptKind, code string, codeUC string)(*DefinedUnit,error){
	b := &DefinedUnit{}
	b.Kind = UNIT
	b.Code = code
	b.CodeUC = codeUC
	return b, nil
}

func (d *DefinedUnit)GetDescription()string {
	return strings.ToLower(d.Kind.String()) + " " + d.Code + " ('" + d.Names[0] + "')" + " (" + d.Property + ")" + " = " + d.Value.GetDescription()
}

//Prefix
type Prefix struct{
	Concept
	Value float64
}

func NewPrefix(kind ConceptKind, code string, codeUC string)(*Prefix,error){
	b := &Prefix{}
	b.Kind = PREFIX
	b.Code = code
	b.CodeUC = codeUC
	return b, nil
}

func (p *Prefix)GetDescription()string {
	return strings.ToLower(p.Kind.String()) + " " + p.Code + " ('" + p.Names[0] + "')" + " = " + strconv.FormatFloat(p.Value, 'E', -1, 64)
}

//Value
type Value struct{
	Text string
	Unit string
	UnitUC string
	Value float64
}

func NewValue(unit, unitUC string, value float64)(*Value, error){
	v := &Value{}
	v.Unit = unit
	v.UnitUC = unitUC
	v.Value = value
	return v, nil
}

func (v *Value)GetDescription()string {
	if v.Value == 0.0 {
		return v.Unit
	}
	return strconv.FormatFloat(v.Value, 'E', -1, 64) + v.Unit
}

//Canonical
type Canonical struct {
	Units []*CanonicalUnit
	Value float64
}

func NewCanonical(value float64)(*Canonical, error){
	v := &Canonical{
		Value:value,
		Units : make([]*CanonicalUnit, 0),
	}
	return v, nil
}

func (c *Canonical)MultiplyValueDecimal(multiplicand float64){
	c.Value = c.Value * multiplicand
}

func (c *Canonical)MultiplyValueInt(multiplicand int){
	c.Value = c.Value * float64(multiplicand)
}

func (c *Canonical)DivideValueDecimal(divisor float64){
	c.Value = c.Value / divisor
}

func (c *Canonical)DivideValueInt(divisor int){
	c.Value = c.Value / float64(divisor)
}


//CanonicalUnit
type CanonicalUnit struct {
	base *BaseUnit
	Exponent int
}

func NewCanonicalUnit(base *BaseUnit, exponent int)(*CanonicalUnit, error){
	v := &CanonicalUnit{}
	if base != nil {v.base = base}
	v.Exponent = exponent
	return v, nil
}

func (c *CanonicalUnit)Base()(*BaseUnit){
	return c.base
}

//Component
type Componenter interface{

}

type Component struct{

}

//Factor
type Factor struct {
	Component
	Value int
}

func NewFactor(value int)(*Factor, error){
	v := &Factor{
		Value:value,
	}
	return v,nil
}

//Symbol
type Symbol struct {
	Component
	Unit Uniter
	Prefix *Prefix
	Exponent int
}

func NewSymbol(unit Uniter, prefix *Prefix, exponent int)(*Symbol, error){
	v := &Symbol{}
	if unit != nil { v.Unit = unit }
	if prefix != nil { v.Prefix = prefix }
	v.Exponent = exponent
	return v,nil
}

func (s *Symbol)HasPrefix()(bool){
	return s.Prefix != nil
}

func (s *Symbol)InvertExponent(){
	s.Exponent = -s.Exponent
}

//Term
type Term struct {
	Component
	Comp Componenter
	Op *Operator
	Term *Term
}

func NewTerm()(*Term, error){
	return &Term{}, nil
}

func (t *Term)SetTermCheckOp(term *Term){
	if term!=nil {
		t.Term = term
		t.Op = term.Op
	}else{
		t.Term = nil
		t.Op = nil
	}
}

