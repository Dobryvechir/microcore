/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strings"
)

func regexp_init() {
	dvevaluation.RegExpMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("dotAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_dotAll,
					Immediate: true,
				},
			},
			{
				Name: []byte("flags"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_flags,
					Immediate: true,
				},
			},
			{
				Name: []byte("global"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_global,
					Immediate: true,
				},
			},
			{
				Name: []byte("hasIndices"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_hasIndices,
					Immediate: true,
				},
			},
			{
				Name: []byte("ignoreCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_ignoreCase,
					Immediate: true,
				},
			},
			{
				Name: []byte("lastIndex"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_global,
					Immediate: true,
				},
			},
			{
				Name: []byte("multiline"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_multiline,
					Immediate: true,
				},
			},
			{
				Name: []byte("source"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_source,
					Immediate: true,
				},
			},
			{
				Name: []byte("sticky"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_sticky,
					Immediate: true,
				},
			},
			{
				Name: []byte("unicode"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        RegExp_unicode,
					Immediate: true,
				},
			},
			{
				Name: []byte("@@match"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("@@matchAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("@@replace"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("@@search"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("@@split"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("exec"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
				},
			},
			{
				Name: []byte("test"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_test,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: RegExp_global,
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
	patternStr := dvevaluation.AnyToString(pattern)
	flagsStr := dvevaluation.AnyToString(flags)
	patternAlt := getRegExpression(pattern)
	if patternAlt != nil {
		patternStr = patternAlt.Pattern
	}
	if len(patternStr) == 0 {
		patternStr = "(?:)"
	}
	rex, err := dvtextutils.NewRegExpression(patternStr, flagsStr)
	if err != nil {
		return nil, err
	}
	v.Extra = rex
	v.Fields = []*dvevaluation.DvVariable{
		{
			Kind:  dvevaluation.FIELD_NUMBER,
			Name:  []byte("lastIndex"),
			Value: []byte("0"),
		},
	}
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
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	return rex.Flags, nil
}

func RegExp_dotAll(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "s")
	return v, nil
}

func RegExp_global(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "g")
	return v, nil
}

func RegExp_hasIndices(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "d")
	return v, nil
}

func RegExp_ignoreCase(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "i")
	return v, nil
}

func RegExp_multiline(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "m")
	return v, nil
}

func RegExp_sticky(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "y")
	return v, nil
}

func RegExp_source(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	return rex.Pattern, nil
}

func RegExp_unicode(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	v := strings.Contains(rex.Flags, "u")
	return v, nil
}

func RegExp_test(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	rex := getRegExpression(thisVariable)
	if rex == nil {
		return "", nil
	}
	s := ""
	if len(params) > 0 {
		s = dvevaluation.AnyToString(params[0])
	}
	res := false
	if rex.Compiled.MatchString(s) {
		res = true
	}
	return res, nil
}
