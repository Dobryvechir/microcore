/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func CalculatorDataGetter(token *dvgrammar.Token, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	name := ""
	if token.DataType == dvgrammar.TYPE_DATA {
		if newType, ok := reservedWords[token.Value]; ok {
			return &dvgrammar.ExpressionValue{Value: buildinTypes[token.Value], DataType: newType}, nil
		} else {
			name = token.Value
			v, ok := context.Scope.Get(name)
			if !ok {
				return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NULL, Name: name}, errors.New(token.Value + " is not defined")
			}
			rv:=AnyToDvGrammarExpressionValue(v)
			if rv!=nil {
				rv.Name = name
			}
			return rv, nil
		}
	}
	var v interface{} = token.Value
	switch token.DataType {
	case dvgrammar.TYPE_NUMBER:
		v = AnyToNumber(token.Value)
	case dvgrammar.TYPE_NUMBER_INT:
		v, _ = AnyToNumberInt(token.Value)
	}
	return &dvgrammar.ExpressionValue{Value: v, DataType: token.DataType, Name: name}, nil
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
	"IN":  ProcessorContainsIn,
	":":   ProcessorColon,
	"?":   ProcessorQuestion,
	"=":   ProcessorAssign,
	"+=":  ProcessorPlusAssign,
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

var LanguageOperatorMap = map[string]dvgrammar.LanguageOperatorVisitor{
	"return":   ReturnOperator,
	"break":    ReturnOperator,
	"continue": ReturnOperator,
	"=>": ArrowFunctionOperator,
}

var CalculatorPostUnaryMap = map[string]dvgrammar.UnaryVisitor{
	"++": PostPlusPlusOperator,
	"--": PostMinusMinusOperator,
}

var CalculatorRules = &dvgrammar.GrammarRuleDefinitions{
	Visitors:          CalculatorOperators,
	BracketVisitor:    BracketProcessors,
	LanguageOperator:  LanguageOperatorMap,
	DataGetter:        CalculatorDataGetter,
	EvaluateOptions:   0,
	UnaryPreVisitors:  CalculatorUnaryMap,
	UnaryPostVisitors: CalculatorPostUnaryMap,
}

func SetNodeValue(tree *dvgrammar.BuildNode, v interface{}, dataType int, context *dvgrammar.ExpressionContext, lastVarName string, lastParent *dvgrammar.ExpressionValue) error {
	if lastVarName != "" {
		context.Scope.SetDeep(lastVarName, v)
	}
	return nil
}
