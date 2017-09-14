package ucum

import (
	"testing"
	"os"
	"UCUM_Golang/ucum"
	"github.com/smartystreets/goconvey/convey"
	"encoding/xml"
	"io/ioutil"
	"fmt"
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
	convey.Convey("Service-creation", t, func() {
		definitions := os.Getenv("GOPATH") + "/src/UCUM_Golang/ucum/terminology_data/ucum-essence.xml"
		service, err = ucum.GetInstanceOfUcumEssenceService(definitions)
		convey.So(err, convey.ShouldBeNil)
		convey.So(service, convey.ShouldNotBeNil)
	})
	convey.Convey("First test run", t, func() {
		test = os.Getenv("GOPATH") + "/src/UCUM_Golang/convey/resources/UcumFunctionalTests.2.xml"
		testStructures, err = UnmarshalTerminology(test)
		convey.So(err, convey.ShouldBeNil)
		convey.So(testStructures, convey.ShouldNotBeNil)
	})
	RunValidationTest(t, testStructures)
	convey.Convey("First test run", t, func() {
		test = os.Getenv("GOPATH") + "/src/UCUM_Golang/convey/resources/UcumFunctionalTests.2.xml"
		testStructures, err = UnmarshalTerminology(test)
		convey.So(err, convey.ShouldBeNil)
		convey.So(testStructures, convey.ShouldNotBeNil)
	})
}

func RunValidationTest(t *testing.T, testStructures *TestStructures){
	convey.Convey("Validation test", t, func() {
		for _,v := range testStructures.validationCases{
			msg := service.Validate(v.Unit)
			fmt.Println(msg)
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


