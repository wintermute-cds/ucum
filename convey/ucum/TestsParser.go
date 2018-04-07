package ucum

import (
	"encoding/xml"
	"os"
	"io/ioutil"
)

type TestStructures struct {
	ValidationCases            []XMLValidationCase
	DisplayNameGenerationCases []XMLDisplayNameGenerationCase
	conversionCases            []XMLConversionCase
	multiplicationCases        []XMLMultiplicationCase
}


type XMLUcumTests struct {
	Validations            XMLValidation            `xml:"validation"`
	DisplayNameGenerations XMLDisplayNameGeneration `xml:"displayNameGeneration"`
	Conversions            XMLConversion            `xml:"conversion"`
	Multiplications        XMLMultiplication        `xml:"multiplication"`
}

type XMLValidation struct {
	Cases []XMLValidationCase `xml:"case"`
}

type XMLValidationCase struct {
	Id     string `xml:"id,attr"`
	Unit   string `xml:"unit,attr"`
	Valid  string `xml:"valid,attr"`
	Reason string `xml:"reason,attr"`
}

type XMLDisplayNameGeneration struct {
	Cases []XMLDisplayNameGenerationCase `xml:"case"`
}

type XMLDisplayNameGenerationCase struct {
	Id      string `xml:"id,attr"`
	Unit    string `xml:"unit,attr"`
	Display string `xml:"display,attr"`
}

type XMLConversion struct {
	Cases []XMLConversionCase `xml:"case"`
}
type XMLConversionCase struct {
	Id      string `xml:"id,attr"`
	Value   string `xml:"value,attr"`
	SrcUnit string `xml:"srcUnit,attr"`
	DstUnit string `xml:"dstUnit,attr"`
	Outcome string `xml:"outcome,attr"`
}

type XMLMultiplication struct {
	Cases []XMLMultiplicationCase `xml:"case"`
}

type XMLMultiplicationCase struct {
	Id   string `xml:"id,attr"`
	V1   string `xml:"v1,attr"`
	U1   string `xml:"u1,attr"`
	V2   string `xml:"v2,attr"`
	U2   string `xml:"u2,attr"`
	VRes string `xml:"vRes,attr"`
	URes string `xml:"uRes,attr"`
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
	t.DisplayNameGenerationCases = make([]XMLDisplayNameGenerationCase, 0)
	t.conversionCases = make([]XMLConversionCase, 0)
	t.multiplicationCases = make([]XMLMultiplicationCase, 0)
	for _, xmlItem := range xmlTest.Validations.Cases {
		validationCase := xmlItem
		t.validationCases = append(t.validationCases, validationCase)
	}
	for _, xmlItem := range xmlTest.DisplayNameGenerations.Cases {
		displayNameGenerationCase := xmlItem
		t.DisplayNameGenerationCases = append(t.DisplayNameGenerationCases, displayNameGenerationCase)
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
