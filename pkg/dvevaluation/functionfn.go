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
	return thisVariable.Extra.(DvvFunction)(context, selfParam, params)
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
	return thisVariable.Extra.(DvvFunction)(context, selfParam, dvParams)
}

func FunctionBind(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	if len(params) == 0 {
		return thisVariable, nil
	}
	paramThis := params[0]
	dvParams := make([]*DvVariable, len(params)-1)
	copy(dvParams, params[1:])
	return &DvVariable{Kind: FIELD_FUNCTION, Prototype: FunctionMaster, Extra: func(context *DvContext, thisVar *DvVariable, pars []*DvVariable) (*DvVariable, error) {
		newParams := append(dvParams, pars...)
		return thisVariable.Extra.(DvvFunction)(context, paramThis, newParams)
	}}, nil
}

func FunctionBindLite(context *DvContext, thisVariable *DvVariable, params []*DvVariable) (*DvVariable, error) {
	if len(params) == 0 {
		return thisVariable, nil
	}
	paramThis := params[0]
	dvParams := make([]*DvVariable, len(params)-1)
	copy(dvParams, params[1:])
	return &DvVariable{Kind: FIELD_FUNCTION, Extra: func(context *DvContext, thisVar *DvVariable, pars []*DvVariable) (*DvVariable, error) {
		newParams := append(dvParams, pars...)
		return thisVariable.Extra.(DvvFunction)(context, paramThis, newParams)
	}}, nil
}

var FunctionMaster *DvVariable = RegisterMasterVariable("Function", &DvVariable{
	Fields: make([]*DvVariable,0,7),
	Kind:   FIELD_OBJECT,
	Prototype: &DvVariable{
		Fields: []*DvVariable{
			{
				Name: []byte("call"),
				Kind: FIELD_FUNCTION,
				Extra:   FunctionCall,
			},
			{
				Name: []byte("apply"),
				Kind: FIELD_FUNCTION,
				Extra:   FunctionApply,
			},
			{
				Name: []byte("bind"),
				Kind: FIELD_FUNCTION,
				Extra:   FunctionBindLite,
			},
		},
		Kind: FIELD_OBJECT,
	},
})

func functionfn_init() {
	FunctionMaster.Prototype.Fields[2].Extra = FunctionBind
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
	return variable.Extra.(DvvFunction)(context, thisVariable, params)
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
	return &DvVariable{Kind: FIELD_FUNCTION, Extra: fn, Prototype: FunctionMaster}
}
