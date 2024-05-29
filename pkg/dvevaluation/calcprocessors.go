/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
)

func ProcessorPlus(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res interface{}
	var resStr string
	var resNum float64 = 0
	dataType := dvgrammar.TYPE_NUMBER
	for i := 0; i < l; i++ {
		dtp := dvgrammar.TYPE_NULL
		v := values[i]
		if v != nil {
			dtp = v.DataType
		}
		if dtp != dvgrammar.TYPE_NUMBER && dtp != dvgrammar.TYPE_NUMBER_INT && dtp != dvgrammar.TYPE_NAN {
			if i == 0 {
				resStr = ""
			} else {
				resStr = AnyToString(resNum)
			}
			for ; i < l; i++ {
				resStr += AnyToString(values[i])
			}
			res = resStr
			dataType = dvgrammar.TYPE_STRING
			break
		}
		resNum += AnyToNumber(v)
	}
	if dataType == dvgrammar.TYPE_NUMBER {
		if math.IsNaN(resNum) {
			dataType = dvgrammar.TYPE_NAN
		}
		res = resNum
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorMinus(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res float64 = 0
	dataType := dvgrammar.TYPE_NUMBER
	for i := 0; i < l; i++ {
		if i == 0 {
			res = AnyToNumber(values[i])
		} else {
			res -= AnyToNumber(values[i])
		}
	}
	if math.IsNaN(res) {
		dataType = dvgrammar.TYPE_NAN
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorMultiply(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res float64 = 0
	dataType := dvgrammar.TYPE_NUMBER
	for i := 0; i < l; i++ {
		v := values[i]
		if v == nil {
			res = 0
		} else {
			if i == 0 {
				res = AnyToNumber(v.Value)
			} else {
				res *= AnyToNumber(v.Value)
			}
		}
	}
	if math.IsNaN(res) {
		dataType = dvgrammar.TYPE_NAN
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorDivision(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res float64 = 0
	dataType := dvgrammar.TYPE_NUMBER
	for i := 0; i < l; i++ {
		if i == 0 {
			res = AnyToNumber(values[i])
		} else {
			res /= AnyToNumber(values[i])
		}
	}
	if math.IsNaN(res) {
		dataType = dvgrammar.TYPE_NAN
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorPower(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res float64 = 0
	dataType := dvgrammar.TYPE_NUMBER
	for i := 0; i < l; i++ {
		if i == 0 {
			res = AnyToNumber(values[i])
		} else {
			res = math.Pow(res, AnyToNumber(values[i]))
		}
	}
	if math.IsNaN(res) {
		dataType = dvgrammar.TYPE_NAN
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorPercent(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	ok := false
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i])
			if !ok {
				break
			}
		} else {
			f, ok = AnyToNumberInt(values[i])
			if !ok || f == 0 {
				ok = false
				break
			}
			res = res % f
		}
	}
	var resFinal interface{} = res
	if !ok {
		dataType = dvgrammar.TYPE_NAN
		resFinal = math.NaN()
	}
	return &dvgrammar.ExpressionValue{Value: resFinal, DataType: dataType}, nil
}

func ProcessorLeftShift(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			if values[i] != nil {
				res, ok = AnyToNumberInt(values[i])
				if !ok {
					res = 0
				}
			}
		} else {
			if values[i] != nil {
				f, ok = AnyToNumberInt(values[i])
				if !ok {
					f = 0
				}
				if f < 0 || f > 64 {
					res = 0
				} else {
					res = res << uint(f)
				}
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorContainsIn(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res bool = false
	dataType := dvgrammar.TYPE_BOOLEAN
	if l >= 2 && values[0] != nil && values[1] != nil {
		res = ContainInProcess(values[0].Value, values[1].Value)
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorColon(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(tree.Children)
	if l != 2 {
		return nil, errors.New("Insufficient arguments for ternary operators")
	}
	condNode, err := GetLeftestQuestionNode(tree)
	if err != nil {
		return nil, err
	}
	_, condition, err := condNode.ExecuteExpression(context)
	if err != nil {
		return nil, err
	}
	right := AnyToBoolean(condition)
	if right {
		_, val, err := tree.Children[0].ExecuteExpression(context)
		return val, err
	}
	_, val, err := tree.Children[1].ExecuteExpression(context)
	return val, err
}

func GetLeftestQuestionNode(tree *dvgrammar.BuildNode) (*dvgrammar.BuildNode, error) {
	if tree == nil || len(tree.Children) != 2 || tree.Children[0] == nil {
		return nil, errors.New("No appropriate ? clause for :")
	}
	t := tree.Children[0]
	operator := t.Operator
	if operator == ":" {
		return GetLeftestQuestionNode(tree.Children[0])
	}
	if operator != "?" {
		return nil, errors.New("No relevant ? clause for :")
	}
	for len(t.Children) == 2 && t.Children[0] != nil && (t.Children[0].Operator == "?" || t.Children[0].Operator == ":") {
		if t.Children[0].Operator == ":" {
			return GetLeftestQuestionNode(t.Children[0])
		}
		t = t.Children[0]
	}
	if len(t.Children) != 2 {
		return nil, errors.New("Not enough clauses for ? (question) in the ternary operator")
	}
	res := t.Children[0]
	t.CloneFrom(t.Children[1])
	return res, nil
}

func ProcessorQuestion(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(tree.Children)
	if l != 2 {
		return nil, errors.New("Ternary operator requires condition ? value1 : value2 format")
	}
	node1, node2, ok := GetColumnSubNodes(tree.Children[1])
	if !ok {
		return nil, errors.New("Missing : clause in ternary operation")
	}
	return CalculateTernaryOperator(tree.Children[0], node1, node2, context)
}

func CalculateTernaryOperator(condNode, node1, node2 *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	_, condition, err := condNode.ExecuteExpression(context)
	if err != nil {
		return nil, err
	}
	right := AnyToBoolean(condition)
	if right {
		_, val, err := node1.ExecuteExpression(context)
		return val, err
	}
	_, val, err := node2.ExecuteExpression(context)
	return val, err
}

func ProcessorRightShift(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				f = 0
			}
			if f < 0 {
				res = 0
			} else {
				res = res >> uint(f)
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorLogicalRightShift(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				f = 0
			}
			if f < 0 || f >= 64 {
				if res < 0 {
					res = -1
				} else {
					res = 0
				}
			} else {
				res = int64(uint64(res) >> uint(f))
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorBoolAnd(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				f = 0
			}
			if f < 0 {
				res = 0
			} else {
				res = res & f
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorBoolOr(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				f = 0
			}
			if f < 0 {
				res = 0
			} else {
				res = res | f
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorBoolXor(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	var res int64 = 0
	var f int64
	var ok bool
	dataType := dvgrammar.TYPE_NUMBER_INT
	for i := 0; i < l; i++ {
		if i == 0 {
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				f = 0
			}
			if f < 0 {
				res = 0
			} else {
				res = res ^ f
			}
		}
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
}

func ProcessorEqualExact(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := values[0].DataType == values[1].DataType
	if res && AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) != 0 {
		res = false
	}
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorNotEqualExact(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := values[0].DataType == values[1].DataType
	if res && AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) != 0 {
		res = false
	}
	return &dvgrammar.ExpressionValue{Value: !res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorNotEqual(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
        res:=false
        if values[0]==nil {
           res  = values[1]==nil
        } else if values[1]!=nil {
	   res = AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) != 0
        }
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorEqual(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	leftValue := values[0]
	rightValue := values[1]
	if leftValue == nil {
		leftValue = &dvgrammar.ExpressionValue{Value: nil, DataType: dvgrammar.TYPE_NULL}
	}
	if rightValue == nil {
		rightValue = &dvgrammar.ExpressionValue{Value: nil, DataType: dvgrammar.TYPE_NULL}
	}
	res := AnyCompareAnyWithTypes(leftValue.DataType, leftValue.Value, rightValue.DataType, rightValue.Value) == 0
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorGreaterThan(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) == 1
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorLessThan(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) == -1
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorLessEqual(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value)
	return &dvgrammar.ExpressionValue{Value: res == 0 || res == -1, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorGreaterEqual(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value)
	return &dvgrammar.ExpressionValue{Value: res == 0 || res == 1, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorInInsideFor(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	return nil, errors.New("'in' should be used only inside 'for' declaration")
}

func ProcessorOfInsideFor(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	return nil, errors.New("'of' should be used only inside 'for' declaration")
}

func ProcessorElseInsideIf(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	return nil, errors.New("'else' should be used only inside 'if' declaration")
}
