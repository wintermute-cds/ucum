package ucum

import (
	"strings"
	"sort"
	"strconv"
	"io"
	"encoding/xml"
	"time"
)

type DefinitionParser struct {
}

func (d *DefinitionParser)Parse(){

}

func unmarshalTerminology(reader io.Reader)(*UcumModel, error){
	xmlUCUM := &XMLRoot{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(xmlUCUM); err!=nil {
		return nil, err
	}
	return xmlUCUM.UcumModel()
}


type XMLRoot struct{
	Version string	`xml:"version"`
	Revision string		`xml:"revision"`
	RevisionDate string	`xml:"revision-date"`
	Prefixes []XMLPrefix	`xml:"prefix"`
	BaseUnits []XMLBaseUnit	`xml:"base-unit"`
	DefinedUnits []XMLDefinedUnit	`xml:"unit"`
}

func (x *XMLRoot)UcumModel()(*UcumModel, error){
	var err error
	dateTime, err := x.ProcessRevisionDate(x.RevisionDate[7:32])
	if err != nil {
		return nil, err
	}
	ucumModel := &UcumModel{
		Version: x.Version,
		Revision: x.Revision,
		RevisionDate: dateTime,
		Prefixes : make([]*Prefix,0),
		BaseUnits : make([]*BaseUnit,0),
		DefinedUnits : make([]*DefinedUnit,0),
	}
	for _, xmlPrefix := range x.Prefixes {
		names := make([]string,1)
		name := xmlPrefix.Name
		names = append(names, name)
		value, err := NewDecimalAndPrecision(xmlPrefix.Value, 24)
		if err != nil {
			return nil, err
		}
		prefix := &Prefix{}
		prefix.Code = xmlPrefix.Code
		prefix.CodeUC = xmlPrefix.CodeUC
		prefix.Names = names
		prefix.PrintSymbol = xmlPrefix.PrintSymbol
		prefix.Value = value
		ucumModel.Prefixes = append(ucumModel.Prefixes, prefix)
	}
	for _, xmlbaseUnit := range x.BaseUnits {
		baseUnit := &BaseUnit{

		}
		ucumModel.BaseUnits = append(ucumModel.BaseUnits, baseUnit)
	}
	for _, xmlUnit := range x.DefinedUnits {
		unit := &DefinedUnit{

		}
		ucumModel.DefinedUnits = append(ucumModel.DefinedUnits, unit)
	}
	return ucumModel,err
}

func (x *XMLRoot)ProcessRevisionDate(revisionDate string)(time.Time, error){
	time,err := time.Parse("2013-10-21 21:24:43 -0700", revisionDate)
	if err!=nil {
		return time.Time{}, err
	}
	return time, nil
}

type XMLPrefix struct{
	XMLConcept
	//private Prefix parsePrefix
	//prefix.setValue(new Decimal(xpp.getAttributeValue(null, "value"), 24));
	//<prefix xmlns="" Code="h" CODE="H">
	//<name>hecto</name>
	//<printSymbol>h</printSymbol>
	//<value value="1e2">1 &#215; 10<sup>2</sup>
	//</value>
	//</prefix>
	Value string	`xml:"value"` //"1e2" //precision 24
}

type XMLConcept struct{
	Code string		`xml:"Code,attr"`
	CodeUC string	`xml:"CODE,attr"`
	Name string	`xml:"name"`
	PrintSymbol string	`xml:"printSymbol"`
}

type XMLDecimal struct{
	Decimal int
	Digits string
	Negative bool
	Precision int
	Scientific bool
}

type XMLUnit struct{
	XMLConcept
	Property string
}

type XMLBaseUnit struct{
	XMLUnit
	Dim rune
}

type XMLDefinedUnit struct{
	XMLUnit
	Class string
	IsSpecial bool
	Metric bool
	Value Value
}


type XMLValue struct{
	Text string
	Unit string
	UnitUC string
	Value *Decimal
}

type XMLCanonical struct {
	Units []*XMLCanonicalUnit
	Value *Decimal
}

type XMLCanonicalUnit struct {
	base *XMLBaseUnit
	Exponent int
}

type XMLComponent struct{

}

//Factor
type XMLFactor struct {
	XMLComponent
	Value int
}

type XMLSymbol struct {
	XMLComponent
	Unit Uniter
	Prefix *Prefix
	Exponent int
}

type XMLTerm struct {
	Component
	Comp Componenter
	Op Operator
	Term *Term
}

