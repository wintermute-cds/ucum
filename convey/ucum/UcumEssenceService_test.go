package ucum

import (
	"github.com/bertverhees/ucum"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"fmt"
	"reflect"
	"github.com/bertverhees/ucum/decimal"
	"strings"
)

var test string
var service *ucum.UcumEssenceService
var testStructures *TestStructures

func InitService() {
	var err error
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err = ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf("Error instantiating service:"+err.Error())
	}
	test = os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/convey/resources/UcumFunctionalTests.xml"
	testStructures, err = UnmarshalTerminology(test)
	if err != nil {
		fmt.Errorf("Error unmarshaling:"+err.Error())
	}
}

func TestIsComparableTests(t *testing.T) {
	InitService()
	Convey("TestIsComparableTests", t, func() {
		validated, err := service.IsComparable("mm", "rad")
		So(err, ShouldBeNil)
		So(validated, ShouldBeFalse)
		validated, err = service.IsComparable("mm", "cm")
		So(err, ShouldBeNil)
		So(validated, ShouldBeTrue)
		validated, err = service.IsComparable("mm", "m")
		So(err, ShouldBeNil)
		So(validated, ShouldBeTrue)
	})
}


func TestGetDefinedFormsTests(t *testing.T) {
	InitService()
	Convey("TestGetDefinedFormsTests", t, func() {
		validated, err := service.GetDefinedForms("mm")
		So(err, ShouldBeNil)
		So(len(validated), ShouldEqual, 0)
		validated, err = service.GetDefinedForms("rad")
		So(err, ShouldBeNil)
		So(len(validated), ShouldBeGreaterThan, 0)
	})
}


func TestGetCanonicalUnitsTests(t *testing.T) {
	InitService()
	Convey("TestGetCanonicalUnitsTests", t, func() {
		validated, err := service.GetCanonicalUnits("mm")
		So(err, ShouldBeNil)
		So(validated, ShouldEqual, "m")
		validated = service.ValidateInProperty("cm", "length" )
		So(validated, ShouldBeEmpty)
	})
}


func TestValidateInPropertyTests(t *testing.T) {
	InitService()
	Convey("TestValidateInPropertyTests", t, func() {
		validated := service.ValidateInProperty("mm", "number" )
		So(validated, ShouldEqual, "unit mm is of the property type length (m), not number as required.")
		validated = service.ValidateInProperty("cm", "length" )
		So(validated, ShouldBeEmpty)
	})
}

func TestValidateCanonicalUnitsTests(t *testing.T) {
	InitService()
	Convey("TestValidateCanonicalUnitsTests", t, func() {
		validated := service.ValidateCanonicalUnits("mm", "l" )
		So(validated, ShouldEqual, "unit mm has the base units m, not l as required.")
		validated = service.ValidateCanonicalUnits("cm", "m" )
		So(validated, ShouldBeEmpty)
	})
}

func TestGetPropertiesTests(t *testing.T) {
	InitService()
	Convey("TestGetPropertiesTests", t, func() {
		list := service.GetProperties()
		So(len(list), ShouldBeGreaterThan, 300)
	})
}

func TestSearchUnitsTests(t *testing.T) {
	InitService()
	Convey("TestSearchUnitsTests", t, func() {
		list, err := service.Search(ucum.UNIT, "sr", false)
		So(err, ShouldBeNil)
		p1 := list[0]
		list, err = service.Search(ucum.UNIT, "SR", false)
		So(err, ShouldBeNil)
		p2 := list[0]
		So(reflect.DeepEqual(p1,p2), ShouldBeTrue)
		list, err = service.Search(ucum.UNIT, "steradian", false)
		So(err, ShouldBeNil)
		p3 := list[0]
		So(reflect.DeepEqual(p1,p3), ShouldBeTrue)
		list, err = service.Search(ucum.UNIT, "solid angle", false)
		So(err, ShouldBeNil)
		p4 := list[0]
		So(reflect.DeepEqual(p1,p4), ShouldBeTrue)
		list, err = service.Search(ucum.UNIT, "^m([a-z]+)r", true)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 9)
		list, err = service.Search(ucum.UNIT, "m([a-z]+)r", true)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 31)
	})
}


func TestSearchBaseUnitsTests(t *testing.T) {
	InitService()
	Convey("TestSearchBaseUnitsTests", t, func() {
		list, err := service.Search(ucum.BASEUNIT, "meter", false)
		So(err, ShouldBeNil)
		p1 := list[0]
		list, err = service.Search(ucum.BASEUNIT, "length", false)
		So(err, ShouldBeNil)
		p2 := list[0]
		So(reflect.DeepEqual(p1,p2), ShouldBeTrue)
		list, err = service.Search(ucum.BASEUNIT, "m", false)
		So(err, ShouldBeNil)
		p3 := list[0]
		So(reflect.DeepEqual(p1,p3), ShouldBeTrue)
		list, err = service.Search(ucum.BASEUNIT, "M", false)
		So(err, ShouldBeNil)
		p4 := list[0]
		So(reflect.DeepEqual(p1,p4), ShouldBeTrue)
		list, err = service.Search(ucum.BASEUNIT, "L", false)
		So(err, ShouldBeNil)
		p6 := list[0]
		So(reflect.DeepEqual(p1,p6), ShouldBeTrue)
		list, err = service.Search(ucum.BASEUNIT, "^m([a-z]+)r", true)
		So(err, ShouldBeNil)
		p5 := list[0]
		So(reflect.DeepEqual(p1,p5), ShouldBeTrue)
		So(len(list), ShouldEqual, 1)
		f := false
		for _, s := range list {
			if s.GetNames()[0] == "meter" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
		list, err = service.Search(ucum.BASEUNIT, "m([a-z]+)r", true)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)
		f = false
		for _, s := range list {
			if s.GetNames()[0] == "meter" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
		f = false
		for _, s := range list {
			if s.(*ucum.BaseUnit).Property == "temperature" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
	})
}

func TestSearchPrefixTests(t *testing.T) {
	InitService()
	Convey("TestSearchPrefixTests", t, func() {
		list, err := service.Search(ucum.PREFIX, "micro", false)
		So(err, ShouldBeNil)
		p1 := list[0]
		list, err = service.Search(ucum.PREFIX, "μ", false)
		So(err, ShouldBeNil)
		p2 := list[0]
		So(reflect.DeepEqual(p1,p2), ShouldBeTrue)
		list, err = service.Search(ucum.PREFIX, "u", false)
		So(err, ShouldBeNil)
		p3 := list[0]
		So(reflect.DeepEqual(p1,p3), ShouldBeTrue)
		list, err = service.Search(ucum.PREFIX, "U", false)
		So(err, ShouldBeNil)
		p4 := list[0]
		So(reflect.DeepEqual(p1,p4), ShouldBeTrue)
		list, err = service.Search(ucum.PREFIX, "^m([a-z]+)o", true)
		So(err, ShouldBeNil)
		p5 := list[0]
		So(reflect.DeepEqual(p1,p5), ShouldBeTrue)
		So(len(list), ShouldEqual, 1)
		f := false
		for _, s := range list {
			if s.GetNames()[0] == "micro" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
		list, err = service.Search(ucum.PREFIX, "m([a-z]+)o", true)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)
		f = false
		for _, s := range list {
			if s.GetNames()[0] == "micro" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
		f = false
		for _, s := range list {
			if s.GetNames()[0] == "femto" {
				f = true
			}
		}
		So(f, ShouldBeTrue)
	})
}

func TestUcumValidateUCUMTest(t *testing.T) {
	InitService()
	Convey("TestUcumValidateUCUMTest", t, func() {
		s := service.ValidateUCUM()
		for _, e := range s{
			fmt.Println(e)
		}
	})
}


func TestUcumIdentificationTest(t *testing.T) {
	InitService()
	Convey("TestUcumIdentificationTest", t, func() {
		So(service.UcumIdentification().Version, ShouldNotBeEmpty)
		So(service.UcumIdentification().ReleaseDate.String(), ShouldNotBeEmpty)
	})
}

func TestValidationTest(t *testing.T) {
	InitService()
	Convey("TestValidationTest", t, func() {
		for _, v := range testStructures.ValidationCases {
			Convey(v.Id+": "+v.Unit, func() {
				validated, _ := service.Validate(v.Unit)
				So(validated, ShouldEqual, v.Valid == "true")
			})
		}
	})
}

func TestDisplayNameGenerationTest(t *testing.T) {
	InitService()
	Convey("TestDisplayNameGenerationTest", t, func() {
		for _, v := range testStructures.DisplayNameGenerationCases {
			Convey(v.Id+": "+v.Unit, func() {
				analysed, _ := service.Analyse(v.Unit)
				So(analysed, ShouldEqual, v.Display)
			})
		}
	})
}

func TestConversionTest(t *testing.T) {
	InitService()
	Convey("TestConversionTest", t, func() {
		for _, v := range testStructures.conversionCases {
			Convey(v.Id+": "+v.Value, func() {
				d, err := decimal.NewFromString(v.Value)
				So(err, ShouldBeNil)
				o, err := decimal.NewFromString(v.Outcome)
				So(err, ShouldBeNil)
				res, _ := service.Convert(d, v.SrcUnit, v.DstUnit)
				So(res.Cmp(o), ShouldEqual, 0)
			})
		}
	})
}

func TestConvert(t *testing.T){
	InitService()
	Convey("TestConvert", t,func() {
		dec := decimal.New(63, -1)
		definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
		service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		result, err := service.Convert(dec, "s.mm-1", "s.m-1")
		if err != nil {
			fmt.Errorf(err.Error())
		}
		So(dec.Mul(decimal.New(1000,0)), ShouldEqual, result)
	})
}

func TestMultiplicationTest(t *testing.T) {
	InitService()
	Convey("TestMultiplicationTest", t, func() {
		for _, v := range testStructures.multiplicationCases {
			Convey(v.Id, func() {
				d, err := decimal.NewFromString(v.V1)
				So(err, ShouldBeNil)
				o1 := ucum.NewPair(d, v.U1)
				d, err = decimal.NewFromString(v.V2)
				So(err, ShouldBeNil)
				o2 := ucum.NewPair(d, v.U2)
				o3, err := service.Multiply(o1, o2)
				So(err, ShouldBeNil)
				d, err = decimal.NewFromString(v.VRes)
				test := o3.Value.Cmp(d)
				So(test, ShouldEqual, 0)
			})
		}
	})
}

func TestListAllClasses(t *testing.T) {
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	classes := service.ListAllClasses()
	if len(classes) !=20 {
		t.Error("Not the right number of classes found")
	}
}

func TestSearchClasses(t *testing.T) {
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	classes := service.SearchClass("len")
	if len(classes) !=2 {
		t.Error("Not the right number of classes found")
	}
}

func TestListAllProperties(t *testing.T) {
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	classes := service.ListAllProperties()
	if len(classes) !=99 {
		t.Error("Not the right number of properties found")
	}
}

func TestSearchProperties(t *testing.T) {
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	classes := service.SearchProperty("len")
	if len(classes) !=3 {
		t.Error("Not the right number of classes found")
	}
}

func TestGetClassInfo(t *testing.T) {
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	class := service.GetClassInfo("us-lengths")
	if !strings.HasPrefix(class.Description, "The older U.S. units according"){
		t.Error("Class description wrong")
	}
}

func TestFilterDefinedUnits(t *testing.T){
	definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
	service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	definedUnits := service.FilterDefinedModels("us-lengths", "", false, false, false, false, false, false)
	if len(definedUnits) != 17 {
		t.Error("Not the right number of defined units found on filtering on class")
	}
	definedUnits = service.FilterDefinedModels("us-lengths", "length", false, false, false, false, false, false)
	if len(definedUnits) != 12 {
		t.Error("Not the right number of defined units found on filtering on class and property")
	}
	definedUnits = service.FilterDefinedModels("", "length", false, false, false, false, false, false)
	if len(definedUnits) != 43 {
		t.Error("Not the right number of defined units found on filtering on class and property")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "", false, false, false, false, false, false)
	if len(definedUnits) != 25 {
		t.Error("Not the right number of defined units found on filtering on class")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "", true, false, false, false, false, false)
	if len(definedUnits) != 17 {
		t.Error("Not the right number of defined units found on filtering on class and metric")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "time", true, false, false, false, false, false)
	if len(definedUnits) != 12 {
		t.Error("Not the right number of defined units found on filtering on class, property and metric")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "time", false, false, false, false, false, false)
	if len(definedUnits) != 12 {
		t.Error("Not the right number of defined units found on filtering on class and property")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "length", false, false, false, false, false, false)
	if len(definedUnits) != 2 {
		t.Error("Not the right number of defined units found on filtering on class and property")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "length", true, false, false, false, false, false)
	if len(definedUnits) != 1 {
		t.Error("Not the right number of defined units found on filtering on class and property and metric")
	}
	definedUnits = service.FilterDefinedModels("iso1000", "length", true, true, false, false, false, false)
	if len(definedUnits) != 1 {
		t.Error("Not the right number of defined units found on filtering on class and property and metric")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, false, false, false, false)
	if len(definedUnits) != 89 {
		t.Error("Not the right number of defined units found on filtering on metric")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, true, false, false)
	if len(definedUnits) != 10 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, false, false, false)
	if len(definedUnits) != 79 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, true, false, false)
	if len(definedUnits) != 11 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, false, false, false)
	if len(definedUnits) != 203 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, false, false, false)
	if len(definedUnits) != 282 {
		t.Error("Not the right number of defined units found on filtering on isSpecial")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, true, false, false)
	if len(definedUnits) != 21 {
		t.Error("Not the right number of defined units found on filtering on isSpecial")
	}

	definedUnits = service.FilterDefinedModels("", "", true, true, false, false, true, false)
	if len(definedUnits) != 87 {
		t.Error("Not the right number of defined units found on filtering on metric, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, true, true, false)
	if len(definedUnits) != 10 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, false, true, false)
	if len(definedUnits) != 77 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, true, true, false)
	if len(definedUnits) != 11 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, false, true, false)
	if len(definedUnits) != 164 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, false, true, false)
	if len(definedUnits) != 241 {
		t.Error("Not the right number of defined units found on filtering on isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, true, true, false)
	if len(definedUnits) != 21 {
		t.Error("Not the right number of defined units found on filtering on isSpecial, isArbitrary")
	}

	definedUnits = service.FilterDefinedModels("", "", true, true, false, false, true, true)
	if len(definedUnits) != 2 {
		t.Error("Not the right number of defined units found on filtering on metric, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, true, true, true)
	if len(definedUnits) != 0 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, true, true, false, true, true)
	if len(definedUnits) != 2 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, true, true, true)
	if len(definedUnits) != 0 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", true, false, true, false, true, true)
	if len(definedUnits) != 39 {
		t.Error("Not the right number of defined units found on filtering on metric and isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, false, true, true)
	if len(definedUnits) != 41 {
		t.Error("Not the right number of defined units found on filtering on isSpecial, isArbitrary")
	}
	definedUnits = service.FilterDefinedModels("", "", false, false, true, true, true, true)
	if len(definedUnits) != 0 {
		t.Error("Not the right number of defined units found on filtering on isSpecial, isArbitrary")
	}

}

