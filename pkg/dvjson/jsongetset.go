/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvjson

import ()

func ConvertDvFieldInfoArrayIntoMap(data []*DvFieldInfo) map[string]*DvFieldInfo {
	res := make(map[string]*DvFieldInfo)
	for _, v := range data {
		res[string(v.Name)] = v
	}
	return res
}

func GetIntValueFromFieldMap(data map[string]*DvFieldInfo, fieldName string) (int, bool) {
	item, ok := data[fieldName]
	if !ok || item.Kind != FIELD_NUMBER {
		return 0, false
	}
	val, ok1 := ConvertByteArrayToIntOrDouble(item.Value)
	if !ok1 {
		return 0, false
	}
	res, ok2 := val.(int)
	if !ok2 {
		return 0, false
	}
	return res, true
}

func GetStringValueFromFieldMap(data map[string]*DvFieldInfo, fieldName string) (string, bool) {
	item, ok := data[fieldName]
	if !ok || item.Kind != FIELD_STRING {
		return "", false
	}
	return string(item.Value), true
}

func CreateDvFieldInfoObject() *DvFieldInfo {
	return &DvFieldInfo{Kind: FIELD_OBJECT, Fields: make([]*DvFieldInfo, 0, 7)}
}

func (field *DvFieldInfo) AddStringField(key string, value string) bool {
	if field.Fields == nil {
		field.Fields = make([]*DvFieldInfo, 0, 7)
	}
	field.Fields = append(field.Fields, &DvFieldInfo{Kind: FIELD_STRING, Name: []byte(key), Value: []byte(value)})
	return true
}

func (field *DvFieldInfo) AddField(item *DvFieldInfo) bool {
	if field.Fields == nil {
		field.Fields = make([]*DvFieldInfo, 0, 7)
	}
	field.Fields = append(field.Fields, item)
	return true
}
