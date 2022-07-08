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
	mode, err := checkForInOfCase(initValues, condValue, stepValues)
	if err != nil {
		return dvgrammar.FLOW_NORMAL, nil, err
	}
	if mode >= 0 {
		return ExecuteCycleInOf(context, initValues[0], cycleBody, mode)
	}
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
		if err != nil {
			return flow, val, err
		}
		if flow == dvgrammar.FLOW_BREAK {
			flow = dvgrammar.FLOW_NORMAL
			break
		} else if flow == dvgrammar.FLOW_CONTINUE {
			flow = dvgrammar.FLOW_NORMAL
		} else if flow != dvgrammar.FLOW_NORMAL {
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

func checkForInOfCase(src []*dvgrammar.BuildNode, lev1 []*dvgrammar.BuildNode, lev2 []*dvgrammar.BuildNode) (r int, err error) {
	r = -1
	if len(src) >= 1 && src[0] != nil {
		c := src[0].Operator
		if c == "in" {
			r = 0
		}
		if c == "of" {
			r = 1
		}
	}
	if r >= 0 {
		if len(src) > 1 || len(lev1) > 0 || len(lev2) > 0 {
			return 0, errors.New("Only one operator in 'for' is acceptable for 'in'/'of'")
		}
	}
	return
}

func ExecuteCycleInOf(context *dvgrammar.ExpressionContext, inof *dvgrammar.BuildNode, cycleBody []*dvgrammar.BuildNode, mode int) (int, *dvgrammar.ExpressionValue, error) {
	if inof == nil || len(inof.Children) != 2 || inof.Children[0] == nil || inof.Children[1] == nil || inof.Children[0].Operator != "" {
		return 0, nil, errors.New("in/of parameters are incorrect")
	}
	flow, v, err := inof.Children[1].ExecuteExpression(context)
	if err != nil {
		return flow, v, err
	}
	dv := AnyToDvVariable(v)
	if dv == nil || dv.Kind != FIELD_OBJECT && dv.Kind != FIELD_ARRAY && dv.Kind != FIELD_STRING {
		return 0, nil, errors.New("Cannot cycle by non-object")
	}
	n := len(dv.Fields)
	var r []interface{}
	switch dv.Kind {
	case FIELD_STRING:
		s := string(dv.Value)
		n = len(s)
		if mode == 1 {
			r = make([]interface{}, n)
			for i := 0; i < n; i++ {
				r[i] = string([]byte{s[i]})
			}
		} else {
			r = createIndexArray(n)
		}
	case FIELD_OBJECT:
		if mode == 1 {
			r = createDvVariableArray(dv.Fields)
		} else {
			r = createDvVariableArrayKeys(dv.Fields)
		}
	case FIELD_ARRAY:
		if mode == 1 {
			r = createDvVariableArray(dv.Fields)
		} else {
			r = createIndexArray(n)
		}
	}
	p, err := ExtractPureName(inof.Children[0])
	if err != nil {
		return 0, nil, err
	}
	pos := dvtextutils.IsValidVariableName(p, true)
	if pos >= 0 {
		return 0, nil, fmt.Errorf("Bad variable name %s at %d", p, pos)
	}
	for i := 0; i < n; i++ {
		context.Scope.Set(p, r[i])
		flow, val, err := dvgrammar.BuildNodeExecution(cycleBody, context)
		if err != nil {
			return 0, nil, err
		}
		if flow == dvgrammar.FLOW_CONTINUE {
			flow = dvgrammar.FLOW_NORMAL
		} else if flow == dvgrammar.FLOW_BREAK {
			flow = dvgrammar.FLOW_NORMAL
			break
		} else if flow != dvgrammar.FLOW_NORMAL {
			return flow, val, nil
		}
	}
	return flow, nil, err
}

func createIndexArray(n int) []interface{} {
	r := make([]interface{}, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

func createDvVariableArray(src []*DvVariable) []interface{} {
	n := len(src)
	r := make([]interface{}, n)
	for i := 0; i < n; i++ {
		r[i] = src[i]
	}
	return r
}

func createDvVariableArrayKeys(src []*DvVariable) []interface{} {
	n := len(src)
	r := make([]interface{}, n)
	for i := 0; i < n; i++ {
		if src[i] == nil {
			r[i] = ""
		} else {
			r[i] = string(src[i].Name)
		}
	}
	return r
}

func IfClauseOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n != 1 || tree.Children[0] == nil || len(tree.Children[0].Children) < 2 {
		return dvgrammar.FLOW_NORMAL, nil, errors.New("If clause requires parentheses and curly brackets")
	}
	ifClauseCondition := tree.Children[0].Children[0].Children
	ifThenClause := tree.Children[0].Children[1].Children
	var ifElseClause []*dvgrammar.BuildNode = nil
	if len(tree.Children[0].Children) >= 3 {
		if tree.Children[0].Children[2].Operator == "{" {
			ifElseClause = tree.Children[0].Children[2].Children
		} else {
			ifElseClause = tree.Children[0].Children[2:3]
		}
	}
	flow, val, err := ExecuteIfClauseCommon(context, ifClauseCondition, ifThenClause, ifElseClause)
	return flow, val, err
}

func ExecuteIfClauseCommon(context *dvgrammar.ExpressionContext, ifClauseCondition []*dvgrammar.BuildNode, ifThenClause []*dvgrammar.BuildNode, ifElseClause []*dvgrammar.BuildNode) (int, *dvgrammar.ExpressionValue, error) {
	_, val, err := dvgrammar.BuildNodeExecution(ifClauseCondition, context)
	flow := dvgrammar.FLOW_NORMAL
	if err != nil {
		return flow, nil, err
	}
	v := AnyToBoolean(val)
	val = nil
	if v {
		if len(ifThenClause) > 0 {
			flow, val, err = dvgrammar.BuildNodeExecution(ifThenClause, context)
		}
	} else {
		if len(ifElseClause) > 0 {
			flow, val, err = dvgrammar.BuildNodeExecution(ifElseClause, context)
		}
	}
	return flow, val, err
}

func BreakOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n > 0 {
		if tree.Children[0] != nil && (tree.Children[0].Value != nil || tree.Children[0].Operator != "") {
			return dvgrammar.FLOW_RETURN, nil, errors.New("'break' has no parameters")
		}
	}
	return dvgrammar.FLOW_BREAK, nil, nil
}

func DeleteOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n != 1 || tree.Children[0] == nil {
		return 0, nil, errors.New("'delete' requires one argument")
	}
	oldVisitOptions := context.VisitorOptions
	context.VisitorOptions |= dvgrammar.EVALUATE_OPTION_PARENT | dvgrammar.EVALUATE_OPTION_NAME
	_, val, err := tree.Children[0].ExecuteExpression(context)
	context.VisitorOptions = oldVisitOptions
	if err != nil {
		return 0, nil, err
	}
	if val != nil && val.Name != "" && val.Parent != nil {
		dv := AnyToDvVariable(val.Parent)
		if dv != nil {
			DeleteVariable(dv, []string{val.Name}, true)
		}
	}
	res := &dvgrammar.ExpressionValue{
		DataType: dvgrammar.TYPE_BOOLEAN,
		Value:    true,
	}
	return dvgrammar.FLOW_NORMAL, res, nil
}

func ContinueOperator(tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (int, *dvgrammar.ExpressionValue, error) {
	n := len(tree.Children)
	if n > 0 {
		if tree.Children[0] != nil && (tree.Children[0].Value != nil || tree.Children[0].Operator != "") {
			return dvgrammar.FLOW_RETURN, nil, errors.New("'continue' has no parameters")
		}
	}
	return dvgrammar.FLOW_CONTINUE, nil, nil
}
