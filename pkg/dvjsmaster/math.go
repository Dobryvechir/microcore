/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvcrypt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
	"strconv"
	"strings"
)

var MathMaster *dvevaluation.DvVariable

const (
	VERSION_DEFAULT = "0.0.0.0"
	VERSION_LIMIT   = 100
)

func math_init() {
	MathMaster = dvevaluation.RegisterMasterVariable("Math", &dvevaluation.DvVariable{
		Fields: make([]*dvevaluation.DvVariable, 0, 7),
		Kind:   dvevaluation.FIELD_OBJECT,
		Prototype: &dvevaluation.DvVariable{
			Fields: []*dvevaluation.DvVariable{
				{
					Name: []byte("compareVersions"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_CompareVersions,
					},
				},
				{
					Name: []byte("increaseVersion"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_IncreaseVersion,
					},
				},
				{
					Name: []byte("generateUUID"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_GenerateUUID,
					},
				},
				{
					Name: []byte("validUUID"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_ValidUUID,
					},
				},
				{
					Name: []byte("abs"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Abs,
					},
				},
				{
					Name: []byte("acos"),
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Acos,
					},
				},
				{
					Name: []byte("acosh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Acosh,
					},
				},
		                {
					Name: []byte("asin"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Asin,
					},
				},
		                {
					Name: []byte("asinh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Asinh,
					},
				},
                                {
					Name: []byte("atan"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Atan,
					},
				},
                                {
					Name: []byte("atanh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Atanh,
					},
				},
                                {
					Name: []byte("cbrt"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Cbrt,
					},
				},
                                {
					Name: []byte("ceil"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Ceil,
					},
				},
				{
					Name: []byte("clz32"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Clz32,
					},
				},                                 
                                {
					Name: []byte("cos"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Cos,
					},
				},
                                {
					Name: []byte("cosh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Cosh,
					},
				},
                                {
					Name: []byte("exp"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Exp,
					},
				},
                                {
					Name: []byte("expm1"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Expm1,
					},
				},
                                {
					Name: []byte("floor"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Floor,
					},
				},
                                {
					Name: []byte("log"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Log,
					},
				},
                                {
					Name: []byte("log1p"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Log1p,
					},
				},
                                {
					Name: []byte("log10"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Log10,
					},
				},
                                {
					Name: []byte("log2"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Log2,
					},
				},
                                {
					Name: []byte("round"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Round,
					},
				},  
                                {
					Name: []byte("sign"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Sign,
					},
				},
                                {
					Name: []byte("sin"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Sin,
					},
				},
                                {
					Name: []byte("sinh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Sinh,
					},
				},
                                {
					Name: []byte("sqrt"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Sqrt,
					},
				},
                                {
					Name: []byte("tan"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Tan,
					},
				},
                                {
					Name: []byte("tanh"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Tanh,
					},
				},
                                {
					Name: []byte("trunc"),                      
					Kind: dvevaluation.FIELD_FUNCTION,
					Extra: &dvevaluation.DvFunction{
						Fn: Math_Trunc,
					},
				},
                                {
					Name: []byte("E"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("2.718281828459045"),
				},
                                {
					Name: []byte("LN10"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("2.302585092994046"),
				},
                                {
					Name: []byte("LN2"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("0.6931471805599453"),
				},
                                {
					Name: []byte("LOG10E"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("0.4342944819032518"),
				},
                                {
					Name: []byte("LOG2E"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("1.4426950408889634"),
				},
                                {
					Name: []byte("PI"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("3.14159265359"),
				},
                                {
					Name: []byte("SQRT1_2"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("0.7071067811865476"),
				},
                                {
					Name: []byte("SQRT2"),                      
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte("1.4142135623730951"),
				},


			},
			Kind: dvevaluation.FIELD_OBJECT,
		},
	})
}

func Math_CompareVersions(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	s1 := ""
	s2 := ""
	defVersion := ""
	if n >= 1 {
		s1 = dvevaluation.AnyToString(params[0])
	}
	if n >= 2 {
		s2 = dvevaluation.AnyToString(params[1])
	}
	if n >= 3 {
		defVersion = dvevaluation.AnyToString(params[2])
	}
	comp := MathCompareVersions(s1, s2, defVersion)
	res := strconv.Itoa(comp)
	return res, nil
}

func Math_IncreaseVersion(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	s := ""
	limit := 0
	defVersion := ""
	if n >= 1 {
		s = dvevaluation.AnyToString(params[0])
	}
	if n >= 2 {
		lim, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && lim > 0 {
			limit = int(lim)
		}
	}
	if n >= 3 {
		defVersion = dvevaluation.AnyToString(params[2])
	}
	version := MathIncreaseVersion(s, limit, defVersion)
	return version, nil
}

func MathSplitVersion(s string, defVersion string) []string {
	if s == "" || !(s[0] >= '0' && s[0] <= '9') {
		if defVersion == "" || !(defVersion[0] >= '0' && defVersion[0] <= '9') {
			defVersion = VERSION_DEFAULT
		}
		return MathSplitVersion(defVersion, defVersion)
	}
	i := 1
	n := len(s)
	m := 1
	for ; i < n; i++ {
		c := s[i]
		if c == '.' && i+1 < n && s[i+1] >= '0' && s[i+1] <= '9' {
			m++
			i++
		} else if !(c >= '0' && c <= '9') {
			break
		}
	}
	r := make([]string, m)
	n = i
	p := 0
	m = 0
	for i = 0; i < n; i++ {
		if s[i] == '.' {
			r[m] = s[p:i]
			p = i + 1
			m++
		}
	}
	r[m] = s[p:]
	return r
}

func MathCompareVersions(s1 string, s2 string, defVersion string) int {
	v1 := MathSplitVersion(s1, defVersion)
	v2 := MathSplitVersion(s2, defVersion)
	n1 := len(v1)
	n2 := len(v2)
	mn := n1
	if n2 < mn {
		mn = n2
	}
	for i := 0; i < mn; i++ {
		k1, _ := strconv.Atoi(v1[i])
		k2, _ := strconv.Atoi(v2[i])
		dif := k1 - k2
		if dif != 0 {
			if dif > 0 {
				dif = 1
			} else {
				dif = -1
			}
			return dif
		}
	}
	if n1 > mn {
		for i := mn + 1; i < n1; i++ {
			k, _ := strconv.Atoi(v1[i])
			if k != 0 {
				return 1
			}
		}
	} else if n2 > mn {
		for i := mn + 1; i < n2; i++ {
			k, _ := strconv.Atoi(v2[i])
			if k != 0 {
				return -1
			}
		}
	}
	return 0
}

func MathJoinVersion(v []string) string {
	return strings.Join(v, ".")
}

func MathIncreaseVersion(s string, limit int, defVersion string) string {
	if limit <= 0 {
		limit = VERSION_LIMIT
	}
	v := MathSplitVersion(s, defVersion)
	n := len(v)
	for i := n - 1; i >= 0; i-- {
		k, _ := strconv.Atoi(v[i])
		if k < limit || i == 0 {
			v[i] = strconv.Itoa(k + 1)
			break
		} else {
			v[i] = "0"
		}
	}
	return MathJoinVersion(v)
}

func Math_GenerateUUID(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	uuid := dvcrypt.GetRandomUuid()
	return uuid, nil
}

func Math_ValidUUID(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return false, nil
	}
	uuid := dvevaluation.AnyToString(params[0])
	res := dvcrypt.IsValidUUID(uuid)
	return res, nil
}

func Math_Abs(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Abs(val)
	return res, nil
}

func Math_Acos(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Acos(val)
	return res, nil
}

func Math_Acosh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Acosh(val)
	return res, nil
}

func Math_Asin(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Asin(val)
	return res, nil
}

func Math_Asinh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Asinh(val)
	return res, nil
}

func Math_Atan(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Atan(val)
	return res, nil
}

func Math_Atanh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Atanh(val)
	return res, nil
}

func Math_Cbrt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Cbrt(val)
	return res, nil
}

func Math_Ceil(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Ceil(val)
	return res, nil
}

func Math_Cos(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Cos(val)
	return res, nil
}

func Math_Cosh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Cosh(val)
	return res, nil
}

func Math_Exp(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Exp(val)
	return res, nil
}

func Math_Expm1(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Expm1(val)
	return res, nil
}

func Math_Floor(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Floor(val)
	return res, nil
}

func Math_Log(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Log(val)
	return res, nil
}

func Math_Log1p(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Log1p(val)
	return res, nil
}

func Math_Log10(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Log10(val)
	return res, nil
}

func Math_Log2(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Log2(val)
	return res, nil
}

func Math_Round(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Round(val)
	return res, nil
}

func Math_Sign(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:= 0
        if val < 0 {
         res = -1
        } else if val > 0 {
           res = 1
        } else if val == -0 {
           res = -0
        }
	return res, nil
}

func Math_Sin(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Sin(val)
	return res, nil
}

func Math_Sinh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Sinh(val)
	return res, nil
}

func Math_Sqrt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Sqrt(val)
	return res, nil
}

func Math_Tan(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Tan(val)
	return res, nil
}

func Math_Tanh(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Tanh(val)
	return res, nil
}

func Math_Trunc(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return 0, nil
	}
	val := dvevaluation.AnyToNumber(params[0])
	res:=math.Trunc(val)
	return res, nil
}
func Math_Clz32(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
        n := len(params)
	if n == 0 {
		return 0, nil
	}
	m, ok := dvevaluation.AnyToNumberInt(params[0])
	if !ok {
		return 0, nil
	}
	
	t := 0
        for i := 31; i>=0; i-- {
	    if m&(1<<i) != 0 {
	       break
	    }
	    t++
        }

        return t, nil
}