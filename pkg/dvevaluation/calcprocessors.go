/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
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
		if values[i].DataType == dvgrammar.TYPE_STRING {
			if i == 0 {
				resStr = ""
			} else {
				resStr = AnyToString(resNum)
			}
			for ; i < l; i++ {
				resStr += AnyToString(values[i].Value)
			}
			res = resStr
			dataType = dvgrammar.TYPE_STRING
			break
		}
		resNum += AnyToNumber(values[i].Value)
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
			res = AnyToNumber(values[i].Value)
		} else {
			res -= AnyToNumber(values[i].Value)
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
		if i == 0 {
			res = AnyToNumber(values[i].Value)
		} else {
			res *= AnyToNumber(values[i].Value)
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
			res = AnyToNumber(values[i].Value)
		} else {
			res /= AnyToNumber(values[i].Value)
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
			res = AnyToNumber(values[i].Value)
		} else {
			res = math.Pow(res, AnyToNumber(values[i].Value))
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
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				break
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
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
			res, ok = AnyToNumberInt(values[i].Value)
			if !ok {
				res = 0
			}
		} else {
			f, ok = AnyToNumberInt(values[i].Value)
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
	return &dvgrammar.ExpressionValue{Value: res, DataType: dataType}, nil
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
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) != 0
	return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_BOOLEAN}, nil
}

func ProcessorEqual(values []*dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, operator string) (*dvgrammar.ExpressionValue, error) {
	l := len(values)
	if l != 2 {
		return nil, errors.New("Only 2 parameters are allowed for " + operator)
	}
	res := AnyCompareAnyWithTypes(values[0].DataType, values[0].Value, values[1].DataType, values[1].Value) == 0
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
