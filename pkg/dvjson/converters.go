/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"errors"
	"strconv"
)

const MaxInt = int64(^uint(0) >> 1)

func ConvertByteArrayToIntOrDouble(data []byte) (interface{}, bool) {
	n := len(data)
	i := 0
	for ; i < n && data[i] <= ' '; i++ {
	}
	positive := true
	if data[i] == '-' {
		positive = false
		i++
	}
	if data[i] == '+' {
		i++
	}
	v := int64(0)
	for ; i < n; i++ {
		c := data[i]
		if c >= '0' && c <= '9' {
			v = v*10 + int64(c) - 48
		} else if c == '.' || c == 'e' || c == 'E' {
			break
		} else {
			return 0, false
		}
	}
	if i == n {
		if !positive {
			v = -v
		}
		if v <= MaxInt && v >= -MaxInt {
			return int(v), true
		}
		return float64(v), true
	}
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, false
	}
	vf := int(f)
	if float64(vf) == f {
		return vf, true
	}
	return f, true
}

func ConvertSimpleKindAndValueToInterface(kind int, data []byte) (interface{}, bool) {
	switch kind {
	case FIELD_EMPTY:
		return nil, true
	case FIELD_NULL:
		return nil, true
	case FIELD_NUMBER:
		return ConvertByteArrayToIntOrDouble(data)
	case FIELD_BOOLEAN:
		return len(data) != 0 && (data[0] == 't' || data[0] == 'T'), true
	case FIELD_STRING:
		return string(data), true
	}
	return nil, false
}

func (parseInfo *DvCrudParsingInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return ConvertSimpleKindAndValueToInterface(parseInfo.Kind, parseInfo.Value)
}

func (item *DvFieldInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return ConvertSimpleKindAndValueToInterface(item.Kind, item.Value)
}

func (item *DvFieldInfo) ReadSimpleStringMap(data map[string]string) error {
	if item.Kind != FIELD_OBJECT {
		return errors.New(string(item.Name) + " must be an object { }")
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		p := item.Fields[i]
		k := string(p.Name)
		v := string(p.Value)
		data[k] = v
	}
	return nil
}

func (item *DvFieldInfo) ReadSimpleStringList(data []string) ([]string, error) {
	if item.Kind != FIELD_ARRAY {
		return data, errors.New(string(item.Name) + " must be an object { }")
	}
	n := len(item.Fields)
	if data == nil {
		data = make([]string, 0, n)
	}
	for i := 0; i < n; i++ {
		p := item.Fields[i]
		v := string(p.Value)
		data = append(data, v)
	}
	return data, nil
}

func (item *DvFieldInfo) ReadSimpleString() (string, error) {
	if item.Kind == FIELD_OBJECT || item.Kind == FIELD_ARRAY {
		return "[]", errors.New(string(item.Name) + " must be a simple type")
	}
	return string(item.Value), nil
}

func (item *DvFieldInfo) ConvertValueToInterface() (interface{}, bool) {
	switch item.Kind {
	case FIELD_ARRAY:
		fields := item.Fields
		n := len(fields)
		data := make([]interface{}, n)
		var ok bool
		for i := 0; i < n; i++ {
			data[i], ok = fields[i].ConvertValueToInterface()
			if !ok {
				return nil, false
			}
		}
		return data, true
	case FIELD_OBJECT:
		fields := item.Fields
		n := len(fields)
		data := make(map[string]interface{}, n)
		var ok bool
		for i := 0; i < n; i++ {
			data[string(fields[i].Name)], ok = fields[i].ConvertValueToInterface()
			if !ok {
				return nil, false
			}
		}
		return data, true
	}
	return item.ConvertSimpleValueToInterface()
}

func ConvertDvFieldInfoToProperties(item *DvFieldInfo, index int) map[string]interface{} {
	res := map[string]interface{}{"_index": index}
	if item == nil {
		return res
	}
	if len(item.Name) != 0 {
		res["_name"] = string(item.Name)
	}
	if len(item.Value) != 0 {
		res["_value"] = string(item.Value)
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		current := item.Fields[i]
		if current == nil {
			res["_"+strconv.Itoa(i)] = ""
		} else {
			var name string
			if len(current.Name) == 0 {
				name = "_" + strconv.Itoa(i)
			} else {
				name = string(current.Name)
			}
			var value interface{}
			if current.Kind == FIELD_ARRAY || current.Kind == FIELD_OBJECT {
				value = current
			} else {
				value = string(current.Value)
			}
			res[name] = value
		}
	}
	return res
}

func (item *DvFieldInfo) GetStringValue() string {
	if item == nil || item.Kind == FIELD_EMPTY {
		return ""
	}
	switch item.Kind {
	case FIELD_OBJECT:
		res := "{"
		subfields := item.Fields
		n := len(subfields)
		for i := 0; i < n; i++ {
			if i != 0 {
				res += ","
			}
			res += "\"" + string(subfields[i].Name) + "\":"
			data := subfields[i].GetStringValue()
			res += data
		}
		return res + "}"
	case FIELD_ARRAY:
		res := "["
		subfields := item.Fields
		n := len(subfields)
		for i := 0; i < n; i++ {
			if i != 0 {
				res += ","
			}
			data := subfields[i].GetStringValue()
			res += data
		}
		return res + "]"
	case FIELD_STRING:
		return QuoteEscapedJsonBytesToString(item.Value)
	case FIELD_NULL:
		return "null"
	}
	return string(item.Value)
}
func DvFieldArrayToBytes(val []*DvFieldInfo) []byte {
	buf := make([]byte, 1, 102400)
	buf[0] = '['
	n := len(val)
	for i := 0; i < n; i++ {
		if i != 0 {
			buf = append(buf, ',')
		}
		b := val[i].PrintToJson(2)
		buf = append(buf, b...)
	}
	buf = append(buf, ']')
	return buf
}

func DvFieldInfoToStringConverter(v interface{}) (string, bool) {
	switch v.(type) {
	case *DvFieldInfo:
		return string(v.(*DvFieldInfo).PrintToJson(2)), true
	case []*DvFieldInfo:
		return string(DvFieldArrayToBytes(v.([]*DvFieldInfo))), true
	}
	return "", false
}
