/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func (parseInfo *DvCrudParsingInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return dvevaluation.ConvertSimpleKindAndValueToInterface(parseInfo.Kind, parseInfo.Value)
}

func DvFieldArrayToBytes(val []*dvevaluation.DvVariable) []byte {
	buf := make([]byte, 1, 102400)
	buf[0] = '['
	n := len(val)
	for i := 0; i < n; i++ {
		if i != 0 {
			buf = append(buf, ',')
		}
		b := PrintToJson(val[i], 2)
		buf = append(buf, b...)
	}
	buf = append(buf, ']')
	return buf
}

func DvFieldInfoToStringConverter(v interface{}) (string, bool) {
	switch v.(type) {
	case *dvevaluation.DvVariable:
		return string(PrintToJson(v.(*dvevaluation.DvVariable), 2)), true
	case []*dvevaluation.DvVariable:
		return string(DvFieldArrayToBytes(v.([]*dvevaluation.DvVariable))), true
	}
	return "", false
}

func ConvertStringToStringIntoStringToInterface(v map[string]string) map[string]interface{} {
	res := make(map[string]interface{}, len(v))
	for key, val := range v {
		res[key] = val
	}
	return res
}

func ConvertStringToInterfaceIntoStringToString(v map[string]interface{}) map[string]string {
	res := make(map[string]string, len(v))
	for key, val := range v {
		res[key] = dvevaluation.AnyToString(val)
	}
	return res
}

func ConvertStringToDvVariableIntoStringToInterface(v map[string]*dvevaluation.DvVariable) map[string]interface{} {
	res := make(map[string]interface{}, len(v))
	for key, val := range v {
		res[key] = val
	}
	return res
}

func ConvertStringToDvVariableIntoStringToString(v map[string]*dvevaluation.DvVariable) map[string]string {
	res := make(map[string]string, len(v))
	for key, val := range v {
		res[key] = val.GetStringValue()
	}
	return res
}

func ConvertInterfaceIntoMap(v interface{}) (map[string]interface{}, bool) {
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{}), true
	case map[string]string:
		return ConvertStringToStringIntoStringToInterface(v.(map[string]string)), true
	case map[string]*dvevaluation.DvVariable:
		return ConvertStringToDvVariableIntoStringToInterface(v.(map[string]*dvevaluation.DvVariable)), true
	case *dvevaluation.DvVariable:
		f := v.(*dvevaluation.DvVariable)
		if f.Kind == dvevaluation.FIELD_OBJECT || f.Kind == dvevaluation.FIELD_ARRAY {
			return f.GetStringInterfaceMap(), true
		}
	}
	return nil, false
}

func ConvertInterfaceIntoStringMap(v interface{}) (map[string]string, bool) {
	switch v.(type) {
	case map[string]interface{}:
		return ConvertStringToInterfaceIntoStringToString(v.(map[string]interface{})), true
	case map[string]string:
		return v.(map[string]string), true
	case map[string]*dvevaluation.DvVariable:
		return ConvertStringToDvVariableIntoStringToString(v.(map[string]*dvevaluation.DvVariable)), true
	case *dvevaluation.DvVariable:
		f := v.(*dvevaluation.DvVariable)
		if f.Kind == dvevaluation.FIELD_OBJECT || f.Kind == dvevaluation.FIELD_ARRAY {
			return f.GetStringMap(), true
		}
	}
	return nil, false
}
