/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
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
	return nil, nil
}

func String_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(s)
	return n, nil
}
