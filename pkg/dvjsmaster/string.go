/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strings"
)

func string_init() {
	dvevaluation.StringMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("contains"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_contains,
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
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        String_length,
					Immediate: true,
				},
			},
			{
				Name: []byte("endsWith"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_endsWith,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func String_contains(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
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
