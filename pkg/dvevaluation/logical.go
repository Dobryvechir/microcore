/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func LogicalDataGetter(token *dvgrammar.Token, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	if token.DataType == dvgrammar.TYPE_DATA {
		if newType, ok := reservedWords[token.Value]; ok {
			return &dvgrammar.ExpressionValue{Value: buildinTypes[token.Value], DataType: newType}, nil
		} else {
			_, v := context.Scope.Get(token.Value)
			if (context.VisitorOptions & dvgrammar.EVALUATE_OPTION_UNDEFINED) != 0 {
				v = !v
			}
			return &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_BOOLEAN}, nil
		}
	}
	return &dvgrammar.ExpressionValue{Value: token.Value, DataType: token.DataType}, nil
}

func LogicalNotOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v := !AnyToBoolean(value.Value)
	value = &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_BOOLEAN}
	return value, nil
}

func ProcessorBooleanOr(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (res *dvgrammar.ExpressionValue, err error) {
	l := tree.GetChildrenNumber()
	for i := 0; i < l; i++ {
		res, err = tree.GetChildrenExpressionValue(i, context)
		if err != nil {
			return nil, err
		}
		if AnyToBoolean(res) {
			break
		}
	}
	return
}

func ProcessorBooleanAnd(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (res *dvgrammar.ExpressionValue, err error) {
	l := tree.GetChildrenNumber()
	for i := 0; i < l; i++ {
		res, err = tree.GetChildrenExpressionValue(i, context)
		if err != nil {
			return nil, err
		}
		if !AnyToBoolean(res) {
			break
		}
	}
	return
}

func ProcessorBooleanOrNullable(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (res *dvgrammar.ExpressionValue, err error) {
	l := tree.GetChildrenNumber()
	for i := 0; i < l; i++ {
		res, err = tree.GetChildrenExpressionValue(i, context)
		if err != nil {
			return nil, err
		}
		if IsNotNullish(res) {
			break
		}
	}
	return
}

var LogicalOperators = map[string]dvgrammar.InterOperatorVisitor{
	"||": ProcessorBooleanOr,
	"&&": ProcessorBooleanAnd,
}

func CalculateDefined(data []byte, scope dvgrammar.ScopeInterface, reference *dvgrammar.SourceReference, visitorOptions int) (*dvgrammar.ExpressionValue, error) {
	context := &dvgrammar.ExpressionContext{
		Scope:          scope,
		Reference:      reference,
		Rules:          LogicalRules,
		VisitorOptions: visitorOptions,
	}
	return dvgrammar.FastEvaluation(data, context)
}

var LogicalUnaryMap = map[string]dvgrammar.UnaryVisitor{
	"!": LogicalNotOperator,
}
var LogicalRules = &dvgrammar.GrammarRuleDefinitions{
	BaseGrammar:      dvgrammar.LogicalGrammarBaseDefinition,
	Visitors:         LogicalOperators,
	DataGetter:       LogicalDataGetter,
	EvaluateOptions:  0,
	UnaryPreVisitors: LogicalUnaryMap,
}
