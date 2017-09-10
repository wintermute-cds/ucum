package ucum

import (
	"time"
	"UCUM_Golang/ucum/special"
	"os"
	"fmt"
)

type UcumService interface {
	/**
	 * return Ucum Identification details for the version in use
	 */
	UcumIdentification() UcumVersionDetails
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
	Search(kind ConceptKind, text string, isRegex bool)[]Concepter
	/**
	 * return a list of the defined types of units in this UCUM version
	 *
	 * @return
	 */
	GetProperties()[]string
	/**
	 * validate whether a unit code are valid UCUM units
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	Validate(unit string)string
	/**
	 * given a unit, return a formal description of what the units stand for using
	 * full names
	 * @param units the unit code
	 * @return formal description
	 * @throws UcumException
	 * @throws OHFException
	 */
	Analyse(unit string)(string,error)
	/**
	 * validate whether a units are valid UCUM units and additionally require that the
	 * units from a particular property
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	ValidateInProperty(unit, property string)string
	/**
	 * validate whether a units are valid UCUM units and additionally require that the
	 * units match a particular base canonical unit
	 *
	 * @param units - the unit code to check
	 * @return nil if valid, or an error message describing the problem
	 */
	ValidateCanonicalUnits(unit,  canonical string)string
	/**
	 * given a set of units, return their canonical form
	 * @param unit
	 * @return the canonical form
	 * @throws UcumException
	 * @throws OHFException
	 */
	GetCanonicalUnits(unit error)(string, error)
	/**
	   * given two pairs of units, return true if they sahre the same canonical base
	   *
	   * @param units1
	   * @param units2
	   * @return
	   * @throws UcumException
	   * @
	   */
	IsComparable( units1,  units2 string)(bool, error)
	/**
	 * for a given canonical unit, return all the defined units that have the
	 * same canonical unit.
	 *
	 * @param code
	 * @return
	 * @throws UcumException
	 * @throws OHFException
	 */
	GetDefinedForms(code string)([]*DefinedUnit,error)
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
	Convert(value *Decimal, sourceUnit, destUnit string)(*Decimal, error)
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
	Multiply( o1,  o2 *Pair)(*Pair, error)
	/**
		 * given a set of UCUM units, return a likely preferred human dense form
		 *
		 * SI units - as is.
		 * Other units - improved by manual fixes, or the removal of []
		 *
		 * @param code
		 * @return the preferred human display form
		 */
	GetCommonDisplay(code string)string;

}

// UcumVersionDetails======================================================
type UcumVersionDetails struct{
	ReleaseDate *time.Time
	Version string
}

func NewUcumVersionDetails(releaseDate *time.Time, version string)*UcumVersionDetails{
	r := &UcumVersionDetails{}
	r.ReleaseDate = releaseDate
	r.Version = version
	return r
}
//UcumEssenceService=======================================================
const UCUM_OID = "2.16.840.1.113883.6.8"
var xmlFileName = "../terminology_data/ucum-essence.xml"

type UcumEssenceService struct{
	Model *UcumModel
	Handlers *special.Registry
}

var instanceOfOpenEhrTerminologyService *UcumEssenceService

func GetInstanceOfOpenEhrTerminologyService()(*UcumEssenceService, error){
	if instanceOfOpenEhrTerminologyService==nil{
		instanceOfOpenEhrTerminologyService = new(UcumEssenceService)
		xmlFile, err := os.Open(xmlFileName)
		if err != nil {
			return nil, err
		}
		defer xmlFile.Close()
		d := new(DefinitionParser)
		instanceOfOpenEhrTerminologyService.Model,err = d.UnmarshalTerminology(xmlFile)
		if err != nil {
			return nil, err
		}
	}
	return instanceOfOpenEhrTerminologyService, nil
}

func (u *UcumEssenceService)UcumIdentification() UcumVersionDetails{

}
func (u *UcumEssenceService)ValidateUCUM() []string{

}
func (u *UcumEssenceService)Search(kind ConceptKind, text string, isRegex bool)[]Concepter{

}
func (u *UcumEssenceService)GetProperties()[]string{

}
func (u *UcumEssenceService)Validate(unit string)string{

}
func (u *UcumEssenceService)Analyse(unit string)(string,error){

}
func (u *UcumEssenceService)ValidateInProperty(unit, property string)string{

}
func (u *UcumEssenceService)ValidateCanonicalUnits(unit,  canonical string)string{

}
func (u *UcumEssenceService)GetCanonicalUnits(unit error)(string, error){

}
func (u *UcumEssenceService)IsComparable( units1,  units2 string)(bool, error){

}
func (u *UcumEssenceService)GetDefinedForms(code string)([]*DefinedUnit,error){

}
func (u *UcumEssenceService)GetCanonicalForm(value *Pair) (*Pair, error){

}
func (u *UcumEssenceService)Convert(value *Decimal, sourceUnit, destUnit string)(*Decimal, error){

}
func (u *UcumEssenceService)Multiply( o1,  o2 *Pair)(*Pair, error){

}
func (u *UcumEssenceService)GetCommonDisplay(code string)string;{

}
