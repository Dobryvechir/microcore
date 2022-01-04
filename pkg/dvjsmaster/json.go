/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
)

var JSONMaster *dvevaluation.DvVariable

func json_init() {
	JSONMaster = dvevaluation.RegisterMasterVariable("JSON", &dvevaluation.DvVariable{
		Fields: make([]*dvevaluation.DvVariable, 0, 7),
		Kind:   dvevaluation.FIELD_OBJECT,
		Prototype: &dvevaluation.DvVariable{
			Fields: []*dvevaluation.DvVariable{
				{
					Name: []byte("stringify"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra:   JSON_stringify,
				},
				{
					Name: []byte("parse"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra:   JSON_parse,
				},
			},
			Kind: dvevaluation.FIELD_OBJECT,
		},
	})

}

func JSON_stringify(context , thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
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

func convert_DvFieldInfo_to_DvVariable(field *dvevaluation.DvVariable) *dvevaluation.DvVariable {
	parent := &dvevaluation.DvVariable{Kind: field.Kind, Value: field.Value, Name: field.Name, Fields: field.Fields}
	switch field.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Prototype = dvevaluation.ObjectMaster
	case dvevaluation.FIELD_ARRAY:
		parent.Prototype = dvevaluation.ArrayMaster
	}
	return parent
}

func convert_Object_DvCrudItem_to_DvVariableMap(data []*dvjson.DvCrudItem) (res []*dvevaluation.DvVariable) {
	n := len(data)
	res = make([]*dvevaluation.DvVariable, n)
	for i := 0; i < n; i++ {
		fld := data[i]
		res[i] = convert_DvCrudItem_to_DvVariable(fld)
	}
	return
}

func convert_Array_DvCrudItem_to_DvVariableMap(data []*dvjson.DvCrudItem) (res []*dvevaluation.DvVariable) {
	n := len(data)
	res = make([]*dvevaluation.DvVariable, n)
	for i := 0; i < n; i++ {
		res[i] = convert_DvCrudItem_to_DvVariable(data[i])
	}
	return
}

func convert_DvCrudItem_to_DvVariable(field *dvjson.DvCrudItem) *dvevaluation.DvVariable {
	parent := &dvevaluation.DvVariable{Kind: field.Kind, Value: field.Value, Name: field.Name, Fields: field.Fields}
	switch field.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Prototype = dvevaluation.ObjectMaster
	case dvevaluation.FIELD_ARRAY:
		parent.Prototype = dvevaluation.ArrayMaster
	}
	return parent
}

func JSON_parse_direct(body []byte, info string) (*dvevaluation.DvVariable, error) {
	if len(body) == 0 {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_NULL}, nil
	}
	crudDetails := &dvjson.DvCrudDetails{}
	highLevelObject := false
	parsed := dvjson.JsonQuickParser(body, crudDetails, highLevelObject, dvjson.OPTIONS_FIELDS_DETAILED)
	if parsed.Err != "" {
		return nil, errors.New(parsed.Err)
	}
	parent := &dvevaluation.DvVariable{Kind: parsed.Kind, Value: parsed.Value}
	switch parsed.Kind {
	case dvevaluation.FIELD_OBJECT:
		parent.Prototype = dvevaluation.ObjectMaster
		parent.Fields = convert_Object_DvCrudItem_to_DvVariableMap(parsed.Items)
	case dvevaluation.FIELD_ARRAY:
		parent.Fields = convert_Array_DvCrudItem_to_DvVariableMap(parsed.Items)
		parent.Prototype = dvevaluation.ArrayMaster
	}
	return parent, nil
}
