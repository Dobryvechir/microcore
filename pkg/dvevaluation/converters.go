/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"log"
	"math"
	"strconv"
	"strings"
)

type MethodToStringConverter func(v interface{}) (string, bool)

var poolStringConverters = make([]MethodToStringConverter, 0, 7)

func RegisterToStringConverter(converter MethodToStringConverter) {
	poolStringConverters = append(poolStringConverters, converter)
}

func AnyToString(v interface{}) string {
	return AnyToStringWithOptions(v, ConversionOptionJsonLike)
}

func AnyToByteArray(v interface{}) []byte {
	switch v.(type) {
	case []byte:
		return v.([]byte)
	case string:
		return []byte(v.(string))
	}
	return []byte(AnyToStringWithOptions(v, ConversionOptionJsonLike))
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
	case []byte:
		f = string(v.([]byte))
	case nil:
		f = nullValueVersion[options]
	case *dvgrammar.ExpressionValue:
		b := v.(*dvgrammar.ExpressionValue)
		if b == nil {
			return ""
		}
		switch b.DataType {
		case dvgrammar.TYPE_NULL, dvgrammar.TYPE_NAN:
			return ""
		default:
			return AnyToString(b.Value)
		}
	default:
		n := len(poolStringConverters)
		done := false
		for i := 0; i < n; i++ {
			f, done = poolStringConverters[i](v)
			if done {
				break
			}
		}
		if !done {
			f = ConvertAnyTypeToJsonString(v)
		}
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
	case *DvVariable:
		f = dvgrammar.TYPE_OBJECT
	case nil:
		f = dvgrammar.TYPE_UNDEFINED
	case *dvgrammar.ExpressionValue:
		f = v.(*dvgrammar.ExpressionValue).DataType
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
	case *dvgrammar.ExpressionValue:
		d := v.(*dvgrammar.ExpressionValue)
		switch d.DataType {
		case dvgrammar.TYPE_NULL:
			f = 0
		default:
			return AnyToNumber(d.Value)
		}
	case *DvVariable:
		b := v.(*DvVariable)
		switch b.Kind {
		case FIELD_OBJECT, FIELD_ARRAY:
			f = float64(len(b.Fields))
		case FIELD_NULL, FIELD_UNDEFINED:
			f = 0
		default:
			return AnyToNumber(string(b.Value))
		}
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
	if v == nil {
		return false
	}
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
	case *DvVariable:
		d := v.(*DvVariable)
		f = d != nil
		if f {
			switch d.Kind {
			case FIELD_NULL, FIELD_UNDEFINED:
				f = false
			case FIELD_STRING:
				f = len(d.Value) > 0
			case FIELD_NUMBER:
				f = len(d.Value) > 0 && string(d.Value) != "0"
			case FIELD_BOOLEAN:
				f = len(d.Value) > 0 && d.Value[0] != 'f'
			}
		}
	case *dvgrammar.ExpressionValue:
		b := v.(*dvgrammar.ExpressionValue)
		f = b != nil
		if f {
			switch b.DataType {
			case dvgrammar.TYPE_NULL, dvgrammar.TYPE_NAN:
				f = false
				/*			case dvgrammar.TYPE_NUMBER:
								f = AnyToNumber(b.Value)!=0
							case dvgrammar.TYPE_NUMBER_INT:
								i,ok := AnyToNumberInt(b.Value)
								f = ok && i==0
				*/
			default:
				f = AnyToBoolean(b.Value)
			}
		}
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
	case *DvVariable:
		vr := v.(*DvVariable)
		switch vr.Kind {
		case FIELD_OBJECT, FIELD_ARRAY, FIELD_FUNCTION:
			f = int64(len(vr.Fields))
		case FIELD_UNDEFINED, FIELD_NULL:
			f = 0
		case FIELD_BOOLEAN:
			if len(vr.Value) > 0 && vr.Value[0] == 't' {
				f = 1
			} else {
				f = 0
			}
		default:
			return AnyToNumberInt(string(vr.Value))
		}
	case nil:
		f = 0
	case *dvgrammar.ExpressionValue:
		b := v.(*dvgrammar.ExpressionValue)
		if b.DataType == dvgrammar.TYPE_NULL || b.DataType == dvgrammar.TYPE_NAN {
			f = 0
		} else {
			return AnyToNumberInt(b.Value)
		}
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

func AnyToDvVariable(v interface{}) *DvVariable {
	switch v.(type) {
	case *DvVariable:
		return v.(*DvVariable)
	case bool:
		d := &DvVariable{
			Kind: FIELD_BOOLEAN,
		}
		if v.(bool) {
			d.Value = []byte("true")
		} else {
			d.Value = []byte("false")
		}
		return d
	case string:
		return &DvVariable{Kind: FIELD_STRING, Value: []byte(v.(string))}
	case int:
		return &DvVariable{Kind: FIELD_NUMBER, Value: []byte(strconv.Itoa(v.(int)))}
	case float64, float32, int64, int32:
		s := AnyToString(v)
		return &DvVariable{Kind: FIELD_NUMBER, Value: []byte(s)}
	case *dvgrammar.ExpressionValue:
		return AnyToDvVariable(v.(*dvgrammar.ExpressionValue).Value)
	case []string:
		rs := v.([]string)
		rn := len(rs)
		rd := &DvVariable{Kind: FIELD_ARRAY, Fields: make([]*DvVariable, rn)}
		for ri := 0; ri < rn; ri++ {
			rd.Fields[ri] = &DvVariable{
				Kind:  FIELD_STRING,
				Value: []byte(rs[ri]),
			}
		}
		return rd
	}
	return nil
}

func ConvertInterfaceListsMapToStringListsMap(listMap map[string][][]interface{}, options int) map[string][][]string {
	r := make(map[string][][]string, len(listMap))
	for k, v := range listMap {
		r[k] = ConvertInterfaceListsToStringLists(v, options)
	}
	return r
}

func CreateDvVariableByPathAndData(path string, data interface{}, parent *DvVariable) *DvVariable {
	p := strings.Index(path, ".")
	if p == 0 {
		return CreateDvVariableByPathAndData(path[1:], data, parent)
	}
	if p > 0 {
		name := path[:p]
		if parent == nil {
			parent = &DvVariable{Kind: FIELD_OBJECT, Fields: make([]*DvVariable, 0, 7)}
		}
		item := &DvVariable{Kind: FIELD_OBJECT, Value: []byte(name)}
		if parent.Fields == nil {
			parent.Fields = make([]*DvVariable, 0, 7)
		}
		parent.Fields = append(parent.Fields, item)
		CreateDvVariableByPathAndData(path[p+1:], data, item)
		return parent
	}
	dvvar := ConvertAnyToDvVariable(data)
	if path == "" {
		if parent == nil {
			return dvvar
		}
		parent.CloneExceptKey(dvvar, false)
		return parent
	}
	dvvar.Name = []byte(path)
	if parent == nil {
		return &DvVariable{Kind: FIELD_OBJECT, Fields: []*DvVariable{dvvar}}
	}
	parent.CloneWithKey(dvvar, false)
	return parent
}

func AnyToDvGrammarExpressionValue(v interface{}) *dvgrammar.ExpressionValue {
	if v == nil {
		return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NULL}
	}
	switch v.(type) {
	case string:
		return &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_STRING}
	case *DvVariable:
		return v.(*DvVariable).ToDvGrammarExpressionValue()
	case int64:
		return &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_NUMBER_INT}
	case int:
		return &dvgrammar.ExpressionValue{Value: int64(v.(int)), DataType: dvgrammar.TYPE_NUMBER_INT}
	case float64:
		return &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_NUMBER}
	case bool:
		return &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_BOOLEAN}
	case nil:
		return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NULL}
	default:
		d := AnyToDvVariable(v)
		if d != nil {
			return d.ToDvGrammarExpressionValue()
		}
	}
	log.Printf("Unknown type %v", v)
	return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NULL}
}
