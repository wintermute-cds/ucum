package ucum


import (
	"encoding/xml"
	"io"
	"strings"
	"time"
	"github.com/bertverhees/ucum/decimal"
)

type DefinitionParser struct {
}

func (d *DefinitionParser) UnmarshalTerminology(reader io.Reader) (*UcumModel, error) {
	xmlUCUM := &XMLRoot{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(xmlUCUM); err != nil {
		return nil, err
	}
	return xmlUCUM.UcumModel()
}

type XMLRoot struct {
	Version      string           `xml:"version,attr"`
	Revision     string           `xml:"revision,attr"`
	RevisionDate string           `xml:"revision-date,attr"`
	Prefixes     []XMLPrefix      `xml:"prefix"`
	BaseUnits    []XMLBaseUnit    `xml:"base-unit"`
	DefinedUnits []XMLDefinedUnit `xml:"unit"`
}

func (x *XMLRoot) UcumModel() (*UcumModel, error) {
	var err error
	dateTime, err := x.ProcessRevisionDate(x.RevisionDate[7:32])
	if err != nil {
		return nil, err
	}
	ucumModel := &UcumModel{
		Version:      x.Version,
		Revision:     x.Revision,
		RevisionDate: dateTime,
		Prefixes:     make([]*Prefix, 0),
		BaseUnits:    make([]*BaseUnit, 0),
		DefinedUnits: make([]*DefinedUnit, 0),
		BaseUnitsByCode : make(map[string]*BaseUnit),
		DefinedUnitsByCode : make(map[string]*DefinedUnit),
	}
	for _, xmlItem := range x.Prefixes {
		names := make([]string, 0)
		name := xmlItem.Name
		names = append(names, name)
		value, err := decimal.NewFromString(xmlItem.Value.Value)
		if err!=nil {
			return nil, err
		}
		prefix := &Prefix{}
		prefix.Code = xmlItem.Code
		prefix.CodeUC = xmlItem.CodeUC
		prefix.Names = names
		prefix.PrintSymbol = xmlItem.PrintSymbol
		prefix.Value = value
		prefix.Kind = PREFIX
		ucumModel.Prefixes = append(ucumModel.Prefixes, prefix)
	}
	for _, xmlItem := range x.BaseUnits {
		names := make([]string, 0)
		name := xmlItem.Name
		names = append(names, name)
		if err != nil {
			return nil, err
		}
		baseUnit := &BaseUnit{}
		baseUnit.Code = xmlItem.Code
		baseUnit.CodeUC = xmlItem.CodeUC
		baseUnit.Names = names
		baseUnit.PrintSymbol = xmlItem.PrintSymbol
		baseUnit.Property = xmlItem.Property
		baseUnit.Dim = xmlItem.Dim
		baseUnit.Kind = BASEUNIT
		ucumModel.BaseUnits = append(ucumModel.BaseUnits, baseUnit)
		ucumModel.BaseUnitsByCode[baseUnit.Code] = baseUnit
	}
	for _, xmlItem := range x.DefinedUnits {
		names := make([]string, 0)
		name := xmlItem.Name
		names = append(names, name)
		value := &Value{}
		xmlItem2 := xmlItem.Value
		value.Unit = xmlItem2.Unit
		value.UnitUC = xmlItem2.UnitUC
		if strings.Trim(xmlItem2.Value, " ") != "" {
			value.Value, err = decimal.NewFromString(xmlItem2.Value)
			if err != nil {
				return nil, err
			}
		}
		unit := &DefinedUnit{}
		unit.Code = xmlItem.Code
		unit.CodeUC = xmlItem.CodeUC
		unit.Names = names
		unit.PrintSymbol = xmlItem.PrintSymbol
		property := strings.Trim(xmlItem.Property," ")
		unit.Property = property
		unit.Class = xmlItem.Class
		unit.IsSpecial = xmlItem.IsSpecial == "yes"
		unit.Metric = xmlItem.Metric == "yes"
		unit.Value = value
		unit.Kind = UNIT
		ucumModel.DefinedUnits = append(ucumModel.DefinedUnits, unit)
		ucumModel.DefinedUnitsByCode[unit.Code] = unit
	}
	return ucumModel, err
}

func (x *XMLRoot) ProcessRevisionDate(revisionDate string) (time.Time, error) {
	//add the T for correct datetime parsing
	if strings.Index(revisionDate, "T") == -1 {
		revisionDate = revisionDate[:10] + "T" + revisionDate[11:]
	}
	//remove space between timezone and datetime
	revisionDate = strings.Replace(revisionDate, " ", "", -1)
	//add colon between timezone-hours and timezone-minutes
	if len(revisionDate)>20 {
		revisionDate = revisionDate[:22] + ":" + revisionDate[22:]
	}
	time_, err := time.Parse(time.RFC3339, revisionDate)
	if err != nil {
		//suppress error
		//ucum-essence.xml has no known date-time notation.
		//I changed it in the current file, but wrong dates may slip through in the future.
		//This is how it should be:
		//revision-date="$Date: 2013-10-21T21:24:43-07:00 (Mon, 21 Oct 2013) $">
		//just normal ISO date-time notation, what else can would you expect from
		//Regenstrief Institute, Inc.
		//But it is: revision-date="$Date: 2015-11-13 15:13:19 -0500 (Fri, 13 Nov 2015) $
		return time.Now(), nil
	}
	return time_, nil
}

type XMLPrefix struct {
	XMLConcept
	Value XMLValue `xml:"value"` //"1e2" //precision 24
}

type XMLConcept struct {
	Code        string `xml:"Code,attr"`
	CodeUC      string `xml:"CODE,attr"`
	Name        string `xml:"name"`
	PrintSymbol string `xml:"printSymbol"`
}

type XMLDecimal struct {
	Decimal    int
	Digits     string
	Negative   bool
	Precision  int
	Scientific bool
}

type XMLUnit struct {
	XMLConcept
	Property string `xml:"property"`
}

type XMLBaseUnit struct {
	XMLUnit
	Dim rune `xml:"dim"`
}

type XMLDefinedUnit struct {
	XMLUnit
	Class     string   `xml:"class,attr"`
	IsSpecial string   `xml:"isSpecial,attr"`
	Metric    string   `xml:"isMetric,attr"`
	Value     XMLValue `xml:"value"`
}

type XMLValue struct {
	Unit   string `xml:"Unit,attr"`
	UnitUC string `xml:"UNIT,attr"`
	Value  string `xml:"value,attr"`
}
