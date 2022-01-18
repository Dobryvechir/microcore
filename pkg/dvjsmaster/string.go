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
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        String_length,
					Immediate: true,
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
