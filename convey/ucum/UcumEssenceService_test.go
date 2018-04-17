package ucum

import (
	"github.com/bertverhees/ucum"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"fmt"
	"reflect"
)

var test string
var service *ucum.UcumEssenceService

func TestService(t *testing.T) {
	var testStructures *TestStructures
	var err error
	Convey("Service-creation", t, func() {
		Convey("Service-creation", func() {
			definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
			service, err = ucum.GetInstanceOfUcumEssenceService(definitions)
			So(err, ShouldBeNil)
			So(service, ShouldNotBeNil)
		})
		Convey("First test run", func() {
			test = os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/convey/resources/UcumFunctionalTests.xml"
			testStructures, err = UnmarshalTerminology(test)
			So(err, ShouldBeNil)
			So(testStructures, ShouldNotBeNil)
		})
		RunValidationTest(t, testStructures, "Validation test")
		RunDisplayNameGenerationTest(t, testStructures, "DisplayNameGenerationTest")
		RunConversionTest(t, testStructures, "ConversionTest")
		RunMultiplicationTest(t, testStructures, "RunMultiplicationTest")
		RunUcumIdentificationTest(t, testStructures, "RunUcumIdentificationTest")
		RunUcumValidateUCUMTest(t, testStructures, "RunUcumValidateUCUMTest")
		RunSearchPrefixTests(t, testStructures, "RunSearchPrefixTests")
		RunSearchBaseUnitsTests(t, testStructures, "RunSearchBaseUnitsTests")
		RunSearchUnitsTests(t, testStructures, "RunSearchUnitsTests")
		RunGetPropertiesTests(t, testStructures, "RunGetPropertiesTests")
		RunValidateInPropertyTests(t, testStructures, "RunValidateInPropertyTests")
		RunValidateCanonicalUnitsTests(t, testStructures, "RunValidateCanonicalUnitsTests")
		RunGetCanonicalUnitsTests(t, testStructures, "RunGetCanonicalUnitsTests")
		RunGetDefinedFormsTests(t, testStructures, "RunGetDefinedFormsTests")
		RunIsComparableTests(t, testStructures, "RunIsComparableTests")
	})
}

func RunIsComparableTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
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


func RunGetDefinedFormsTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		validated, err := service.GetDefinedForms("mm")
		So(err, ShouldBeNil)
		So(len(validated), ShouldEqual, 0)
		validated, err = service.GetDefinedForms("rad")
		So(err, ShouldBeNil)
		So(len(validated), ShouldBeGreaterThan, 0)
	})
}


func RunGetCanonicalUnitsTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		validated, err := service.GetCanonicalUnits("mm")
		So(err, ShouldBeNil)
		So(validated, ShouldEqual, "m")
		validated = service.ValidateInProperty("cm", "length" )
		So(validated, ShouldBeEmpty)
	})
}


func RunValidateInPropertyTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		validated := service.ValidateInProperty("mm", "number" )
		So(validated, ShouldEqual, "unit mm is of the property type length (m), not number as required.")
		validated = service.ValidateInProperty("cm", "length" )
		So(validated, ShouldBeEmpty)
	})
}

func RunValidateCanonicalUnitsTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		validated := service.ValidateCanonicalUnits("mm", "l" )
		So(validated, ShouldEqual, "unit mm has the base units m, not l as required.")
		validated = service.ValidateCanonicalUnits("cm", "m" )
		So(validated, ShouldBeEmpty)
	})
}

func RunGetPropertiesTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		list := service.GetProperties()
		So(len(list), ShouldBeGreaterThan, 300)
	})
}

func RunSearchUnitsTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
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


func RunSearchBaseUnitsTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
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

func RunSearchPrefixTests(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
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

func RunUcumValidateUCUMTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		s := service.ValidateUCUM()
		for _, e := range s{
			fmt.Println(e)
		}
	})
}


func RunUcumIdentificationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		So(service.UcumIdentification().Version, ShouldNotBeEmpty)
		So(service.UcumIdentification().ReleaseDate.String(), ShouldNotBeEmpty)
	})
}

func RunValidationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.ValidationCases {
			Convey(v.Id+": "+v.Unit, func() {
				validated, _ := service.Validate(v.Unit)
				So(validated, ShouldEqual, v.Valid == "true")
			})
		}
	})
}

func RunDisplayNameGenerationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.DisplayNameGenerationCases {
			Convey(v.Id+": "+v.Unit, func() {
				analysed, _ := service.Analyse(v.Unit)
				So(analysed, ShouldEqual, v.Display)
			})
		}
	})
}

func RunConversionTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.conversionCases {
			Convey(v.Id+": "+v.Value, func() {
				d, err := ucum.NewDecimal(v.Value)
				So(err, ShouldBeNil)
				o, err := ucum.NewDecimal(v.Outcome)
				So(err, ShouldBeNil)
				res, _ := service.Convert(d, v.SrcUnit, v.DstUnit)
				So(res.AsDecimal(), ShouldEqual, o.AsDecimal())
				fmt.Println(v.SrcUnit+":"+v.DstUnit)
				fmt.Println(d.AsScientific()+":"+d.StringFixed(int32(d.GetPrecision()))+":"+d.String())
				fmt.Println(o.AsScientific()+":"+o.StringFixed(int32(o.GetPrecision()))+":"+o.String())
			})
		}
	})
}

func TestConvert(t *testing.T){
	Convey("TestConvert", t,func() {
		decimal := ucum.NewFromInt64Precision(63, -1, 1)
		fmt.Println(decimal)
		fmt.Println(decimal.Exponent())
		fmt.Println(decimal.GetPrecision())
		fmt.Println(decimal.AsInteger())
		fmt.Println(decimal.GetValue().String())
		fmt.Println(decimal.GetValue().IsInt64())
		fmt.Println("--------")
		fmt.Println("Income")
		fmt.Println(decimal)
		definitions := os.Getenv("GOPATH") + "/src/github.com/bertverhees/ucum/terminology_data/ucum-essence.xml"
		service, err := ucum.GetInstanceOfUcumEssenceService(definitions)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		result, err := service.Convert(decimal, "s.mm-1", "s.m-1")
		if err != nil {
			fmt.Errorf(err.Error())
		}
		So(decimal.Multiply(ucum.NewFromInt(1000,0)), ShouldEqual, result)
		fmt.Println(result)
		fmt.Println(result.Exponent())
		fmt.Println(result.GetPrecision())
		fmt.Println(result.AsInteger())
		fmt.Println(result.GetValue().String())
		fmt.Println(result.GetValue().IsInt64())
		fmt.Println("--------")
	})
}

func RunMultiplicationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.multiplicationCases {
			Convey(v.Id, func() {
				d, err := ucum.NewDecimal(v.V1)
				So(err, ShouldBeNil)
				o1 := ucum.NewPair(d, v.U1)
				d, err = ucum.NewDecimal(v.V2)
				So(err, ShouldBeNil)
				o2 := ucum.NewPair(d, v.U2)
				o3, err := service.Multiply(o1, o2)
				So(err, ShouldBeNil)
				d, err = ucum.NewDecimal(v.VRes)
				test := o3.Value.ComparesTo(d)
				So(test, ShouldEqual, 0)
			})
		}
	})
}

