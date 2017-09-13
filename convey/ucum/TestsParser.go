package ucum


type XMLUcumTests struct {
	validations  XMLValidation  	`xml:"validation"`
	displayNameGeneration	XMLDisplayNameGeneration	`xml:"displayNameGeneration"`
	conversion	XMLConversion	`xml:"conversion"`
	multiplication	XMLMultiplication	`xml:"multiplication"`
}

type XMLValidation struct{
	case_ XMLValidationCase					`xml:"case"`
}

type XMLValidationCase struct {
	id string						`xml:"id, attr"`
	unit string						`xml:"unit, attr"`
	valid string					`xml:"valid, attr"`
	reason string					`xml:"reason, attr"`
}

type XMLDisplayNameGeneration struct {
	id string						`xml:"id, attr"`
	unit string						`xml:"unit, attr"`
	display string					`xml:"valid, attr"`
}

type XMLConversion struct {
	id string			`xml:"id, attr"`
	value string		`xml:"value, attr"`
	srcUnit string		`xml:"srcUnit, attr"`
	dstUnit string		`xml:"dstUnit, attr"`
	outcome string		`xml:"outcome, attr"`
}

type XMLMultiplication struct {
	id string			`xml:"id, attr"`
	v1 string			`xml:"v1, attr"`
	u1 string			`xml:"u1, attr"`
	v2 string			`xml:"v2, attr"`
	u2 string			`xml:"u2, attr"`
	vRes string			`xml:"vRes, attr"`
	uRes string			`xml:"uRes, attr"`
}
