/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strings"
)

func ProcessorAssign(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(tree.Children)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	valueRight, err := tree.GetChildrenExpressionValue(1, context)
	if err != nil {
		return nil, err
	}
	valueLeft, err := tree.GetChildrenExpressionValue(0, context)
	if err != nil && !strings.Contains(err.Error(),"is not defined") {
		return nil, err
	}
	if valueLeft == nil || valueLeft.Name == "" {
		return nil, errors.New("Invalid left-hand side in assignment")
	}
	var valueRightDirect interface{} = nil
	if valueRight!=nil {
		valueRightDirect = valueRight.Value
	}
	context.Scope.SetDeep(valueLeft.Name, valueRightDirect)
	return valueRight, nil
}

func reassign(res *dvgrammar.ExpressionValue, err error, values []*dvgrammar.ExpressionValue, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	if err != nil {
		return nil, err
	}
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	valueLeft := values[0]
	if valueLeft == nil || valueLeft.Name == "" {
		return nil, errors.New("Invalid left-hand side in assignment")
	}
	context.Scope.SetDeep(valueLeft.Name, res)
	return res, nil
}

func ProcessorPlusAssign(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	res, err := ProcessorPlus(values, tree, context, operator)
	return reassign(res, err, values, context, operator)
}
