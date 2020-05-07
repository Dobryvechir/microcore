/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
)

func DvObjectInternalToString(obj *DvObject) string {
	if v, ok := AnyWithTypeToString(obj.GetObjectType(), obj.Value); ok {
		return v
	}
	return "[object Object]"
}

func (obj *DvObject) GetObjectType() int {
	if obj == nil {
		return 0
	}
	return obj.Options & dvgrammar.TYPE_MASK
}

func (obj *DvObject) ToString() string {
	if fn, ok := obj.Get("toString"); ok && IsFunction(fn) {
		v, err := ExecFunction(fn, obj, nil)
		if err == nil {
			return AnyToString(v)
		}
	}
	return DvObjectInternalToString(obj)
}

func (obj *DvObject) ToNumber() float64 {
	if obj == nil {
		return math.NaN()
	}
	return AnyWithTypeToNumber(obj.GetObjectType(), obj.Value)
}

func (obj *DvObject) ToNumberInt() (int64, bool) {
	if obj == nil {
		return 0, false
	}
	return AnyWithTypeToNumberInt(obj.GetObjectType(), obj.Value)
}

func (obj *DvObject) ToBoolean() bool {
	if obj == nil {
		return false
	}
	return AnyWithTypeToBoolean(obj.GetObjectType(), obj.Value, true)
}

func (obj *DvObject) AssignProperties(properties map[string]interface{}) {
	if obj == nil {
		return
	}
	if obj.Properties == nil {
		obj.Properties = make(map[string]interface{})
	}
	for k, v := range properties {
		obj.Properties[k] = v
	}
}

func (obj *DvObject) SetProperty(name string, value interface{}) {
	if obj == nil {
		return
	}
	if obj.Properties == nil {
		obj.Properties = make(map[string]interface{})
	}
	obj.Properties[name] = value
}
