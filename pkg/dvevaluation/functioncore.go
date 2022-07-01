/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

const (
	FUNCTION_KIND_ARROW = iota
	FUNCTION_KIND_NORMAL
)

type CustomJsFunction struct {
	Params  []string
	Body    []*dvgrammar.BuildNode
	Options int
}

func CreateFunctionContainer(params []string, body []*dvgrammar.BuildNode, options int, name string) (*dvgrammar.ExpressionValue, error) {
	jsfunc := &CustomJsFunction{
		Params:  params,
		Body:    body,
		Options: options,
	}
	gram := &dvgrammar.ExpressionValue{
		Value:    jsfunc,
		DataType: dvgrammar.TYPE_FUNCTION,
		Name:     name,
	}
	return gram, nil
}

func GetFunctionParameterList(tree *dvgrammar.BuildNode) ([]string, error) {
	if tree == nil || len(tree.Children) != 1 || tree.Value == nil || tree.Value.DataType != dvgrammar.TYPE_FUNCTION {
		return nil, errors.New("Expected parameters in round brackets only")
	}
	tree = tree.Children[0]
	if tree == nil || tree.Operator != "(" {
		return nil, errors.New("Only expected parameters in round brackets")
	}
	n := len(tree.Children)
	if n == 0 {
		return nil, nil
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		p, err := ExtractPureName(tree.Children[i])
		if err != nil {
			return nil, err
		}
		pos := dvtextutils.IsValidVariableName(p, true)
		if pos >= 0 {
			return nil, fmt.Errorf("Bad variable name %s at %d", p, pos)
		}
		res[i] = p
	}
	return res, nil
}

func ExtractPureName(tree *dvgrammar.BuildNode) (string, error) {
	if tree == nil || tree.Value == nil || tree.Value.DataType != dvgrammar.TYPE_DATA && tree.Value.DataType != dvgrammar.TYPE_STRING {
		return "", errors.New("Parameter name required")
	}
	return tree.Value.Value, nil
}

func GetFunctionCodeList(tree *dvgrammar.BuildNode) ([]*dvgrammar.BuildNode, error) {
	if tree == nil || len(tree.Children) == 0 || tree.Value == nil || tree.Value.DataType != dvgrammar.TYPE_FUNCTION {
		return nil, nil
	}
	tree = tree.Children[0]
	if tree == nil || tree.Operator != "{" {
		return nil, nil
	}
	return tree.Children, nil
}

func CalculateAllNodeParams(args []*dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) ([]interface{}, error) {
	n := len(args)
	interfaceArgs := make([]interface{}, n)
	var err error
	for i := 0; i < n; i++ {
		_, interfaceArgs[i], err = args[i].ExecuteExpression(context)
		if err != nil {
			return nil, err
		}
	}
	return interfaceArgs, nil
}

func PutVariablesInScope(params []string, args []interface{}, context *dvgrammar.ExpressionContext) error {
	n := len(params)
	m := len(args)
	var v interface{}
	for i := 0; i < n; i++ {
		v = nil
		if i < m {
			v = args[i]
		}
		context.Scope.Set(params[i], v)
	}
	return nil
}

func ExecuteAnyFunction(context *dvgrammar.ExpressionContext, fn interface{}, thisArg interface{}, args []interface{}) (value interface{}, err error) {
	switch fn.(type) {
	case *DvVariable:
		dv := fn.(*DvVariable)
		if dv.Kind == FIELD_FUNCTION && dv.Extra != nil {
			switch dv.Extra.(type) {
			case *DvFunctionObject:
				value, err = dv.Extra.(*DvFunctionObject).ExecuteDvFunctionWithTreeArguments(args, context)
				return
			case *DvFunction:
				functionObject := &DvFunctionObject{
					SelfRef:  thisArg,
					Context:  context,
					Executor: dv.Extra.(*DvFunction),
				}
				value, err = functionObject.ExecuteDvFunctionWithTreeArguments(args, context)
				return
			}
		}
	case *CustomJsFunction:
		cf := fn.(*CustomJsFunction)
		context.Scope.StackPush(cf.Options)
		err = PutVariablesInScope(cf.Params, args, context)
		if err == nil {
			_, value, err = dvgrammar.BuildNodeExecution(cf.Body, context)
		}
		context.Scope.StackPop()
		return
	case *dvgrammar.ExpressionValue:
		dvg := fn.(*dvgrammar.ExpressionValue)
		if dvg != nil && dvg.DataType == dvgrammar.TYPE_FUNCTION && dvg.Value != nil {
			return ExecuteAnyFunction(context, dvg.Value, thisArg, args)
		}
	}
	return nil, fmt.Errorf("Value of %v is not a function", fn)
}
