/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strings"
)

func FunctionCall(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	var selfParam interface{} = nil
	if len(params) > 0 {
		selfParam = params[0]
		params = params[1:]
	}
	return functionCallInner(context, thisVariable, params, selfParam)
}
func functionCallInner(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{},selfParam interface{})  (interface{}, error) {
	switch thisVariable.(type) {
	case *DvVariable:
		d:=thisVariable.(*DvVariable)
		if d!=nil && d.Kind==FIELD_FUNCTION && d.Extra!=nil {
			return d.Extra.(*DvFunction).Fn(context, selfParam, params)
		}
	case *DvFunction:
		return thisVariable.(*DvFunction).Fn(context, selfParam, params)
	}
	return nil, errors.New("Cannot execute call for non functions")
}

func FunctionApply(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	var selfParam interface{} = nil
	var dvParams []interface{}
	if len(params) > 0 {
		selfParam = params[0]
		if len(params) > 1 {
			v:=params[1]
			switch v.(type) {
			case []interface{}:
				dvParams = v.([]interface{})
			}
		}
	}
	return functionCallInner(context, thisVariable, dvParams, selfParam)
}

func FunctionBind(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if len(params) == 0 {
		return thisVariable, nil
	}
	paramThis := params[0]
	dvParams := make([]interface{}, len(params)-1)
	copy(dvParams, params[1:])
	return &DvVariable{
		Kind: FIELD_FUNCTION,
		Extra: &DvFunction{
			Fn: func(context *dvgrammar.ExpressionContext, thisVar interface{}, pars []interface{}) (interface{}, error) {
				newParams := append(dvParams, pars...)
				return functionCallInner(context, thisVariable, newParams, paramThis)
			},
		},
	}, nil
}

var FunctionMaster *DvVariable = RegisterMasterVariable("Function", &DvVariable{
	Fields: make([]*DvVariable,0,7),
	Kind:   FIELD_OBJECT,
	Prototype: &DvVariable{
		Fields: []*DvVariable{
			{
				Name: []byte("call"),
				Kind: FIELD_FUNCTION,
				Extra: &DvFunction{
					Fn: FunctionCall,
				},
			},
			{
				Name: []byte("apply"),
				Kind: FIELD_FUNCTION,
				Extra:   &DvFunction{
					Fn: FunctionApply,
				},
			},
			{
				Name: []byte("bind"),
				Kind: FIELD_FUNCTION,
				Extra:  &DvFunction{
					Fn: FunctionBind,
				},
			},
		},
		Kind: FIELD_OBJECT,
	},
})

func (context *DvContext) FunctionCallByKeys(variable *DvVariable, keys []string, params []interface{}, thisVariable *DvVariable) (interface{}, error) {
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
	if variable == nil || variable.Kind != FIELD_FUNCTION || variable.Extra==nil {
		return nil, errors.New("Cannot call not function:" + strings.Join(keys, "."))
	}
	return variable.Extra.(*DvFunction).Fn(nil, thisVariable, params)
}

func (context *DvContext) FunctionCallByVariableDefinition(variable *DvVariable, variableDefinition string, params []interface{}, thisVariable *DvVariable) (interface{}, error) {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return nil, err
	}
	return context.FunctionCallByKeys(variable, keys, params, thisVariable)
}

func (context *DvContext) FunctionCallByVariableDefinitionWithStringParams(variable *DvVariable, variableDefinition string, params []string, thisVariable *DvVariable) (interface{}, error) {
	dvparams := ConvertStringArrayToInterfaceArray(params)
	return context.FunctionCallByVariableDefinition(variable, variableDefinition, dvparams, thisVariable)
}

func ConvertDvFunctionToDvVariable(fn *DvFunction) *DvVariable {
	return &DvVariable{Kind: FIELD_FUNCTION, Extra: fn}
}
