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
	if WindowMaster.Prototype != nil && len(GlobalFunctionPool) < len(WindowMaster.Prototype.Fields) {
		for _, v := range WindowMaster.Prototype.Fields {
			GlobalFunctionPool[string(v.Name)] = v
		}
	}
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
	switch v.(type) {
	case string:
		return v.(string)
	}
	return AnyToString(v)
}

func (obj *DvObject) GetInt(key string) int {
	v, ok := obj.Get(key)
	if !ok || v == nil {
		return 0
	}
	switch v.(type) {
	case int:
		return v.(int)
	case int64:
		return int(v.(int64))
	}
	res, ok := AnyToNumberInt(v)
	if !ok {
		return 0
	}
	return int(res)
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

func (obj *DvObject) Delete(key string) {
	if obj == nil || obj.Properties == nil {
		return
	}
	delete(obj.Properties, key)
}

func (obj *DvObject) SetAtParent(key string, value interface{}, level int) {
	place := obj
	for ; level > 0 && place != nil; level-- {
		place = place.Prototype
	}
	if place == nil {
		return
	}
	if place.Properties == nil {
		place.Properties = make(map[string]interface{})
	}
	place.Properties[key] = value
}

func (obj *DvObject) DeleteAtParent(key string, level int) {
	place := obj
	for ; level > 0 && place != nil; level-- {
		place = place.Prototype
	}
	if place == nil || place.Properties == nil {
		return
	}
	delete(place.Properties, key)
}

func (obj *DvObject) ReadAtParent(key string, level int) (res interface{}, ok bool) {
	place := obj
	for ; level > 0 && place != nil; level-- {
		place = place.Prototype
	}
	if place == nil || place.Properties == nil {
		return
	}
	res, ok = place.Properties[key]
	return
}

func (obj *DvObject) FindFirstNotEmptyString(keys []string) string {
	n := len(keys)
	for i := 0; i < n; i++ {
		v, err := obj.EvaluateAnyTypeExpression(keys[i])
		if err == nil && v != nil {
			s := AnyToString(v)
			if s != "" {
				return s
			}
		}
	}
	return ""
}

func NewDvObjectWithSpecialValues(value interface{}, kind int, proto *DvObject, properties map[string]interface{}) *DvObject {
	return &DvObject{Value: value, Options: kind, Prototype: proto, Properties: properties}
}

var DvObject_null *DvObject = NewDvObjectWithSpecialValues(0, dvgrammar.TYPE_NULL, nil, nil)
