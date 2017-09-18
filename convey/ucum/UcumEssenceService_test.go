package ucum

import (
	"UCUM_Golang/ucum"
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

var test string
var service *ucum.UcumEssenceService

type TestStructures struct {
	validationCases            []XMLValidationCase
	displayNameGenerationCases []XMLDisplayNameGenerationCase
	conversionCases            []XMLConversionCase
	multiplicationCases        []XMLMultiplicationCase
}

func TestService(t *testing.T) {
	var testStructures *TestStructures
	var err error
	Convey("Service-creation", t, func() {
		Convey("Service-creation", func() {
			definitions := os.Getenv("GOPATH") + "/src/UCUM_Golang/ucum/terminology_data/ucum-essence.xml"
			service, err = ucum.GetInstanceOfUcumEssenceService(definitions)
			So(err, ShouldBeNil)
			So(service, ShouldNotBeNil)
		})
		Convey("First test run", func() {
			test = os.Getenv("GOPATH") + "/src/UCUM_Golang/convey/resources/UcumFunctionalTests.xml"
			testStructures, err = UnmarshalTerminology(test)
			So(err, ShouldBeNil)
			So(testStructures, ShouldNotBeNil)
		})
		RunValidationTest(t, testStructures, "Validation test")
		RunDisplayNameGenerationTest(t, testStructures, "DisplayNameGenerationTest")
		RunConversionTest(t, testStructures, "ConversionTest")
		RunMultiplicationTest(t, testStructures, "RunMultiplicationTest")
	})
}

func RunValidationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.validationCases {
			Convey(v.Id+": "+v.Unit, func() {
				validated, _ := service.Validate(v.Unit)
				So(validated, ShouldEqual, v.Valid == "true")
			})
		}
	})
}

func RunDisplayNameGenerationTest(t *testing.T, testStructures *TestStructures, name string) {
	Convey(name, func() {
		for _, v := range testStructures.displayNameGenerationCases {
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
			})
		}
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

func UnmarshalTerminology(xmlFileName string) (*TestStructures, error) {
	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var xmlTest XMLUcumTests
	xml.Unmarshal(byteValue, &xmlTest)
	//decoder := xml.NewDecoder(xmlFile)
	//if err := decoder.Decode(xmlTest); err!=nil {
	//	return nil, err
	//}
	t := &TestStructures{}
	t.validationCases = make([]XMLValidationCase, 0)
	t.displayNameGenerationCases = make([]XMLDisplayNameGenerationCase, 0)
	t.conversionCases = make([]XMLConversionCase, 0)
	t.multiplicationCases = make([]XMLMultiplicationCase, 0)
	for _, xmlItem := range xmlTest.Validations.Cases {
		validationCase := xmlItem
		t.validationCases = append(t.validationCases, validationCase)
	}
	for _, xmlItem := range xmlTest.DisplayNameGenerations.Cases {
		displayNameGenerationCase := xmlItem
		t.displayNameGenerationCases = append(t.displayNameGenerationCases, displayNameGenerationCase)
	}
	for _, xmlItem := range xmlTest.Conversions.Cases {
		conversionCase := xmlItem
		t.conversionCases = append(t.conversionCases, conversionCase)
	}
	for _, xmlItem := range xmlTest.Multiplications.Cases {
		multiplicationCase := xmlItem
		t.multiplicationCases = append(t.multiplicationCases, multiplicationCase)
	}

	return t, nil
}
