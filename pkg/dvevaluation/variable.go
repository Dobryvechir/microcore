/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"strings"
)

const (
	JS_TYPE_UNDEFINED = iota
	JS_TYPE_NULL      = iota
	JS_TYPE_NUMBER    = iota
	JS_TYPE_BOOLEAN   = iota
	JS_TYPE_STRING    = iota
	JS_TYPE_OBJECT    = iota
	JS_TYPE_ARRAY     = iota
	JS_TYPE_FUNCTION  = iota
)

var typeOfSpecific map[int]string = map[int]string{
	0: "undefined",
	1: "null",
	2: "number",
	3: "boolean",
	4: "string",
	5: "object",
	6: "array",
	7: "function",
}

var typeOf map[int]string = map[int]string{
	0: "undefined",
	1: "object",
	2: "number",
	3: "boolean",
	4: "string",
	5: "object",
	6: "object",
	7: "function",
}

type DvvFunction func(*DvContext, *DvVariable, []*DvVariable) (*DvVariable, error)

type DvVariable struct {
	Refs      map[string]*DvVariable
	Value     string
	Tp        int
	Fn        DvvFunction
	Prototype *DvVariable
}

type DvVariable_DumpInfo struct {
	used map[*DvVariable]bool
	buf  []byte
}

func DvVariableGetNewObject() *DvVariable {
	variable := &DvVariable{Tp: JS_TYPE_OBJECT, Refs: make(map[string]*DvVariable)}
	return variable
}

func DeleteVariable(parent *DvVariable, keys []string) error {
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
		if parent.Refs == nil {
			return errors.New("Cannot delete " + key + " from undefined [keys:" + strings.Join(keys, ", ") + "]")
		}
		if child, ok := parent.Refs[key]; ok {
			parent = child
		} else {
			return errors.New("Cannot delete from undefined " + key + " [keys:" + strings.Join(keys, ", ") + "]")
		}
	}
	key = keys[l]
	if parent.Refs != nil {
		if _, ok := parent.Refs[key]; ok {
			delete(parent.Refs, key)
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
	return DvVariableFromMap(parent, ConvertStringMapToDvVariableMap(data), true)
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
	if variable.Tp == JS_TYPE_ARRAY || variable.Tp == JS_TYPE_OBJECT {
		info.used[variable] = true
		openQuote := byte('[')
		closeQuote := byte(']')
		if variable.Tp == JS_TYPE_OBJECT {
			openQuote = byte('{')
			closeQuote = byte('}')
		}
		info.buf = append(info.buf, openQuote)
		comma := false
		var vlength *DvVariable = nil
		for k, v := range variable.Refs {
			if k == "length" {
				vlength = v
			} else {
				if comma {
					info.buf = append(info.buf, ',')
				}
				comma = true
				info.buf = append(append(append(info.buf, byte('"')), []byte(k)...), '"', ':')
				if _, ok := info.used[v]; ok {
					info.buf = append(info.buf, '*')
				} else {
					v.dumpDvVariable(info)
				}
			}
		}
		if vlength != nil {
			if comma {
				info.buf = append(info.buf, ',')
			}
			info.buf = append(info.buf, []byte(`"length":`)...)
			if _, ok := info.used[vlength]; ok {
				info.buf = append(info.buf, '*')
			} else {
				vlength.dumpDvVariable(info)
			}
		}
		info.buf = append(info.buf, closeQuote)
	} else {
		switch variable.Tp {
		case JS_TYPE_UNDEFINED, JS_TYPE_NULL:
			info.buf = append(info.buf, []byte(typeOfSpecific[variable.Tp])...)
		case JS_TYPE_STRING:
			info.buf = append(append(append(info.buf, '"'), getEscapedByteArray([]byte(variable.Value))...), '"')
		case JS_TYPE_FUNCTION:
			info.buf = append(info.buf, []byte("function "+variable.Value)...)
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
	if variable == nil || variable.Tp == JS_TYPE_UNDEFINED || variable.Tp == JS_TYPE_NULL {
		return []byte{}
	}
	return variable.JsonStringify()
}

func (variable *DvVariable) GetStringValueAsBytes() []byte {
	if variable == nil || variable.Tp == JS_TYPE_UNDEFINED || variable.Tp == JS_TYPE_NULL {
		return []byte{}
	}
	return []byte(variable.GetStringValue())
}
