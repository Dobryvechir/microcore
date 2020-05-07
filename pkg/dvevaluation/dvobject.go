/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func NewDvObject(properties map[string]string, prototype *DvObject) *DvObject {
	props := make(map[string]interface{}, len(properties))
	for k, v := range properties {
		props[k] = StringToAny(v)
	}
	return NewDvObjectWithSpecialValues(nil, 0, prototype, props)
}

func NewDvObjectWithGlobalPrototype(properties map[string]string) *DvObject {
	return NewDvObject(properties, GlobalFunctionPrototype)
}

func NewObject(properties map[string]interface{}) *DvObject {
	return NewDvObjectWithSpecialValues(nil, 0, nil, properties)
}

func NewObjectWithPrototype(properties map[string]interface{}, prototype *DvObject) *DvObject {
	return NewDvObjectWithSpecialValues(nil, 0, prototype, properties)
}

func (obj *DvObject) Get(key string) (interface{}, bool) {
	if obj == nil {
		return nil, false
	}
	if obj.Properties != nil {
		res, ok := obj.Properties[key]
		if ok {
			return res, ok
		}
	}
	if obj.Prototype != nil {
		return obj.Prototype.Get(key)
	}
	return nil, false
}

func (obj *DvObject) GetString(key string) string {
	v, ok := obj.Get(key)
	if !ok || v == nil {
		return ""
	}
	return AnyToString(v)
}

func (obj *DvObject) Set(key string, value interface{}) {
	if obj == nil {
		return
	}
	if obj.Properties == nil {
		obj.Properties = make(map[string]interface{})
	}
	obj.Properties[key] = value
}

func NewDvObjectWithSpecialValues(value interface{}, kind int, proto *DvObject, properties map[string]interface{}) *DvObject {
	return &DvObject{Value: value, Options: kind, Prototype: proto, Properties: properties}
}

var DvObject_null *DvObject = NewDvObjectWithSpecialValues(0, dvgrammar.TYPE_NULL, nil, nil)
