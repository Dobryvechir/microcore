/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvcrypt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strings"
)

func string_init() {
	dvevaluation.StringMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("charAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_charAt,
				},
			},
			{
				Name: []byte("charCodeAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("codePointAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("concat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("endsWith"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_endsWith,
				},
			},
			{
				Name: []byte("fromCharCode"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("fromCodePoint"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("includes"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("indexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("lastIndexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("localeCompare"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("match"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("matchAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("normalize"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("padEnd"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("padStart"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("raw"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("repeat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("replace"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("replaceAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("search"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("slice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("split"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_split,
				},
			},
			{
				Name: []byte("startsWith"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("substring"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLocaleLowerCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLocaleUpperCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLowerCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toUpperCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trim"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trimEnd"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trimStart"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("valueOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        String_length,
					Immediate: true,
				},
			},
			{
				Name: []byte("isValidUUID"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_isValidUUID,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func String_includes(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	b := strings.Contains(s, s1)
	return b, nil
}

func String_charAt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	m := len(s)
	if m == 0 {
		return "", nil
	}
	n := len(params)
	p := 0
	if n != 0 && params[0] != nil {
		p64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok {
			if p64 > int64(m) || p64 < 0 {
				return "", nil
			}
			p = int(p64)
		}
	}
	return s[p : p+1], nil
}

func String_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(s)
	return n, nil
}

func String_endsWith(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			m := int(p)
			s = s[:m]
		}
	}
	b := strings.HasSuffix(s, s1)
	return b, nil
}

func String_startsWith(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			m := int(p)
			s = s[:m]
		}
	}
	b := strings.HasPrefix(s, s1)
	return b, nil
}

func String_split(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	limit := 0
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			limit = int(p)
		}
	}
	var b []string
	if limit > 0 {
		b = strings.SplitN(s, s1, limit)
	} else {
		b = strings.Split(s, s1)
	}
	return b, nil
}

func String_isValidUUID(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	b := dvcrypt.IsValidUUID(s)
	return b, nil
}
