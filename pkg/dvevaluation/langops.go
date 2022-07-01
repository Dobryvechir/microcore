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
	n := len(tree.Children)
	if n == 0 {
		return dvgrammar.FLOW_RETURN, nil, nil
	}
	_, val, err := tree.Children[n-1].ExecuteExpression(context)
	return dvgrammar.FLOW_RETURN, val, err
}

func ArrowFunctionOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n != 2 {
		return dvgrammar.FLOW_NORMAL, nil, errors.New("Arrow function requires 2 parameters")
	}
	params, err := GetFunctionParameterList(tree.Children[0])
	if err != nil {
		return dvgrammar.FLOW_NORMAL, nil, err
	}
	code, err := GetFunctionCodeList(tree.Children[1])
	if err != nil {
		return dvgrammar.FLOW_NORMAL, nil, err
	}
	val, err := CreateFunctionContainer(params, code, FUNCTION_KIND_ARROW, "")
	return dvgrammar.FLOW_NORMAL, val, err
}

func ForCycleOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n != 1 || tree.Children[0] == nil || len(tree.Children[0].Children) < 2 {
		return dvgrammar.FLOW_NORMAL, nil, errors.New("For cycle requires only parentheses and curly brackets")
	}
	forInside := tree.Children[0].Children[0].Children
	cycleBody := tree.Children[0].Children[1].Children
	initValues := collectAtLevel(forInside, 0)
	condValue := collectAtLevel(forInside, 1)
	stepValues := collectAtLevel(forInside, 2)
	flow, val, err := ExecuteCycleCommon(context, initValues, condValue, stepValues, cycleBody, true)
	return flow, val, err
}

func collectAtLevel(src []*dvgrammar.BuildNode, group int) []*dvgrammar.BuildNode {
	n := len(src)
	m := 0
	for i := 0; i < n; i++ {
		if src[i] != nil && src[i].Group == group {
			m++
		}
	}
	res := make([]*dvgrammar.BuildNode, m)
	m = 0
	for i := 0; i < n; i++ {
		if src[i] != nil && src[i].Group == group {
			res[m] = src[i]
			m++
		}
	}
	return res
}

func ExecuteCycleCommon(context *dvgrammar.ExpressionContext, initValues []*dvgrammar.BuildNode, condValue []*dvgrammar.BuildNode, stepValues []*dvgrammar.BuildNode, cycleBody []*dvgrammar.BuildNode, condAtFirst bool) (int, *dvgrammar.ExpressionValue, error) {
	if len(condValue) > 1 {
		return 0, nil, errors.New("The cycle must have no more than one condition")
	}
	if len(initValues) > 0 {
		_, _, err := dvgrammar.BuildNodeExecution(initValues, context)
		if err != nil {
			return dvgrammar.FLOW_NORMAL, nil, err
		}
	}
	condPresent := len(condValue) == 1 && condValue[0] != nil
	for i := 0; i < 1000000000; i++ {
		flow, val, err := dvgrammar.BuildNodeExecution(cycleBody, context)
		if err != nil || flow != dvgrammar.FLOW_NORMAL {
			return flow, val, err
		}
		_, _, err = dvgrammar.BuildNodeExecution(stepValues, context)
		if err != nil {
			return dvgrammar.FLOW_NORMAL, nil, err
		}
		if condPresent && (i > 0 || condAtFirst) {
			_, val, err = condValue[0].ExecuteExpression(context)
			if err != nil {
				return dvgrammar.FLOW_NORMAL, nil, err
			}
			b := AnyToBoolean(val)
			if !b {
				break
			}
		}
	}
	return dvgrammar.FLOW_NORMAL, nil, nil
}
