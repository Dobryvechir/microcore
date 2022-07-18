/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strconv"
)

var ArrayIteratorPrototype = &dvevaluation.DvVariable {
	Fields: []*dvevaluation.DvVariable{
		{
			Name: []byte("next"),
			Kind: dvevaluation.FIELD_FUNCTION,
			Extra: &dvevaluation.DvFunction{
				Fn: Array_iteratorNext,
			},
		},
	},
	Kind: dvevaluation.FIELD_OBJECT,
}

func Array_iteratorNextValue(value *dvevaluation.DvVariable,done bool) (interface{}, error) {
	var newValue *dvevaluation.DvVariable
	if value == nil {
		newValue = &dvevaluation.DvVariable{
			Kind: dvevaluation.FIELD_NULL,
		}
	} else {
		newValue=value.Clone()
	}
	newValue.Name = []byte("value")
	return &dvevaluation.DvVariable{
		Kind: dvevaluation.FIELD_OBJECT,
		Fields: []*dvevaluation.DvVariable{
			newValue,
			{
				Kind: dvevaluation.FIELD_BOOLEAN,
				Value: dvevaluation.GetBooleanValue(done),
				Name: []byte("done"),
			},
		},
	}, nil
}

func Array_iteratorNextFinish() (interface{}, error) {
	return Array_iteratorNextValue(nil, true)
}

func Array_iteratorNext(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v:=dvevaluation.AnyToDvVariable(thisVariable)
	if v==nil || v.Extra ==nil {
		return Array_iteratorNextFinish()
	}
	num, ok:=v.Extra.(int)
	if !ok || num<0 {
		return Array_iteratorNextFinish()
	}
	v.Extra = num + 1
	if v.Kind==dvevaluation.FIELD_ARRAY || v.Kind==dvevaluation.FIELD_OBJECT {
		if num < len(v.Fields) {
			return Array_iteratorNextValue(v.Fields[num], false)
		}
	} else if v.Kind==dvevaluation.FIELD_STRING {
		if num < len(v.Value) {
			return Array_iteratorNextValue(&dvevaluation.DvVariable{
				Kind: dvevaluation.FIELD_STRING,
				Value: []byte{v.Value[num]},
			}, false)
		}
	}
    return Array_iteratorNextFinish()
}

func Array_entries(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return createArrayIterator(res), nil
	}
	n:=len(v.Fields)
	res.Fields = make([]*dvevaluation.DvVariable, n)
	for i:=0;i<n;i++ {
		res.Fields[i] = &dvevaluation.DvVariable{
			Kind: dvevaluation.FIELD_ARRAY,
			Fields: []*dvevaluation.DvVariable{
				{
					Kind: dvevaluation.FIELD_NUMBER,
					Value: []byte(strconv.Itoa(i)),
				},
				v.Fields[i],
			},
		}
	}
	return createArrayIterator(res), nil
}

func Array_keys(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return createArrayIterator(res), nil
	}
	n:=len(v.Fields)
	res.Fields = make([]*dvevaluation.DvVariable,n)
	for i:=0;i<n;i++ {
		res.Fields[i] = &dvevaluation.DvVariable{
			Kind: dvevaluation.FIELD_NUMBER,
			Value: []byte(strconv.Itoa(i)),
		}
	}
	return createArrayIterator(res), nil
}

func Array_values(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return createArrayIterator(res), nil
	}
	n:=len(v.Fields)
	res.Fields = make([]*dvevaluation.DvVariable, n)
	for i:=0;i<n;i++ {
		res.Fields[i] = v.Fields[i]
	}
	return createArrayIterator(res), nil
}

func createArrayIterator(r *dvevaluation.DvVariable) *dvevaluation.DvVariable {
	r.Extra = 0
	r.Prototype = ArrayIteratorPrototype
	return r
}
