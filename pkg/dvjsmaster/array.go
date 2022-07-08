/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func array_init() {
	dvevaluation.ArrayMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("reduce"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_reduce,
				},
			},
			{
				Name: []byte("push"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_push,
				},
			},
			{
				Name: []byte("slice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        Array_length,
					Immediate: true,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Array_push(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil {
		return nil, errors.New("Cannot convert null to object")
	}
    n:=len(params)
	for i:=0;i<n;i++ {
		d:=dvevaluation.AnyToDvVariable(params[i])
		v.Fields = append(v.Fields, d)
	}
	n = len(v.Fields)
	return n, nil
}

func Array_slice(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	return nil, nil
}

func Array_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	n := 0
	if v != nil {
		n = len(v.Fields)
	}
	return n, nil
}

func Array_reduce(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var result interface{} = nil
	n := len(params)
	if n >= 2 {
		result = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		m:=len(v.Fields)
		for i := 0; i < m; i++ {
			fnParams := []interface{}{result, v.Fields[i], i, v}
			result, err = dvevaluation.ExecuteAnyFunction(context, fn, v, fnParams)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
