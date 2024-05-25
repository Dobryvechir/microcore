/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"strconv"
	"strings"
)

const (
	FIELD_UNDEFINED = iota
	FIELD_NULL
	FIELD_NUMBER
	FIELD_BOOLEAN
	FIELD_STRING
	FIELD_OBJECT
	FIELD_ARRAY
	FIELD_FUNCTION
)

var typeOfSpecific map[int]string = map[int]string{
	FIELD_UNDEFINED: "undefined",
	FIELD_NULL:      "null",
	FIELD_NUMBER:    "number",
	FIELD_BOOLEAN:   "boolean",
	FIELD_STRING:    "string",
	FIELD_OBJECT:    "object",
	FIELD_ARRAY:     "array",
	FIELD_FUNCTION:  "function",
}

var typeOf map[int]string = map[int]string{
	FIELD_UNDEFINED: "undefined",
	FIELD_NULL:      "object",
	FIELD_NUMBER:    "number",
	FIELD_BOOLEAN:   "boolean",
	FIELD_STRING:    "string",
	FIELD_OBJECT:    "object",
	FIELD_ARRAY:     "object",
	FIELD_FUNCTION:  "function",
}

type QuickSearchInfo struct {
	Looker map[string]*DvVariable
	Key    string
}

type DvVariable struct {
	Name        []byte
	Value       []byte
	Kind        int
	Fields      []*DvVariable
	Extra       interface{}
	Prototype   *DvVariable
	QuickSearch *QuickSearchInfo
}

type DvVariable_DumpInfo struct {
	used map[*DvVariable]bool
	buf  []byte
}

func DvVariableGetNewObject() *DvVariable {
	variable := &DvVariable{
		Kind:      FIELD_OBJECT,
		Fields:    make([]*DvVariable, 0, 7),
		Prototype: ObjectMaster,
	}
	return variable
}

func DeleteVariable(parent *DvVariable, keys []string, silent bool) error {
	l := len(keys)
	if l == 0 {
		return nil
	}
	key := keys[0]
	if parent == nil {
		return errors.New("Cannot delete " + key + " from undefined [keys:" + strings.Join(keys, ", ") + "]")
	}
	l--
	for i := 0; i < l; i++ {
		key = keys[i]
		if parent.Fields == nil {
			if silent {
				return nil
			}
			return errors.New("Cannot delete " + key + " from undefined [keys:" + strings.Join(keys, ", ") + "]")
		}
		if child, ok := parent.FindChildByKey(key); ok {
			parent = child
		} else {
			if silent {
				return nil
			}
			return errors.New("Cannot delete from undefined " + key + " [keys:" + strings.Join(keys, ", ") + "]")
		}
	}
	key = keys[l]
	if parent.Fields != nil {
		ind := parent.FindChildIndexByKey(key)
		if ind >= 0 {
			parent.DeleteChildByIndex(ind)
		}
	}
	return nil
}

func ConvertStringArrayToDvVariableArray(data []string) (res []*DvVariable) {
	n := len(data)
	res = make([]*DvVariable, n)
	for i := 0; i < n; i++ {
		res[i] = DvVariableFromString(nil, data[i])
	}
	return
}

func ConvertStringArrayToInterfaceArray(data []string) (res []interface{}) {
	n := len(data)
	res = make([]interface{}, n)
	for i := 0; i < n; i++ {
		res[i] = data[i]
	}
	return
}

func GetBooleanValue(isTrue bool) []byte {
	if isTrue {
		return bytesTrue
	}
	return bytesFalse
}

func ConvertAnyToDvVariable(data interface{}) *DvVariable {
	switch data.(type) {
	case *DvVariable:
		return data.(*DvVariable)
	case string:
		return &DvVariable{Kind: FIELD_STRING, Value: []byte(data.(string))}
	case nil:
		return &DvVariable{Kind: FIELD_NULL}
	}
	buf := make([]byte, 0, 100)
	b, ok, kind := ConvertSimpleTypeToBuf(buf, data)
	if ok {
		return &DvVariable{Kind: kind, Value: b}
	}
	s := AnyToString(data)
	return &DvVariable{Kind: FIELD_STRING, Value: []byte(s)}
}

func ConvertStringMapToDvVariableMap(data map[string]string) (res map[string]*DvVariable) {
	res = make(map[string]*DvVariable)
	if data != nil {
		for k, v := range data {
			res[k] = DvVariableFromString(nil, v)
		}
	}
	return
}

func DvVariableFromStringArray(parent *DvVariable, data []string) *DvVariable {
	return DvVariableFromArray(parent, ConvertStringArrayToDvVariableArray(data))
}

func DvVariableFromStringMap(parent *DvVariable, data map[string]string) *DvVariable {
	return DvVariableFromMap(parent, ConvertStringMapToDvVariableMap(data))
}

func GetEscapedString(s string) string {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] == '\\' {
			s = s[0:i] + s[i+1:n]
		}
	}
	return s
}

func ConvertVariableNameToKeys(data string) (res []string, err error) {
	data = strings.TrimSpace(data)
	l := len(data)
	res = make([]string, 0, 1)
	p := 0
	for i := 0; i < l; i++ {
		c := data[i]
		if c == '.' {
			s := data[p:i]
			res = append(res, s)
			p = i + 1
		} else if c == '[' {
			i++
			p = i
			f := byte(']')
			r := byte(0)
			if i < l {
				r = byte(data[i])
			}
			if r == '\'' || r == '`' || r == '"' {
				i++
				p = i
				f = r
			}
			escape := false
			for ; i < l && data[i] != f; i++ {
				if data[i] == '\\' {
					i++
					escape = true
				}
			}
			if i > l {
				return nil, errors.New("Error in escaped string: " + data)
			}
			s := data[p:i]
			p = i + 1
			if escape {
				s = GetEscapedString(s)
			}
			res = append(res, s)
			if f != ']' {
				for i++; i < l && data[i] != ']'; i++ {
					if data[i] > 32 {
						return nil, errors.New("Error inside []: " + data)
					}
				}

			}
		}
	}
	if p < l {
		s := data[p:l]
		res = append(res, s)
	}
	return
}

func getEscapedByteArray(buf []byte) []byte {
	n := len(buf)
	for i := 0; i < n; i++ {
		if buf[i] == '\\' || buf[i] == '"' {
			rest := buf[i:]
			buf = append(buf[:i:i], byte('\\'))
			buf = append(buf, rest...)
			n++
			i++
		}
	}
	return buf
}

func (variable *DvVariable) dumpDvVariable(info *DvVariable_DumpInfo) *DvVariable_DumpInfo {
	if info == nil {
		info = &DvVariable_DumpInfo{used: make(map[*DvVariable]bool), buf: make([]byte, 0, 16384)}
	}
	if variable == nil {
		info.buf = append(info.buf, 'n', 'u', 'l', 'l')
		return info
	}
	if variable.Kind == FIELD_ARRAY || variable.Kind == FIELD_OBJECT {
		info.used[variable] = true
		openQuote := byte('[')
		closeQuote := byte(']')
		if variable.Kind == FIELD_OBJECT {
			openQuote = byte('{')
			closeQuote = byte('}')
		}
		info.buf = append(info.buf, openQuote)
		comma := false
		for _, v := range variable.Fields {
			if comma {
				info.buf = append(info.buf, ',')
			}
			comma = true
			info.buf = append(append(append(info.buf, byte('"')), v.Name...), '"', ':')
			if _, ok := info.used[v]; ok {
				info.buf = append(info.buf, '*')
			} else {
				v.dumpDvVariable(info)
			}
		}
		info.buf = append(info.buf, closeQuote)
	} else {
		switch variable.Kind {
		case FIELD_UNDEFINED, FIELD_NULL:
			info.buf = append(info.buf, []byte(typeOfSpecific[variable.Kind])...)
		case FIELD_STRING:
			info.buf = append(append(append(info.buf, '"'), getEscapedByteArray([]byte(variable.Value))...), '"')
		case FIELD_FUNCTION:
			info.buf = append(info.buf, []byte("function "+string(variable.Value))...)
		default:
			info.buf = append(info.buf, []byte(variable.Value)...)
		}
	}
	return info
}

func (variable *DvVariable) JsonStringify() []byte {
	info := variable.dumpDvVariable(nil)
	return info.buf
}

func (variable *DvVariable) JsonStringifyNonEmpty() []byte {
	if variable == nil || variable.Kind == FIELD_UNDEFINED || variable.Kind == FIELD_NULL {
		return []byte{}
	}
	return variable.JsonStringify()
}

func (variable *DvVariable) GetStringValueAsBytes() []byte {
	if variable == nil || variable.Kind == FIELD_UNDEFINED || variable.Kind == FIELD_NULL {
		return []byte{}
	}
	return []byte(variable.GetStringValue())
}

func ConvertMapDvVariableToList(varMap map[string]*DvVariable) []*DvVariable {
	n := len(varMap)
	res := make([]*DvVariable, n)
	i := 0
	for k, v := range varMap {
		res[i] = v.MakeCopyWithNewKey(k)
		i++
	}
	return res
}

func (variable *DvVariable) IsEmpty() bool {
	if variable == nil {
		return true
	}
	switch variable.Kind {
	case FIELD_UNDEFINED:
	case FIELD_NULL:
		return true
	case FIELD_OBJECT:
	case FIELD_ARRAY:
		return len(variable.Fields) == 0
	default:
		return len(variable.Value) == 0
	}
	return false
}

func (variable *DvVariable) CleanValue() {
	if variable == nil {
		return
	}
	switch variable.Kind {
	case FIELD_UNDEFINED:
	case FIELD_NULL:
		break
	case FIELD_OBJECT:
	case FIELD_ARRAY:
		variable.Fields = nil
	default:
		variable.Value = make([]byte, 0, 8)
	}
}

func (variable *DvVariable) CleanFields(keys []string) {
	if variable == nil || variable.Kind != FIELD_OBJECT || len(variable.Fields) == 0 {
		return
	}
	keyMap := collectKeysToMap(keys)
	n := len(variable.Fields)
	for i := 0; i < n; i++ {
		f := variable.Fields[i]
		if f == nil {
			continue
		}
		k := string(f.Name)
		if keyMap[k] == 1 {
			f.CleanValue()
		}
	}
}

func (variable *DvVariable) CopyFieldsFromOther(keys []string, other *DvVariable) {
	if variable == nil || variable.Kind != FIELD_OBJECT || other == nil || len(other.Fields) == 0 {
		return
	}
	keyMap := collectKeysToMap(keys)
	n := len(other.Fields)
	for i := 0; i < n; i++ {
		f := variable.Fields[i]
		if f == nil {
			continue
		}
		k := string(f.Name)
		if keyMap[k] == 1 {
			variable.SetField(k, f)
		}
	}
}

func (variable *DvVariable) CopyFieldsToMap(prefix string, data map[string]interface{}) {
	if variable == nil || len(variable.Fields) == 0 {
		return
	}
	f := variable.Fields
	n := len(f)
	if variable.Kind == FIELD_ARRAY {
		for i := 0; i < n; i++ {
			data[prefix+strconv.Itoa(i)] = f[i]
		}
	} else {
		for i := 0; i < n; i++ {
			v := f[i]
			if v == nil || len(v.Name) == 0 {
				continue
			}
			data[prefix+string(v.Name)] = v
		}
	}
}
