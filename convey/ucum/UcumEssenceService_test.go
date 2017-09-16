package ucum

import (
	"testing"
	"os"
	"UCUM_Golang/ucum"
	. "github.com/smartystreets/goconvey/convey"
	"encoding/xml"
	"io/ioutil"
)

var test string
var service *ucum.UcumEssenceService

type TestStructures struct {
	validationCases []XMLValidationCase
	displayNameGenerationCases []XMLDisplayNameGenerationCase
	conversionCases []XMLConversionCase
	multiplicationCases []XMLMultiplicationCase
}

func TestService(t *testing.T){
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
		RunValidationTest(t, testStructures, "Validation test 1")
		Convey("Second test run", func() {
			test = os.Getenv("GOPATH") + "/src/UCUM_Golang/convey/resources/UcumFunctionalTests.2.xml"
			testStructures, err = UnmarshalTerminology(test)
			So(err, ShouldBeNil)
			So(testStructures, ShouldNotBeNil)
		})
		RunValidationTest(t, testStructures, "Validation test 2")
	})
}

func RunValidationTest(t *testing.T, testStructures *TestStructures, name string){
	Convey(name, func() {
		for _,v := range testStructures.validationCases{
			Convey(v.Id + ": " + v.Unit, func() {
				validated, _ := service.Validate(v.Unit)
				So(validated, ShouldEqual, v.Valid == "true")
			})
		}
	})
}

func UnmarshalTerminology(xmlFileName string)(*TestStructures, error){
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
	t.validationCases = make([]XMLValidationCase,0)
	t.displayNameGenerationCases = make([]XMLDisplayNameGenerationCase,0)
	t.conversionCases = make([]XMLConversionCase,0)
	t.multiplicationCases = make([]XMLMultiplicationCase,0)
	for _, xmlItem := range xmlTest.Validations.Cases {
		validationCase := xmlItem
		t.validationCases = append (t.validationCases,validationCase)
	}
	for _, xmlItem := range xmlTest.DisplayNameGenerations.Cases {
		displayNameGenerationCase := xmlItem
		t.displayNameGenerationCases = append (t.displayNameGenerationCases,displayNameGenerationCase)
	}
	for _, xmlItem := range xmlTest.Conversions.Cases {
		conversionCase := xmlItem
		t.conversionCases = append (t.conversionCases,conversionCase)
	}
	for _, xmlItem := range xmlTest.Multiplications.Cases {
		multiplicationCase := xmlItem
		t.multiplicationCases = append (t.multiplicationCases,multiplicationCase)
	}

	return t, nil
}


