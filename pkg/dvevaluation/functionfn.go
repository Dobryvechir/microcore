/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"strings"
)

func FunctionCall(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	var selfParam *DvVariable = nil
	if len(params) > 0 {
		selfParam = params[0]
		params = params[1:]
	}
	return thisVariable.Fn(context, selfParam, params)
}

func FunctionApply(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	var selfParam *DvVariable = nil
	var dvParams []*DvVariable
	var err error
	if len(params) > 0 {
		selfParam = params[0]
		if len(params) > 1 {
			dvParams, err = params[1].GetVariableArray(true)
			if err != nil {
				return nil, err
			}
		}
	}
	return thisVariable.Fn(context, selfParam, dvParams)
}

func FunctionBind(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	if len(params) == 0 {
		return thisVariable, nil
	}
	paramThis := params[0]
	dvParams := make([]*DvVariable, len(params)-1)
	copy(dvParams, params[1:])
	return &DvVariable{Kind: FIELD_FUNCTION, Prototype: FunctionMaster, Fn: func(context *DvContext, thisVar *DvVariable, pars []*DvVariable) (*DvVariable, error) {
		newParams := append(dvParams, pars...)
		return thisVariable.Fn(context, paramThis, newParams)
	}}, nil
}

func FunctionBindLite(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	if len(params) == 0 {
		return thisVariable, nil
	}
	paramThis := params[0]
	dvParams := make([]*DvVariable, len(params)-1)
	copy(dvParams, params[1:])
	return &DvVariable{Kind: FIELD_FUNCTION, Fn: func(context *DvContext, thisVar *DvVariable, pars []*DvVariable) (*DvVariable, error) {
		newParams := append(dvParams, pars...)
		return thisVariable.Fn(context, paramThis, newParams)
	}}, nil
}

var FunctionMaster *DvVariable = RegisterMasterVariable("Function", &DvVariable{
	Fields: make(map[string]*DvVariable),
	Kind:   FIELD_OBJECT,
	Prototype: &DvVariable{
		Fields: map[string]*DvVariable{
			"call": {
				Kind: FIELD_FUNCTION,
				Fn:   FunctionCall,
			},
			"apply": {
				Kind: FIELD_FUNCTION,
				Fn:   FunctionApply,
			},
			"bind": {
				Kind: FIELD_FUNCTION,
				Fn:   FunctionBindLite,
			},
		},
		Kind: FIELD_OBJECT,
	},
})

func functionfn_init() {
	FunctionMaster.Prototype.Fields["bind"].Fn = FunctionBind
}

func (context *DvContext) FunctionCallByKeys(variable *DvVariable, keys []string, params []*DvVariable, thisVariable *DvVariable) (*DvVariable, error) {
	n := len(keys)
	parent := variable
	for i := 0; i < n; i++ {
		v, err := variable.ObjectGet(keys[i])
		if err != nil {
			return nil, err
		}
		parent = variable
		variable = v
	}
	if thisVariable == nil {
		thisVariable = parent
	}
	if variable == nil || variable.Kind != FIELD_FUNCTION {
		return nil, errors.New("Cannot call not function:" + strings.Join(keys, "."))
	}
	return variable.Fn(context, thisVariable, params)
}

func (context *DvContext) FunctionCallByVariableDefinition(variable *DvVariable, variableDefinition string, params []*DvVariable, thisVariable *DvVariable) (*DvVariable, error) {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return nil, err
	}
	return context.FunctionCallByKeys(variable, keys, params, thisVariable)
}

func (context *DvContext) FunctionCallByVariableDefinitionWithStringParams(variable *DvVariable, variableDefinition string, params []string, thisVariable *DvVariable) (*DvVariable, error) {
	dvparams := ConvertStringArrayToDvVariableArray(params)
	return context.FunctionCallByVariableDefinition(variable, variableDefinition, dvparams, thisVariable)
}

func ConvertDvFunctionToDvVariable(fn DvvFunction) *DvVariable {
	return &DvVariable{Kind: FIELD_FUNCTION, Fn: fn, Prototype: FunctionMaster}
}
