// package dvevaluation manages expressions, functions using agrammar
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
	"strings"
)

func AnyToString(v interface{}) string {
	return AnyToStringWithOptions(v, ConversionOptionJSLike)
}

func AnyToStringWithOptions(v interface{}, options int) string {
	f, ok := ConvertSimpleTypeToString(v)
	if ok {
		return f
	}
	switch v.(type) {
	case string:
		f = v.(string)
	case *DvObject:
		f = v.(*DvObject).ToString()
	case *DvFunction:
		f = v.(*DvFunction).ToString()
	case nil:
		f = nullValueVersion[options]
	default:
		f = ConvertAnyTypeToJsonString(v)
	}
	return f
}

func AnyGetType(v interface{}) int {
	f := 0
	switch v.(type) {
	case string:
		f = dvgrammar.TYPE_STRING
	case *DvObject:
		f = v.(*DvObject).GetObjectType()
	case *DvFunction:
		f = dvgrammar.TYPE_FUNCTION
	case int:
		f = dvgrammar.TYPE_NUMBER_INT
	case int64:
		f = dvgrammar.TYPE_NUMBER_INT
	case bool:
		f = dvgrammar.TYPE_BOOLEAN
	case float64:
		f = dvgrammar.TYPE_NUMBER
		if math.IsNaN(v.(float64)) {
			f = dvgrammar.TYPE_NAN
		}
	case nil:
		f = dvgrammar.TYPE_UNDEFINED
	}
	return f
}

func AnyWithTypeToString(kind int, v interface{}) (string, bool) {
	switch kind {
	case dvgrammar.TYPE_STRING, dvgrammar.TYPE_CHAR, dvgrammar.TYPE_NUMBER,
		dvgrammar.TYPE_NUMBER_INT, dvgrammar.TYPE_BOOLEAN:
		return AnyToString(v), true
	case dvgrammar.TYPE_NULL:
		return "null", true
	case dvgrammar.TYPE_UNDEFINED:
		return "undefined", true
	case dvgrammar.TYPE_NAN:
		return "NaN", true
	}
	return "", false
}

func AnyToNumber(v interface{}) float64 {
	var f float64
	switch v.(type) {
	case string:
		f = StringToNumber(v.(string))
	case *DvObject:
		f = v.(*DvObject).ToNumber()
	case int:
		f = float64(v.(int))
	case int64:
		f = float64(v.(int64))
	case bool:
		if v.(bool) {
			f = 0
		} else {
			f = 1
		}
	case float64:
		f = v.(float64)
	default:
		f = math.NaN()
	}
	return f
}

func AnyWithTypeToNumber(kind int, v interface{}) float64 {
	switch kind {
	case dvgrammar.TYPE_STRING,
		dvgrammar.TYPE_CHAR,
		dvgrammar.TYPE_NUMBER,
		dvgrammar.TYPE_NUMBER_INT,
		dvgrammar.TYPE_BOOLEAN:
		return AnyToNumber(v)
	case dvgrammar.TYPE_NULL:
		return 0
	}
	return math.NaN()
}

func AnyToBoolean(v interface{}) bool {
	var f bool = false
	switch v.(type) {
	case string:
		f = v.(string) != ""
	case *DvObject:
		f = v.(*DvObject).ToBoolean()
	case int:
		f = v.(int) != 0
	case int64:
		f = v.(int64) != 0
	case bool:
		f = v.(bool)
	case float64:
		e := v.(float64)
		f = e != 0 && !math.IsNaN(e)
	case nil:
		f = false
	}
	return f
}

func AnyWithTypeToBoolean(kind int, v interface{}, defValue bool) bool {
	switch kind {
	case dvgrammar.TYPE_STRING, dvgrammar.TYPE_CHAR, dvgrammar.TYPE_NUMBER,
		dvgrammar.TYPE_NUMBER_INT, dvgrammar.TYPE_BOOLEAN:
		return AnyToBoolean(v)
	case dvgrammar.TYPE_NULL:
		return false
	}
	return defValue
}

func StringToAny(v string) interface{} {
	v = strings.TrimSpace(v)
	if t, ok := buildinTypes[v]; ok {
		return t
	}
	c := v[0]
	if c >= '0' && c <= '9' || c == '+' || c == '-' {
		n := len(v)
		kind, pos, vint, vfloat := GetNumberKindFromString([]byte(v), n)
		if kind >= 0 && pos == n {
			switch kind {
			case dvgrammar.TYPE_NUMBER:
				return vfloat
			case dvgrammar.TYPE_NUMBER_INT:
				return vint
			}
		}
	} else if len(v) >= 2 && (c == '"' || c == '`') && v[len(v)-1] == c {
		return v[1 : len(v)-1]
	}
	return v
}

func NumberToInt(v float64) (f int64, ok bool) {
	d := math.Trunc(v)
	if math.IsNaN(d) {
		ok = false
	} else {
		f = int64(d)
		ok = true
	}
	return
}

func AnyToNumberInt(v interface{}) (f int64, ok bool) {
	ok = true
	switch v.(type) {
	case string:
		s := v.(string)
		n := len(s)
		kind, pos, vint, vfloat := GetNumberKindFromString([]byte(s), n)
		ok = kind >= 0 && pos == n
		if ok {
			switch kind {
			case dvgrammar.TYPE_NUMBER:
				f, ok = NumberToInt(vfloat)
			case dvgrammar.TYPE_NUMBER_INT:
				f = vint
			}
		}
	case *DvObject:
		f, ok = v.(*DvObject).ToNumberInt()
	case int:
		f = int64(v.(int))
	case int64:
		f = v.(int64)
	case bool:
		if v.(bool) {
			f = 0
		} else {
			f = 1
		}
	case float64:
		f, ok = NumberToInt(v.(float64))
	default:
		ok = false
	}
	return
}

func AnyWithTypeToNumberInt(kind int, v interface{}) (int64, bool) {
	switch kind {
	case dvgrammar.TYPE_STRING,
		dvgrammar.TYPE_CHAR,
		dvgrammar.TYPE_NUMBER,
		dvgrammar.TYPE_NUMBER_INT,
		dvgrammar.TYPE_BOOLEAN:
		return AnyToNumberInt(v)
	case dvgrammar.TYPE_NULL:
		return 0, true
	}
	return 0, false
}

func ConvertInterfaceListToStringList(list []interface{}, options int) []string {
	n := len(list)
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = AnyToStringWithOptions(list[i], options)
	}
	return r
}

func ConvertInterfaceListsToStringLists(list [][]interface{}, options int) [][]string {
	n := len(list)
	r := make([][]string, n)
	for i := 0; i < n; i++ {
		r[i] = ConvertInterfaceListToStringList(list[i], options)
	}
	return r
}

func ConvertInterfaceListsMapToStringListsMap(listMap map[string][][]interface{}, options int) map[string][][]string {
	r := make(map[string][][]string, len(listMap))
	for k, v := range listMap {
		r[k] = ConvertInterfaceListsToStringLists(v, options)
	}
	return r
}
