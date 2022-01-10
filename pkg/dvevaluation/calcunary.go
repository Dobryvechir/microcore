/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
)

func BitwiseNotOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v, ok := AnyToNumberInt(value.Value)
	if ok {
		value = &dvgrammar.ExpressionValue{Value: -1 ^ v, DataType: dvgrammar.TYPE_NUMBER_INT}
	} else {
		value = &dvgrammar.ExpressionValue{Value: math.NaN(), DataType: dvgrammar.TYPE_NUMBER}
	}
	return value, nil
}

func UnaryPlusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v := AnyToNumber(value.Value)
	value = &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_NUMBER}
	return value, nil
}

func UnaryMinusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v := -AnyToNumber(value.Value)
	value = &dvgrammar.ExpressionValue{Value: v, DataType: dvgrammar.TYPE_NUMBER}
	return value, nil
}

func PrePlusPlusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v, ok := AnyToNumberInt(value.Value)
	var res interface{}
	var dataType int
	if ok {
		res = v + 1
		dataType = dvgrammar.TYPE_NUMBER_INT
	} else {
		res = math.NaN()
		dataType = dvgrammar.TYPE_NUMBER
	}
	value = &dvgrammar.ExpressionValue{Value: res, DataType: dataType}
	err := SetNodeValue(tree, res, dataType, context, lastVarName, lastParent)
	return value, err
}

func PreMinusMinusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v, ok := AnyToNumberInt(value.Value)
	var res interface{}
	var dataType int
	if ok {
		res = v - 1
		dataType = dvgrammar.TYPE_NUMBER_INT
	} else {
		res = math.NaN()
		dataType = dvgrammar.TYPE_NUMBER
	}
	value = &dvgrammar.ExpressionValue{Value: res, DataType: dataType}
	err := SetNodeValue(tree, res, dataType, context, lastVarName, lastParent)
	return value, err
}

func PostPlusPlusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v, ok := AnyToNumberInt(value.Value)
	var res interface{}
	var dataType int
	if ok {
		res = v + 1
		dataType = dvgrammar.TYPE_NUMBER_INT
		value = &dvgrammar.ExpressionValue{DataType: dataType, Value: v}
	} else {
		res = math.NaN()
		dataType = dvgrammar.TYPE_NUMBER
		value = &dvgrammar.ExpressionValue{DataType: dataType, Value: res}
	}
	err := SetNodeValue(tree, res, dataType, context, lastVarName, lastParent)
	return value, err
}

func PostMinusMinusOperator(value *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string, lastVarName string, lastParent *dvgrammar.ExpressionValue) (*dvgrammar.ExpressionValue, error) {
	v, ok := AnyToNumberInt(value.Value)
	var res interface{}
	var dataType int
	if ok {
		res = v - 1
		dataType = dvgrammar.TYPE_NUMBER_INT
		value = &dvgrammar.ExpressionValue{DataType: dataType, Value: v}
	} else {
		res = math.NaN()
		dataType = dvgrammar.TYPE_NUMBER
		value = &dvgrammar.ExpressionValue{DataType: dataType, Value: res}
	}
	err := SetNodeValue(tree, res, dataType, context, lastVarName, lastParent)
	return value, err
}
