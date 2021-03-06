// Code generated by "enumer -type=ConceptKind"; DO NOT EDIT

package ucum

import (
	"fmt"
)

const _ConceptKindName = "PREFIXBASEUNITUNIT"

var _ConceptKindIndex = [...]uint8{0, 6, 14, 18}

func (i ConceptKind) String() string {
	if i < 0 || i >= ConceptKind(len(_ConceptKindIndex)-1) {
		return fmt.Sprintf("ConceptKind(%d)", i)
	}
	return _ConceptKindName[_ConceptKindIndex[i]:_ConceptKindIndex[i+1]]
}

var _ConceptKindValues = []ConceptKind{0, 1, 2}

var _ConceptKindNameToValueMap = map[string]ConceptKind{
	_ConceptKindName[0:6]:   0,
	_ConceptKindName[6:14]:  1,
	_ConceptKindName[14:18]: 2,
}

// ConceptKindString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ConceptKindString(s string) (ConceptKind, error) {
	if val, ok := _ConceptKindNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ConceptKind values", s)
}

// ConceptKindValues returns all values of the enum
func ConceptKindValues() []ConceptKind {
	return _ConceptKindValues
}

// IsAConceptKind returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ConceptKind) IsAConceptKind() bool {
	for _, v := range _ConceptKindValues {
		if i == v {
			return true
		}
	}
	return false
}
