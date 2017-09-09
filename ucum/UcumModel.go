package ucum

import (
	"time"
	"regexp"
	"strings"
)

type UcumModel struct {
	Version string
	Revision string
	RevisionDate time.Time
	Prefixes []*Prefix
	BaseUnits []*BaseUnit
	DefinedUnits []*DefinedUnit
}

func NewUcumModel(version, revision string, revisionDate time.Time)*UcumModel{
	r := &UcumModel{}
	r.Version = version
	r.Revision = revision
	r.RevisionDate = revisionDate
	r.Prefixes = make([]*Prefix,0)
	r.BaseUnits = make([]*BaseUnit,0)
	r.DefinedUnits = make([]*DefinedUnit,0)
	return r
}

func (u *UcumModel)GetUnit(code string)Uniter{
	for _,unit := range u.BaseUnits {
		if unit.Code == code {
			return unit
		}
	}
	for _,unit := range u.DefinedUnits {
		if unit.Code == code {
			return unit
		}
	}
	return nil
}

func (u *UcumModel)Search(kind ConceptKind, text string, isRegex bool)[]Concepter{
	concepts := make([]Concepter,0)
	if kind == 0 || kind == PREFIX {
		concepts = append(concepts, u.searchPrefixes(text, isRegex)...)
	}
	if kind == 0 || kind == BASEUNIT || kind == UNIT {
		concepts = append(concepts, u.searchUnits(text, isRegex, kind)...)
	}
}

func (u *UcumModel)searchPrefixes(text string, isRegex bool)[]Concepter{
	concepts := make([]Concepter,0)
	return concepts
}

func (u *UcumModel)searchUnits(text string, isRegex bool, kind ConceptKind)[]Concepter{
	concepts := make([]Concepter,0)
	if kind == BASEUNIT {
		for _,unit := range u.BaseUnits{
			if u.matchesUnit(unit, text, isRegex){
				concepts = append(concepts,unit)
			}
		}
	}
	if kind == UNIT {
		for _,unit := range u.DefinedUnits{
			if u.matchesUnit(unit, text, isRegex){
				concepts = append(concepts,unit)
			}
		}
	}
	return concepts
}

func (u *UcumModel)matchesUnit(unit Uniter, text string, isRegex bool)bool{
	return u.matches(unit.GetProperty(), text, isRegex)||
}

func (u *UcumModel)matches(value, text string, isRegEx bool)bool{
	if isRegEx {
		b,_ := regexp.MatchString(value, text)
		return b
	}else{
		return strings.Contains(strings.ToLower(value), strings.ToLower(text))
	}
}


