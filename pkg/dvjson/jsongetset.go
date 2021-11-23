/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func ConvertDvFieldInfoArrayIntoMap(data []*dvevaluation.DvVariable) map[string]*dvevaluation.DvVariable {
	res := make(map[string]*dvevaluation.DvVariable)
	for _, v := range data {
		res[string(v.Name)] = v
	}
	return res
}

func GetIntValueFromFieldMap(data map[string]*dvevaluation.DvVariable, fieldName string) (int, bool) {
	item, ok := data[fieldName]
	if !ok || item.Kind != dvevaluation.FIELD_NUMBER {
		return 0, false
	}
	val, ok1 := dvevaluation.ConvertByteArrayToIntOrDouble(item.Value)
	if !ok1 {
		return 0, false
	}
	res, ok2 := val.(int)
	if !ok2 {
		return 0, false
	}
	return res, true
}

func GetStringValueFromFieldMap(data map[string]*dvevaluation.DvVariable, fieldName string) (string, bool) {
	item, ok := data[fieldName]
	if !ok || item.Kind != dvevaluation.FIELD_STRING {
		return "", false
	}
	return string(item.Value), true
}

func CreateDvFieldInfoObject() *dvevaluation.DvVariable {
	return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 7)}
}

