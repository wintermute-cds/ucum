package ucum


type XMLUcumTests struct {
	Validations  XMLValidation  	`xml:"validation"`
	DisplayNameGenerations	XMLDisplayNameGeneration	`xml:"displayNameGeneration"`
	Conversions	XMLConversion	`xml:"conversion"`
	Multiplications	XMLMultiplication	`xml:"multiplication"`
}

type XMLValidation struct{
	Cases []XMLValidationCase					`xml:"case"`
}

type XMLValidationCase struct {
	Id string						`xml:"id,attr"`
	Unit string						`xml:"unit,attr"`
	Valid string					`xml:"valid,attr"`
	Reason string					`xml:"reason,attr"`
}

type XMLDisplayNameGeneration struct {
	Cases []XMLDisplayNameGenerationCase		`xml:"case"`
}

type XMLDisplayNameGenerationCase struct {
	Id string						`xml:"id,attr"`
	Unit string						`xml:"unit,attr"`
	Display string					`xml:"display,attr"`
}

type XMLConversion struct {
	Cases []XMLConversionCase		`xml:"case"`
}
type XMLConversionCase struct {
	Id string			`xml:"id,attr"`
	Value string		`xml:"value,attr"`
	SrcUnit string		`xml:"srcUnit,attr"`
	DstUnit string		`xml:"dstUnit,attr"`
	Outcome string		`xml:"outcome,attr"`
}


type XMLMultiplication struct {
	Cases []XMLMultiplicationCase		`xml:"case"`
}

type XMLMultiplicationCase struct {
	Id string			`xml:"id,attr"`
	V1 string			`xml:"v1,attr"`
	U1 string			`xml:"u1,attr"`
	V2 string			`xml:"v2,attr"`
	U2 string			`xml:"u2,attr"`
	VRes string			`xml:"vRes,attr"`
	URes string			`xml:"uRes,attr"`
}
