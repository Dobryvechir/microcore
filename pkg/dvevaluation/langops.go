/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
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
