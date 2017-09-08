package ucum

import "time"

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
