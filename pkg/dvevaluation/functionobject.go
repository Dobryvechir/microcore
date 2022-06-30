/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import "github.com/Dobryvechir/microcore/pkg/dvgrammar"

var dvgrammarTypePrototyper map[int]*DvObject
var dvvariableTypePrototyper map[int]*DvObject

func (fn *DvFunctionObject) ExecuteDvFunctionWithTreeArguments(args []interface{}, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	vl, err := fn.Executor.Fn(context, fn.SelfRef, args)
	value := AnyToDvGrammarExpressionValue(vl)
	return value, err
}

func GetFunctionObjectVariable(fn *DvFunction, selfVal *dvgrammar.ExpressionValue, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	fnObj := &DvFunctionObject{SelfRef: selfVal, Context: context, Executor: fn}
	if fn.Immediate {
		vl, err := fnObj.ExecuteDvFunctionWithTreeArguments(nil, context)
		return vl, err
	}
	dv := &DvVariable{Kind: FIELD_FUNCTION, Extra: fnObj}
	return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_FUNCTION, Value: dv}, nil
}

func RefillDvGrammarTypePrototyper() {
	dvgrammarTypePrototyper = make(map[int]*DvObject)
	dvgrammarTypePrototyper[dvgrammar.TYPE_FUNCTION] = GetFunctionPrototypeFromMasterVariable("Function")
	dvgrammarTypePrototyper[dvgrammar.TYPE_STRING] = GetFunctionPrototypeFromMasterVariable("String")
	dvgrammarTypePrototyper[dvgrammar.TYPE_NUMBER] = GetFunctionPrototypeFromMasterVariable("Number")
	dvgrammarTypePrototyper[dvgrammar.TYPE_BOOLEAN] = GetFunctionPrototypeFromMasterVariable("Boolean")
	dvgrammarTypePrototyper[dvgrammar.TYPE_NAN] = dvgrammarTypePrototyper[dvgrammar.TYPE_NUMBER]
	dvgrammarTypePrototyper[dvgrammar.TYPE_NUMBER_INT] = dvgrammarTypePrototyper[dvgrammar.TYPE_NUMBER]
	dvgrammarTypePrototyper[dvgrammar.TYPE_CHAR] = dvgrammarTypePrototyper[dvgrammar.TYPE_STRING]
}

func RefillDvVariableTypePrototyper() {
	dvvariableTypePrototyper = make(map[int]*DvObject)
	dvvariableTypePrototyper[FIELD_FUNCTION] = GetFunctionPrototypeFromMasterVariable("Function")
	dvvariableTypePrototyper[FIELD_ARRAY] = GetFunctionPrototypeFromMasterVariable("Array")
	dvvariableTypePrototyper[FIELD_OBJECT] = GetFunctionPrototypeFromMasterVariable("Object")
	dvvariableTypePrototyper[FIELD_STRING] = GetFunctionPrototypeFromMasterVariable("String")
	dvvariableTypePrototyper[FIELD_NUMBER] = GetFunctionPrototypeFromMasterVariable("Number")
	dvvariableTypePrototyper[FIELD_BOOLEAN] = GetFunctionPrototypeFromMasterVariable("Boolean")
}

func (item *DvVariable) GetPrototypeByKind() *DvObject {
	if dvvariableTypePrototyper == nil {
		RefillDvVariableTypePrototyper()
	}
	if item == nil {
		return nil
	}
	return dvvariableTypePrototyper[item.Kind]
}

func GetPrototypeForDvGrammarExpressionValue(value *dvgrammar.ExpressionValue) *DvObject {
	if value == nil {
		return nil
	}
	if dvgrammarTypePrototyper == nil {
		RefillDvGrammarTypePrototyper()
	}
	d := dvgrammarTypePrototyper[value.DataType]
	if d != nil {
		return d
	}
	if value.DataType == dvgrammar.TYPE_OBJECT {
		v := AnyToDvVariable(value.Value)
		return v.GetPrototypeByKind()
	}
	return nil
}
