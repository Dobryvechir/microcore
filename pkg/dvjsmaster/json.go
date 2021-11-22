/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"strconv"
)

var JSONMaster *dvevaluation.DvVariable

func json_init() {
	JSONMaster = dvevaluation.RegisterMasterVariable("JSON", &dvevaluation.DvVariable{
		Fields: make(map[string]*dvevaluation.DvVariable),
		Kind:   dvevaluation.FIELD_OBJECT,
		Prototype: &dvevaluation.DvVariable{
			Fields: map[string]*dvevaluation.DvVariable{
				"stringify": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   JSON_stringify,
				},
				"parse": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   JSON_parse,
				},
			},
			Kind: dvevaluation.FIELD_OBJECT,
		},
	})

}

func JSON_stringify(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	var data []byte
	if len(params) == 0 || params[0] == nil {
		data = []byte{}
	} else {
		data = params[1].JsonStringify()
	}
	res := dvevaluation.DvVariableFromString(nil, string(data))
	return res, nil
}

func JSON_parse(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	if len(params) == 0 || params[0] == nil {
		return &dvevaluation.DvVariable{}, nil
	}
	return JSON_parse_direct([]byte(params[1].GetStringValue()), "JSON.parse")
}

func convert_Object_DvFieldInfo_to_DvVariableMap(data []*dvjson.DvFieldInfo) (res map[string]*dvevaluation.DvVariable) {
	res = make(map[string]*dvevaluation.DvVariable)
	n := len(data)
	for i := 0; i < n; i++ {
		fld := data[i]
		res[string(fld.Name)] = convert_DvFieldInfo_to_DvVariable(fld)
	}
	return
}

func convert_Array_DvFieldInfo_to_DvVariableMap(data []*dvjson.DvFieldInfo) (res map[string]*dvevaluation.DvVariable) {
	res = make(map[string]*dvevaluation.DvVariable)
	n := len(data)
	for i := 0; i < n; i++ {
		res[strconv.Itoa(i)] = convert_DvFieldInfo_to_DvVariable(data[i])
	}
	res[dvevaluation.LENGTH_PROPERTY] = dvevaluation.DvVariableFromInt(nil, n)
	return
}

func convert_DvFieldInfo_to_DvVariable(field *dvjson.DvFieldInfo) *dvevaluation.DvVariable {
	parent := &dvevaluation.DvVariable{}
	switch field.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Kind = dvevaluation.FIELD_OBJECT
		parent.Prototype = dvevaluation.ObjectMaster
		parent.Fields = convert_Object_DvFieldInfo_to_DvVariableMap(field.Fields)
	case dvevaluation.FIELD_ARRAY:
		parent.Fields = convert_Array_DvFieldInfo_to_DvVariableMap(field.Fields)
		parent.Kind = dvevaluation.FIELD_ARRAY
		parent.Prototype = dvevaluation.ArrayMaster
	case dvevaluation.FIELD_STRING:
		parent.Kind = dvevaluation.FIELD_STRING
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_NUMBER:
		parent.Kind = dvevaluation.FIELD_NUMBER
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_BOOLEAN:
		parent.Kind = dvevaluation.FIELD_BOOLEAN
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_NULL:
		parent.Kind = dvevaluation.FIELD_NULL
	}
	return parent
}

func convert_Object_DvCrudItem_to_DvVariableMap(data []*dvjson.DvCrudItem) (res map[string]*dvevaluation.DvVariable) {
	res = make(map[string]*dvevaluation.DvVariable)
	n := len(data)
	for i := 0; i < n; i++ {
		fld := data[i]
		res[string(fld.Name)] = convert_DvCrudItem_to_DvVariable(fld)
	}
	return
}

func convert_Array_DvCrudItem_to_DvVariableMap(data []*dvjson.DvCrudItem) (res map[string]*dvevaluation.DvVariable) {
	res = make(map[string]*dvevaluation.DvVariable)
	n := len(data)
	for i := 0; i < n; i++ {
		res[strconv.Itoa(i)] = convert_DvCrudItem_to_DvVariable(data[i])
	}
	res[dvevaluation.LENGTH_PROPERTY] = dvevaluation.DvVariableFromInt(nil, n)
	return
}

func convert_DvCrudItem_to_DvVariable(field *dvjson.DvCrudItem) *dvevaluation.DvVariable {
	parent := &dvevaluation.DvVariable{}
	switch field.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Kind = dvevaluation.FIELD_OBJECT
		parent.Prototype = dvevaluation.ObjectMaster
		parent.Fields = convert_Object_DvFieldInfo_to_DvVariableMap(field.Fields)
	case dvevaluation.FIELD_ARRAY:
		parent.Fields = convert_Array_DvFieldInfo_to_DvVariableMap(field.Fields)
		parent.Kind = dvevaluation.FIELD_ARRAY
		parent.Prototype = dvevaluation.ArrayMaster
	case dvevaluation.FIELD_STRING:
		parent.Kind = dvevaluation.FIELD_STRING
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_NUMBER:
		parent.Kind = dvevaluation.FIELD_NUMBER
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_BOOLEAN:
		parent.Kind = dvevaluation.FIELD_BOOLEAN
		parent.Value = string(field.Value)
	case dvevaluation.FIELD_NULL:
		parent.Kind = dvevaluation.FIELD_NULL
	}
	return parent
}

func JSON_parse_direct(body []byte, info string) (*dvevaluation.DvVariable, error) {
	parent := &dvevaluation.DvVariable{}
	if len(body) == 0 {
		return parent, nil
	}
	crudDetails := &dvjson.DvCrudDetails{}
	highLevelObject := false
	parsed := dvjson.JsonQuickParser(body, crudDetails, highLevelObject, dvjson.OPTIONS_FIELDS_DETAILED)
	if parsed.Err != "" {
		return nil, errors.New(parsed.Err)
	}
	switch parsed.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Kind = dvevaluation.FIELD_OBJECT
		parent.Prototype = dvevaluation.ObjectMaster
		parent.Fields = convert_Object_DvCrudItem_to_DvVariableMap(parsed.Items)
	case dvevaluation.FIELD_ARRAY:
		parent.Fields = convert_Array_DvCrudItem_to_DvVariableMap(parsed.Items)
		parent.Kind = dvevaluation.FIELD_ARRAY
		parent.Prototype = dvevaluation.ArrayMaster
	case dvevaluation.FIELD_STRING:
		parent.Kind = dvevaluation.FIELD_STRING
		parent.Value = string(parsed.Value)
	case dvevaluation.FIELD_NUMBER:
		parent.Kind = dvevaluation.FIELD_NUMBER
		parent.Value = string(parsed.Value)
	case dvevaluation.FIELD_BOOLEAN:
		parent.Kind = dvevaluation.FIELD_BOOLEAN
		parent.Value = string(parsed.Value)
	case dvevaluation.FIELD_NULL:
		parent.Kind = dvevaluation.FIELD_NULL
	}
	return parent, nil
}
