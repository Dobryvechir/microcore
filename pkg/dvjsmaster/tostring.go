/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func Array_toString(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY {
		return dvevaluation.AnyToString(thisVariable), nil
	}
	res:=ArrayJoinWith(v, ",")
	return res, nil
}

func ArrayJoinWith(v *dvevaluation.DvVariable, joiner string) string {
	res := ""
	n := len(v.Fields)
	for i := 0; i < n; i++ {
		if i != 0 {
			res += joiner
		}
		res += dvevaluation.AnyToString(v.Fields[i])
	}
	return res
}

func ArrayJoinWithLocale(v *dvevaluation.DvVariable, locale string, options map[string]string) string {
	res := ""
	joiner:=","
	n := len(v.Fields)
	for i := 0; i < n; i++ {
		if i != 0 {
			res += joiner
		}
		res += ToStringByLocaleByKind(v.Fields[i], locale, options)
	}
	return res
}

func ToStringByLocaleByKind(v *dvevaluation.DvVariable, locale string, options map[string]string) string {
	//TODO there must be special implementation for numbers and dates
   res := dvevaluation.AnyToString(v)
   return res
}

func Array_toLocaleString(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY {
		return dvevaluation.AnyToString(thisVariable), nil
	}
	locale, options:=ToLocaleStringReadLocaleAndOptions(params)
	res:=ArrayJoinWithLocale(v, locale, options)
	return res, nil
}

func ToLocaleStringReadLocaleAndOptions(params []interface{}) (locale string, options map[string]string) {
	m:=len(params)
	if m>=1 {
		locale = dvevaluation.AnyToString(params[0])
	}
	locale = VerifyLocaleOrDefault(locale,"en-US")
	options =make(map[string]string)
	if m>=2 {
		v:=dvevaluation.AnyToDvVariable(params[1])
		if v!=nil {
			options = v.GetStringMap()
		}
	}
	return
}

func VerifyLocaleOrDefault(locale string, defLocale string) string {
    n:=len(locale)
	if n!=2 && n!=5 {
		return defLocale
	}
	if !(locale[0]>='a' && locale[0]<='z' && locale[1]>='a' && locale[1]<='z') {
		return defLocale
	}
	if n==2 {
		return locale
	}
	if locale[2]!='-' || !(locale[3]>='A' && locale[3]<='Z' && locale[4]>='A' && locale[4]<='Z') {
		return defLocale
	}
	return locale
}

func String_toString(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToString(thisVariable)
	return v, nil
}
