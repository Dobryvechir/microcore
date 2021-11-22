/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
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
		Fields: make(map[string]*dvevaluation.DvVariable),
		Kind:   dvevaluation.FIELD_OBJECT,
		Prototype: &dvevaluation.DvVariable{
			Fields: map[string]*dvevaluation.DvVariable{
				"compareVersions": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Math_CompareVersions,
				},
				"increaseVersion": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Math_IncreaseVersion,
				},
			},
			Kind: dvevaluation.FIELD_OBJECT,
		},
	})
}

func Math_CompareVersions(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable,
	params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	n := len(params)
	s1 := ""
	s2 := ""
	defVersion := ""
	if n >= 1 {
		s1 = params[0].Value
	}
	if n >= 2 {
		s2 = params[1].Value
	}
	if n >= 3 {
		defVersion = params[2].Value
	}
	comp := MathCompareVersions(s1, s2, defVersion)
	res := strconv.Itoa(comp)
	return &dvevaluation.DvVariable{Value: res, Kind: dvevaluation.FIELD_NUMBER}, nil
}

func Math_IncreaseVersion(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	n := len(params)
	s := ""
	limit := 0
	defVersion := ""
	if n >= 1 {
		s = params[0].Value
	}
	if n >= 2 {
		lim, err := strconv.Atoi(params[1].Value)
		if err != nil && lim > 0 {
			limit = lim
		}
	}
	if n >= 3 {
		defVersion = params[2].Value
	}
	version := MathIncreaseVersion(s, limit, defVersion)
	return &dvevaluation.DvVariable{Value: version, Kind: dvevaluation.FIELD_STRING}, nil
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
