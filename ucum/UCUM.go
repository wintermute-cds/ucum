package ucum

import (
	"strings"
	"strconv"
	"sort"
)

// Concept
type Concepter interface {
	GetDescription() string
	String() string
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

func (c Concept)GetDescription()string {
	description := strings.ToLower(c.Kind.String()) + " " + c.Code + " ('" + c.Names[0] + "')"
	return description
}

func (c Concept)String()string{
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

func (u Unit)GetDescription()string {
	return strings.ToLower(u.Kind.String()) + " " + u.Code + " ('" + u.Names[0] + "')" + " (" + u.Property + ")"
}

func (u Unit)String()string{
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

func (d DefinedUnit)GetDescription()string {
	return strings.ToLower(d.Kind.String()) + " " + d.Code + " ('" + d.Names[0] + "')" + " (" + d.Property + ")" + " = " + d.Value.GetDescription()
}

//Prefix
type Prefix struct{
	Concept
	Value *Decimal
}

func NewPrefix(kind ConceptKind, code string, codeUC string)(*Prefix,error){
	b := &Prefix{}
	b.Kind = PREFIX
	b.Code = code
	b.CodeUC = codeUC
	return b, nil
}

func (p Prefix)GetDescription()string {
	return strings.ToLower(p.Kind.String()) + " " + p.Code + " ('" + p.Names[0] + "')" + " = " + p.Value.String()
}

//Value
type Value struct{
	Text string
	Unit string
	UnitUC string
	Value *Decimal
}

func NewValue(unit, unitUC string, value *Decimal)(*Value, error){
	v := &Value{}
	v.Unit = unit
	v.UnitUC = unitUC
	v.Value = value
	return v, nil
}

func (v Value)GetDescription()string {
	if v.Value == nil {
		return v.Unit
	}
	return v.Value.String()
}

//Canonical
type Canonical struct {
	Units []*CanonicalUnit
	Value *Decimal
}

func (c *Canonical)RemoveFromUnits(i int){
	c.Units[i] = c.Units[len(c.Units)-1]
	c.Units[len(c.Units)-1] = nil
	c.Units = c.Units[:len(c.Units)-1]
}

type ByCode []*CanonicalUnit

func (a ByCode) Len() int           { return len(a) }
func (a ByCode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCode) Less(i, j int) bool { return a[i].base.Code < a[j].base.Code }

func (c *Canonical)SortUnits(){
	sort.Sort(ByCode(c.Units))
}

func NewCanonical(value *Decimal)(*Canonical, error){
	v := &Canonical{
		Value:value,
		Units : make([]*CanonicalUnit, 0),
	}
	return v, nil
}

func (c *Canonical)MultiplyValueDecimal(multiplicand *Decimal){
	c.Value = c.Value.Multiply(multiplicand)
}

func (c *Canonical)MultiplyValueInt(multiplicand int){
	c.Value = c.Value.Multiply(NewDecimal(strconv.Itoa(multiplicand)))
}

func (c *Canonical)DivideValueDecimal(divisor *Decimal){
	c.Value = c.Value.Divide(divisor)
}

func (c *Canonical)DivideValueInt(divisor int){
	c.Value = c.Value.Divide(NewDecimal(strconv.Itoa(divisor)))
}


//CanonicalUnit
type CanonicalUnit struct {
	base *BaseUnit
	Exponent int
}

func NewCanonicalUnit(base *BaseUnit, exponent int)(*CanonicalUnit, error){
	v := &CanonicalUnit{}
	v.base = base
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
	v.Unit = unit
	v.Prefix = prefix
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
	Op Operator
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
		t.Op = 0
	}
}

