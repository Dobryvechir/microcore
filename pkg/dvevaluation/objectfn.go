/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

const LENGTH_PROPERTY = "length"

func (variable *DvVariable) ObjectInPrototypedChain(key string) *DvVariable {
	if variable == nil {
		return nil
	}
	if variable.Fields != nil {
		if v, ok := variable.FindChildByKey(key); ok {
			return v
		}
	}
	if variable.Prototype != nil {
		return variable.Prototype.ObjectInPrototypedChain(key)
	}
	return nil
}

func (variable *DvVariable) ObjectGet(key string) (*DvVariable, error) {
	if variable == nil {
		return nil, errors.New("Cannot get " + key + " from non-existing variable")
	}
	if variable.Kind == FIELD_UNDEFINED {
		return nil, errors.New("Cannot get " + key + " from undefined")
	}
	if variable.Kind == FIELD_NULL {
		return nil, errors.New("Cannot get " + key + " from null")
	}
	if v, ok := variable.FindChildByKey(key); ok {
		return v, nil
	}
	if variable.Prototype != nil {
		o := variable.Prototype.ObjectInPrototypedChain(key)
		if o != nil {
			return o, nil
		}
	}
	return &DvVariable{}, nil
}

func (variable *DvVariable) ObjectGetByKeys(keys []string) (*DvVariable, error) {
	n := len(keys)
	for i := 0; i < n; i++ {
		v, err := variable.ObjectGet(keys[i])
		if err != nil {
			return nil, err
		}
		variable = v
	}
	return variable, nil
}

func (variable *DvVariable) ObjectGetByVariableDefinition(variableDefinition string) (*DvVariable, error) {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return nil, err
	}
	return variable.ObjectGetByKeys(keys)
}

func (variable *DvVariable) GetVariableArray(force bool) ([]*DvVariable, error) {
	if variable == nil || (variable.Kind != FIELD_ARRAY && variable.Kind != FIELD_OBJECT) {
		if force {
			return []*DvVariable{variable}, nil
		}
		return nil, errors.New("Array is expected")
	}
	if variable.Fields != nil {
		return variable.Fields, nil
	}
	res := make([]*DvVariable, 0, 7)
	return res, nil
}

func (variable *DvVariable) GetIntValue(force bool) (int, error) {
	if variable == nil {
		if force {
			return 0, nil
		} else {
			return 0, errors.New("Cannot convert no variable to integer")
		}
	}
	if variable.Kind > FIELD_OBJECT {
		if force {
			return 1, nil
		} else {
			return 0, errors.New("Cannot convert object to integer")
		}
	}
	if variable.Kind == FIELD_UNDEFINED || variable.Kind == FIELD_NULL || len(variable.Value) == 0 {
		return 0, nil
	}
	if variable.Kind == FIELD_BOOLEAN {
		if len(variable.Value) == 0 || variable.Value[0] == 'f' || variable.Value[0] == 'F' {
			return 0, nil
		}
		return 1, nil
	}
	return strconv.Atoi(string(variable.Value))
}

func (variable *DvVariable) GetStringValueJS() string {
	if variable == nil || variable.Kind == FIELD_UNDEFINED {
		return typeOfSpecific[FIELD_UNDEFINED]
	}
	if variable.Kind == FIELD_NULL || variable.Kind == FIELD_FUNCTION {
		return typeOfSpecific[variable.Kind]
	}
	if variable.Kind >= FIELD_OBJECT {
		return "[object Object]"
	}
	return string(variable.Value)
}

func (variable *DvVariable) GetStringArrayValue() []string {
	if variable == nil || variable.Kind == FIELD_UNDEFINED {
		return []string{typeOfSpecific[FIELD_UNDEFINED]}
	}
	if variable.Kind == FIELD_NULL || variable.Kind == FIELD_FUNCTION {
		return []string{typeOfSpecific[variable.Kind]}
	}
	if variable.Kind >= FIELD_OBJECT {
		res := make([]string, 0, 8)
		for _, v := range variable.Fields {
			res = append(res, v.GetStringValue())
		}
		return res
	}
	return []string{string(variable.Value)}
}

func (variable *DvVariable) GetStringMap() (res map[string]string) {
	res = make(map[string]string)
	if variable == nil || variable.Kind == FIELD_UNDEFINED || variable.Fields == nil {
		return
	}
	for k, v := range variable.Fields {
		key := string(v.Name)
		if key == "" {
			key = strconv.Itoa(k)
		}
		res[key] = v.GetStringValue()
	}
	return
}

func (variable *DvVariable) GetStringInterfaceMap() (res map[string]interface{}) {
	res = make(map[string]interface{})
	if variable == nil || variable.Kind == FIELD_UNDEFINED || variable.Fields == nil {
		return
	}
	for k, v := range variable.Fields {
		key := string(v.Name)
		if key == "" {
			key = strconv.Itoa(k)
		}
		res[key] = v
	}
	return
}

func (variable *DvVariable) GetStringArrayMap() (res map[string][]string) {
	res = make(map[string][]string)
	if variable == nil || variable.Kind == FIELD_UNDEFINED || variable.Fields == nil {
		return
	}
	for k, v := range variable.Fields {
		key := string(v.Name)
		if key == "" {
			key = strconv.Itoa(k)
		}
		res[key] = v.GetStringArrayValue()
	}
	return
}

func (variable *DvVariable) SetSimpleValue(value string, kind int) error {
	if variable == nil {
		return errors.New("Cannot assign to null variable")
	}
	variable.Value = []byte(value)
	variable.Kind = kind
	return nil
}

func (variable *DvVariable) SetField(key string, value *DvVariable) int {
	if variable == nil || variable.Kind != FIELD_OBJECT {
		return -2
	}
	newValue := value.CloneExceptKey(value, true)
	newValue.Name = []byte(key)
	n := variable.FindIndex(key)
	if n < 0 {
		variable.Fields = append(variable.Fields, newValue)
	} else {
		variable.Fields[n] = newValue
	}
	return n
}

func ValidateNumber(data string) error {
	l := len(data)
	i := 0
	for i < l && data[i] <= 32 {
		i++
	}
	if i < l && (data[i] == '-' || data[i] == '+') {
		i++
	}
	if i == l {
		return errors.New("Empty string is not a number")
	}
	b := data[i]
	comma := false
	if b == '.' {
		comma = true
		i++
		if i < l {
			b = data[i]
		} else {
			return errors.New("Only . is not a number")
		}
	}
	if b < '0' || b > '9' {
		return errors.New(data + " is not a number")
	}
	for i < l && (b >= '0' && b <= '9' || b == '.') {
		if b == '.' {
			if comma {
				return errors.New("Only one point is allowed in a number: " + data)
			} else {
				comma = true
			}
		}
		i++
		if i < l {
			b = data[i]
		}
	}
	if i < l-1 && (data[i] == 'e' || data[i] == 'E') && (data[i+1] == '+' || data[i+1] == '-' || data[i+1] >= '0' && data[i+1] <= '9') {
		for i < l && data[i] >= '0' && data[i] <= '9' {
			i++
		}
	}
	for i < l && data[i] <= 32 {
		i++
	}
	if i != l {
		return errors.New("Wrong characters at the end of the number [" + data[i:] + "]")
	}
	return nil
}

func QuickNumberEvaluation(parent *DvVariable, data string) (res *DvVariable, err error) {
	if parent == nil {
		res = &DvVariable{Kind: FIELD_STRING}
	} else {
		res = parent
		res.Kind = FIELD_NUMBER
	}
	res.Value = []byte(data)
	err = ValidateNumber(data)
	return
}

func QuickStringEvaluation(parent *DvVariable, data string) (res *DvVariable, err error) {
	if parent == nil {
		res = &DvVariable{Kind: FIELD_STRING}
	} else {
		res = parent
		res.Kind = FIELD_STRING
	}
	l := len(data)
	for l > 0 && data[l-1] <= 32 {
		l--
	}
	i := 0
	for i < l && data[i] <= 32 {
		i++
	}
	if l-i < 2 {
		err = errors.New("Wrong string: " + data)
		return
	}
	c := data[i]
	if (c != '\'' && c != '"' && c != '`') || data[l-1] != c {
		err = errors.New("Wrong string: " + data)
	}
	res.Value = []byte(dvgrammar.GetEscapedString([]byte(data[i+1 : l-1])))
	return
}

func QuickVariableEvaluation(parent *DvVariable, data string) (res *DvVariable, err error) {
	res = &DvVariable{}
	l := len(data)
	for l > 0 && data[l-1] <= 32 {
		l--
	}
	i := 0
	for i < l && data[i] <= 32 {
		i++
	}
	if i == l {
		return
	}
	data = data[i:l]
	l -= i
	b := data[0]
	switch b {
	case '+', '-':
		return QuickNumberEvaluation(res, data)
	case '"', '`', '\'':
		return QuickStringEvaluation(res, data)
	}
	res.Value = []byte(data)
	switch data {
	case "undefined":
		return
	case "null":
		res.Kind = FIELD_NULL
		return
	case "true", "false":
		res.Kind = FIELD_BOOLEAN
		return
	}
	if b >= '0' && b <= '9' || b == '.' && l > 1 && data[1] >= '0' && data[i] <= '9' {
		return QuickNumberEvaluation(res, data)
	}
	if parent == nil {
		err = errors.New("Cannot evaluate constant " + data)
		return
	}
	_, res, _, err, _ = GetVariableByDefinition(parent, data)
	return
}

func QuickVariableArrayEvaluation(parent *DvVariable, data []string) (res []interface{}, err error) {
	res = make([]interface{}, len(data))
	for i, v := range data {
		res[i], err = QuickVariableEvaluation(parent, v)
		if err != nil {
			err = dvgrammar.EnrichErrorStr(err, "At variable array evaluation "+strings.Join(data, ", "))
			return
		}
	}
	return
}
