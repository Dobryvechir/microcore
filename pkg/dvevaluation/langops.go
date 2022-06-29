/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func ReturnOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n:=len(tree.Children)
	if n==0 {
		return dvgrammar.FLOW_RETURN, nil, nil
	}
	_, val, err:=tree.Children[n-1].ExecuteExpression(context)
	return dvgrammar.FLOW_RETURN, val, err
}

func ArrowFunctionOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n:=len(tree.Children)
	if n!=2 {
		return dvgrammar.FLOW_NORMAL, nil, errors.New("Arrow function requires 2 parameters")
	}
	params, err:=GetFunctionParameterList(tree.Children[0])
	if err!=nil {
		return dvgrammar.FLOW_NORMAL, nil, err
	}
	code, err:=GetFunctionCodeList(tree.Children[1])
	if err!=nil {
		return dvgrammar.FLOW_NORMAL, nil, err
	}
    val, err:=CreateFunctionContainer(params, code, FUNCTION_KIND_ARROW, "")
	return dvgrammar.FLOW_NORMAL, val, err
}

