/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func array_init() {
	dvevaluation.ArrayMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("push"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra:   &dvevaluation.DvFunction {
					Fn: Array_push,
				},
			},
			{
				Name: []byte("slice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction {
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction {
					Fn: Array_length,
					Immediate: true,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Array_push(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	return nil, nil
}

func Array_slice(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	return nil, nil
}

func Array_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v:=dvevaluation.AnyToDvVariable(thisVariable)
	n:=0
	if v!=nil {
		n=len(v.Fields)
	}
	return n, nil
}
