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
