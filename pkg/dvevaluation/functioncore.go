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
	}
	return res, nil
}

func ExtractPureName(tree *dvgrammar.BuildNode) (string, error) {
	return "", nil
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

func PutVariablesInScope(params []string, nodes []*dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) error {
	n := len(params)
	m := len(nodes)
	var v *dvgrammar.ExpressionValue
	var err error
	for i := 0; i < n; i++ {
		v = nil
		if i < m {
			_, v, err = nodes[i].ExecuteExpression(context)
			if err != nil {
				return err
			}
		}
		context.Scope.Set(params[i], v)
	}
	return nil
}
