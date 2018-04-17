package ucum

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type UcumService interface {
	/**
	 * return Ucum Identification details for the version in use
	 */
	UcumIdentification() *UcumVersionDetails
	/**
	 * Check UCUM. Note that this stands as a test of the service
	 * more than UCUM itself (for version 1.7, there are no known
	 * semantic errors in UCUM). But you should always run this test at least
	 * once with the version of UCUM you are using to ensure that
	 * the service implementation correctly understands the UCUM data
	 * to which it is bound
	 *
	 * @return a list of internal errors in the UCUM spec.
	 *
	 */
	ValidateUCUM() []string
	/**
	 * Search through the UCUM concepts for any concept containing matching text.
	 * Search will be limited to the kind of concept defined by kind, or all if kind
	 * is null
	 *
	 * @param kind - can be null. scope of search
	 * @param text - required
	 * @param isRegex
	 * @return
	 */
	Search(kind ConceptKind, text string, isRegex bool) ([]Concepter, error)
	/**
	 * return a list of the defined types of units in this UCUM version
	 *
	 * @return
	 */
	GetProperties() []string
	/**
	 * validate whether a unit code are valid UCUM units
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	Validate(unit string) (bool, string)
	/**
	 * given a unit, return a formal description of what the units stand for using
	 * full names
	 * @param units the unit code
	 * @return formal description
	 * @throws UcumException
	 * @throws OHFException
	 */
	Analyse(unit string) (string, error)
	/**
	 * validate whether a units are valid UCUM units and additionally require that the
	 * units from a particular property
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	ValidateInProperty(unit, property string) string
	/**
	 * validate whether a units are valid UCUM units and additionally require that the
	 * units match a particular base canonical unit
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	ValidateCanonicalUnits(unit, canonical string) string
	/**
	 * given a set of units, return their canonical form
	 * @param unit
	 * @return the canonical form
	 * @throws UcumException
	 * @throws OHFException
	 */
	GetCanonicalUnits(unit string) (string, error)
	/**
	 * given two pairs of units, return true if they share the same canonical base
	 *
	 * @param units1
	 * @param units2
	 * @return
	 * @throws UcumException
	 * @
	 */
	IsComparable(units1, units2 string) (bool, error)
	/**
	 * for a given canonical unit, return all the defined units that have the
	 * same canonical unit.
	 *
	 * @param code
	 * @return
	 * @throws UcumException
	 * @throws OHFException
	 */
	GetDefinedForms(code string) ([]*DefinedUnit, error)
	/**
	 * given a value/unit pair, return the canonical form as a value/unit pair
	 *
	 * 1 mm -> 1e-3 m
	 * @param value
	 * @return
	 * @throws UcumException
	 * @throws OHFException
	 */
	GetCanonicalForm(value *Pair) (*Pair, error)
	/**
	 * given a value and source unit, return the value in the given dest unit
	 * an exception is thrown if the conversion is not possible
	 *
	 * @param value
	 * @param sourceUnit
	 * @param destUnit
	 * @return the value if a conversion is possible
	 * @throws UcumException
	 * @throws OHFException
	 */
	Convert(value *Decimal, sourceUnit, destUnit string) (*Decimal, error)
	/**
	 * multiply two value/units pairs together and return the result in canonical units
	 *
	 * Note: since the units returned are canonical,
	 * @param o1
	 * @param o2
	 * @return
	 * @throws UcumException
	 * @
	 */
	Multiply(o1, o2 *Pair) (*Pair, error)
	/**
	 * given a set of UCUM units, return a likely preferred human dense form
	 *
	 * SI units - as is.
	 * Other units - improved by manual fixes, or the removal of []
	 *
	 * @param code
	 * @return the preferred human display form
	 */
	GetCommonDisplay(code string) string
}

// UcumVersionDetails======================================================
type UcumVersionDetails struct {
	ReleaseDate time.Time
	Version     string
}

func NewUcumVersionDetails(releaseDate time.Time, version string) *UcumVersionDetails {
	r := &UcumVersionDetails{}
	r.ReleaseDate = releaseDate
	r.Version = version
	return r
}

//UcumEssenceService=======================================================
const UCUM_OID = "2.16.840.1.113883.6.8"

type UcumEssenceService struct {
	Model    *UcumModel
	Handlers *Registry
}

var instanceOfUcumEssenceService *UcumEssenceService

func GetInstanceOfUcumEssenceService(xmlFileName string) (*UcumEssenceService, error) {
	if instanceOfUcumEssenceService == nil {
		instanceOfUcumEssenceService = new(UcumEssenceService)
		xmlFile, err := os.Open(xmlFileName)
		if err != nil {
			return nil, err
		}
		defer xmlFile.Close()
		d := new(DefinitionParser)
		instanceOfUcumEssenceService.Model, err = d.UnmarshalTerminology(xmlFile)
		if err != nil {
			return nil, err
		}
	}
	return instanceOfUcumEssenceService, nil
}

func (u *UcumEssenceService) UcumIdentification() *UcumVersionDetails {
	d := &UcumVersionDetails{}
	d.ReleaseDate = u.Model.RevisionDate
	d.Version = u.Model.Version
	return d
}

func (u *UcumEssenceService) ValidateUCUM() []string {
	return NewUcumValidator(u.Model, u.Handlers).Validate()
}

func (u *UcumEssenceService) Search(kind ConceptKind, text string, isRegex bool) ([]Concepter, error) {
	if text == "" {
		return nil, fmt.Errorf("search", "text", "must not be empty")
	}
	return u.Model.Search(kind, text, isRegex), nil
}

func (u *UcumEssenceService) GetProperties() []string {
	result := make([]string, 0)
	for _, unit := range u.Model.DefinedUnits {
		result = append(result, unit.GetProperty())
	}
	return result
}

func (u *UcumEssenceService) Validate(unit string) (bool, string) {
	if unit == "" {
		return true, "search text must not be empty"
	}
	_, err := NewExpressionParser(u.Model).Parse(unit)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func (u *UcumEssenceService) Analyse(unit string) (string, error) {
	if unit == "" {
		return "(unity)", nil
	}
	term, err := NewExpressionParser(u.Model).Parse(unit)
	if err != nil {
		return "", err
	}
	return ComposeFormalStructure(term), nil
}

func (u *UcumEssenceService) ValidateInProperty(unit, property string) string {
	if unit == "" {
		return "validateInProperty: unit must not be null or empty"
	}
	if property == "" {
		return "validateInProperty: property must not be null or empty"
	}
	term, err := NewExpressionParser(u.Model).Parse(unit)
	if err != nil {
		return err.Error()
	}
	can, err := NewConverter(u.Model, u.Handlers).Convert(term)
	if err != nil {
		return err.Error()
	}
	cu := ComposeExpression(can, false)
	if len(can.Units) == 1 {
		if property == can.Units[0].Base.Property {
			return ""
		} else {
			return "unit " + unit + " is of the property type " + can.Units[0].Base.Property + " (" + cu + "), not " + property + " as required."
		}
	}
	if property == "concentration" && (cu == "g/L" || cu == "mol/L") {
		return ""
	}
	return "unit " + unit + " has the base units " + cu + ", and are not from the property " + property + " as required."
}

func (u *UcumEssenceService) ValidateCanonicalUnits(unit, canonical string) string {
	if unit == "" {
		return "ValidateCanonicalUnits: unit must not be null or empty"
	}
	if canonical == "" {
		return "ValidateCanonicalUnits: canonical must not be null or empty"
	}
	term, err := NewExpressionParser(u.Model).Parse(unit)
	if err != nil {
		return err.Error()
	}
	can, err := NewConverter(u.Model, u.Handlers).Convert(term)
	if err != nil {
		return err.Error()
	}
	cu := ComposeExpression(can, false)
	if canonical != cu {
		return "unit " + unit + " has the base units " + cu + ", not " + canonical + " as required."
	}
	return ""
}

func (u *UcumEssenceService) GetCanonicalUnits(unit string) (string, error) {
	if unit == "" {
		return "", fmt.Errorf("GetCanonicalUnits: unit must not be null or empty")
	}
	term, err := NewExpressionParser(u.Model).Parse(unit)
	if err != nil {
		return "", err
	}
	converter := NewConverter(u.Model, u.Handlers)
	can, err := converter.Convert(term)
	if err != nil {
		return "", err
	}
	return ComposeExpression(can, false), nil
}

func (u *UcumEssenceService) IsComparable(units1, units2 string) (bool, error) {
	if units1 == "" {
		return false, nil
	}
	if units2 == "" {
		return false, nil
	}
	u1, err := u.GetCanonicalUnits(units1)
	if err != nil {
		return false, err
	}
	u2, err := u.GetCanonicalUnits(units2)
	if err != nil {
		return false, err
	}
	return u1 == u2, nil
}

func (u *UcumEssenceService) GetDefinedForms(code string) ([]*DefinedUnit, error) {
	if code == "" {
		return nil, fmt.Errorf("getDefinedForms: code must not be null or empty")
	}
	result := make([]*DefinedUnit, 0)
	base := u.Model.getBaseUnit(code)
	if base != nil {
		for _, du := range u.Model.DefinedUnits {
			if !du.IsSpecial {
				s, err := u.GetCanonicalUnits(du.Code)
				if err != nil {
					return nil, err
				}
				if code == s {
					result = append(result, du)
				}
			}
		}
	}
	return result, nil
}

func (u *UcumEssenceService) GetCanonicalForm(value *Pair) (*Pair, error) {
	if value == nil {
		return nil, fmt.Errorf("getCanonicalForm: value must not be null")
	}
	if value.Code == "" {
		return nil, fmt.Errorf("getCanonicalForm: value.code must not be empty")
	}
	term, err := NewExpressionParser(u.Model).Parse(value.Code)
	if err != nil {
		return nil, err
	}
	converter := NewConverter(u.Model, u.Handlers)
	can, err := converter.Convert(term)
	if err != nil {
		return nil, err
	}
	cu := ComposeExpression(can, false)
	if value.Value == nil {
		return NewPair(Zero, cu), nil
	} else {
		return NewPair(value.Value.Multiply(can.Value), cu), nil
	}
}

func (u *UcumEssenceService) Convert(value *Decimal, sourceUnit, destUnit string) (*Decimal, error) {
	if value == nil {
		return nil, fmt.Errorf("Convert: value must not nil")
	}
	if sourceUnit == "" {
		return nil, fmt.Errorf("Convert: sourceUnit must not be empty")
	}
	if destUnit == "" {
		return nil, fmt.Errorf("Convert: destUnit must not be empty")
	}
	if sourceUnit == destUnit {
		return value, nil
	}
	converter := NewConverter(u.Model, u.Handlers)
	srcEp, err := NewExpressionParser(u.Model).Parse(sourceUnit)
	if err != nil {
		return nil, err
	}
	drcEp, err := NewExpressionParser(u.Model).Parse(destUnit)
	if err != nil {
		return nil, err
	}
	src, err := converter.Convert(srcEp)
	if err != nil {
		return nil, err
	}
	dst, err := converter.Convert(drcEp)
	if err != nil {
		return nil, err
	}
	s := ComposeExpression(src, false)
	d := ComposeExpression(dst, false)
	if s != d {
		return nil, fmt.Errorf("Unable to convert between units " + sourceUnit + " and " + destUnit + " as they do not have matching canonical forms (" + s + " and " + d + " respectively)")
	}
	canValue := value.Multiply(src.Value)
	dr := canValue.Divide(dst.Value)
	return dr, nil
}

func (u *UcumEssenceService) Multiply(o1, o2 *Pair) (*Pair, error) {
	res := NewPair(o1.Value.Multiply(o2.Value), o1.Code+"."+o2.Code)
	return u.GetCanonicalForm(res)
}

func (u *UcumEssenceService) GetCommonDisplay(code string) string {
	code = strings.Replace(code, "[", "", -1)
	code = strings.Replace(code, "]", "", -1)
	return code
}

//UcumEssenceService=======================================================
type UcumValidator struct {
	Model    *UcumModel
	Result   []string
	Handlers *Registry
}

func NewUcumValidator(model *UcumModel, handlers *Registry) *UcumValidator {
	v := &UcumValidator{}
	v.Model = model
	if handlers == nil {
		handlers = NewRegistry()
	}
	v.Handlers = handlers
	return v
}

func (v *UcumValidator) Validate() []string {
	v.Result = make([]string, 0)
	v.checkCodes()
	v.checkUnits()
	return v.Result
}

func (v *UcumValidator) checkCodes() {
	for _, u := range v.Model.BaseUnits {
		v.checkUnitCode(u.Code, true)
	}
	for _, u := range v.Model.DefinedUnits {
		v.checkUnitCode(u.Code, true)
	}
}

func (v *UcumValidator) checkUnits() {
	for _, u := range v.Model.DefinedUnits {
		if !u.IsSpecial {
			v.checkUnitCode(u.Value.Unit, false)
		} else if !v.Handlers.Exists(u.Code) {
			v.Result = append(v.Result, "No handler for "+u.Code)
		}
	}
}

func (v *UcumValidator) checkUnitCode(code string, primary bool) {
	term, err := NewExpressionParser(v.Model).Parse(code)
	if err != nil {
		v.Result = append(v.Result, err.Error())
		return
	}
	c := ComposeExpression(term, false)
	if c != code {
		v.Result = append(v.Result, "Round trip failed: "+code+" -> "+c)
	}
	NewConverter(v.Model, v.Handlers).Convert(term)
	if primary {
		isInBrack := false
		nonDigits := false
		for i := 0; i < len(code); i++ {
			ch := code[i]
			if ch == '[' {
				if isInBrack {
					v.Result = append(v.Result, "nested '[' detected")
				} else {
					isInBrack = true
				}
			}
			if ch == ']' {
				if !isInBrack {
					v.Result = append(v.Result, "']' without '[' detected")
				} else {
					isInBrack = false
				}
			}
			nonDigits = nonDigits || !(ch >= '0' && ch <= '9')
			if ch >= '0' && ch <= '9' && !isInBrack && nonDigits {
				v.Result = append(v.Result, "code "+code+" is ambiguous because it has digits outside []")
			}
		}
	}
}
