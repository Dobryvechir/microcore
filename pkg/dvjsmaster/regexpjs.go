/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func regexp_init() {
	dvevaluation.RegExpMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("dotAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_concat,
					Immediate: true,
				},
			},
			{
				Name: []byte("flags"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_flags,
					Immediate: true,
				},
			},
			{
				Name: []byte("global"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_entries,
					Immediate: true,
				},
			},
			{
				Name: []byte("hasIndices"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_every,
					Immediate: true,
				},
			},
			{
				Name: []byte("ignoreCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_fill,
					Immediate: true,
				},
			},
			{
				Name: []byte("lastIndex"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_filter,
					Immediate: true,
				},
			},
			{
				Name: []byte("multiline"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_find,
					Immediate: true,
				},
			},
			{
				Name: []byte("source"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findIndex,
					Immediate: true,
				},
			},
			{
				Name: []byte("sticky"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findLast,
					Immediate: true,
				},
			},
			{
				Name: []byte("unicode"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findLastIndex,
					Immediate: true,
				},
			},
			{
				Name: []byte("@@match"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("@@matchAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flatMap,
				},
			},
			{
				Name: []byte("@@replace"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_foreach,
				},
			},
			{
				Name: []byte("@@search"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_from,
				},
			},
			{
				Name: []byte("@@split"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_includes,
				},
			},
			{
				Name: []byte("exec"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_indexOf,
				},
			},
			{
				Name: []byte("test"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_isArray,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_join,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
	dvevaluation.RegExpMaster.Kind = dvevaluation.FIELD_FUNCTION
	dvevaluation.RegExpMaster.Extra = &dvevaluation.DvFunction{
		Fn: RegExp_constructor,
	}
}

func RegExp_constructor(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	var pattern interface{}
	var flags interface{}
	n := len(params)
	if n > 0 {
		pattern = params[0]
	}
	if n > 1 {
		flags = params[1]
	}
	v, err := regExpQuickCreation(thisVariable, pattern, flags)
	return v, err
}

func regExpQuickCreation(thisVar interface{}, pattern interface{}, flags interface{}) (*dvevaluation.DvVariable, error) {
	v := dvevaluation.AnyToDvVariable(thisVar)
	if v == nil || v.Kind != dvevaluation.FIELD_OBJECT || v.Fields != nil {
		v = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT}
	}
	v.Prototype = dvevaluation.RegExpMaster
	patternStr := ""
	flagsStr := dvevaluation.AnyToString(flags)
	rex, err := dvtextutils.NewRegExpression(patternStr, flagsStr)
	if err != nil {
		return nil, err
	}
	v.Extra = rex
	return v, nil
}

func getRegExpression(item interface{}) *dvtextutils.RegExpession {
	if item == nil {
		return nil
	}
	v := dvevaluation.AnyToDvVariable(item)
	if v == nil || v.Kind != dvevaluation.FIELD_OBJECT || v.Extra == nil {
		return nil
	}
	rex := v.Extra.(*dvtextutils.RegExpession)
	return rex
}

func RegExp_flags(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex:=getRegExpression(thisVariable)
	if rex==nil {
		return "", nil
	}
	return rex.Flags, nil
}
