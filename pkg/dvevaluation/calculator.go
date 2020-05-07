/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func CalculatorDataGetter(token *dvgrammar.Token, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	if token.DataType == dvgrammar.TYPE_DATA {
		if newType, ok := reservedWords[token.Value]; ok {
			return &dvgrammar.ExpressionValue{Value: buildinTypes[token.Value], DataType: newType}, nil
		} else {
			v, ok := context.Scope.Get(token.Value)
			if !ok {
				return nil, errors.New(token.Value + " is not defined")
			}
			return &dvgrammar.ExpressionValue{Value: v, DataType: AnyGetType(v)}, nil
		}
	}
	return &dvgrammar.ExpressionValue{Value: token.Value, DataType: token.DataType}, nil
}

var CalculatorOperators = map[string]dvgrammar.InterOperatorVisitor{
	"+":   ProcessorPlus,
	"-":   ProcessorMinus,
	"*":   ProcessorMultiply,
	"/":   ProcessorDivision,
	"%":   ProcessorPercent,
	"&":   ProcessorBoolAnd,
	"|":   ProcessorBoolOr,
	"^":   ProcessorBoolXor,
	"**":  ProcessorPower,
	"||":  ProcessorBooleanOr,
	"&&":  ProcessorBooleanAnd,
	"<<":  ProcessorLeftShift,
	">>>": ProcessorLogicalRightShift,
	">>":  ProcessorRightShift,
	"===": ProcessorEqualExact,
	"!==": ProcessorNotEqualExact,
	"==":  ProcessorEqual,
	"!=":  ProcessorNotEqual,
	">":   ProcessorGreaterThan,
	">=":  ProcessorGreaterEqual,
	"<":   ProcessorLessThan,
	"<=":  ProcessorLessEqual,
}

func CalculatorEvaluator(data []byte, scope dvgrammar.ScopeInterface, reference *dvgrammar.SourceReference, visitorOptions int) (*dvgrammar.ExpressionValue, error) {
	context := &dvgrammar.ExpressionContext{
		Scope:          scope,
		Reference:      reference,
		Rules:          CalculatorRules,
		VisitorOptions: visitorOptions,
	}
	return dvgrammar.FastEvaluation(data, context)
}

var CalculatorUnaryMap = map[string]dvgrammar.UnaryVisitor{
	"!":  LogicalNotOperator,
	"~":  BitwiseNotOperator,
	"+":  UnaryPlusOperator,
	"-":  UnaryMinusOperator,
	"++": PrePlusPlusOperator,
	"--": PreMinusMinusOperator,
}

var CalculatorPostUnaryMap = map[string]dvgrammar.UnaryVisitor{
	"++": PostPlusPlusOperator,
	"--": PostMinusMinusOperator,
}

var CalculatorRules = &dvgrammar.GrammarRuleDefinitions{
	Visitors:          CalculatorOperators,
	DataGetter:        CalculatorDataGetter,
	EvaluateOptions:   0,
	UnaryPreVisitors:  CalculatorUnaryMap,
	UnaryPostVisitors: CalculatorPostUnaryMap,
}

func SetNodeValue(tree *dvgrammar.BuildNode, v interface{}, dataType int, context *dvgrammar.ExpressionContext) error {
	//TODO: implement
	return nil
}
