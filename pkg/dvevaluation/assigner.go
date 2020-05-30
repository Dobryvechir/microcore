// package dvevaluation covers expression calculations
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strconv"
	"strings"
)

var ObjectMaster *DvVariable = RegisterMasterVariable("Object", &DvVariable{Tp: JS_TYPE_OBJECT})
var ArrayMaster *DvVariable = RegisterMasterVariable("Array", &DvVariable{Tp: JS_TYPE_OBJECT})

func AssignVariableDirect(parent *DvVariable, value *DvVariable) error {
	if value == nil {
		value = &DvVariable{}
	}
	parent.Refs = value.Refs
	parent.Value = value.Value
	parent.Tp = value.Tp
	parent.Fn = value.Fn
	parent.Prototype = value.Prototype
	return nil
}

func AssignVariable(parent *DvVariable, keys []string, value *DvVariable, force bool) error {
	if parent == nil {
		return errors.New("Cannot assign to undefined [keys:" + strings.Join(keys, ",") + "]")
	}
	if value == nil {
		value = &DvVariable{}
	}
	l := len(keys)
	if l == 0 {
		return AssignVariableDirect(parent, value)
	}
	l--
	var ok bool
	var child *DvVariable
	for i := 0; i < l; i++ {
		key := keys[i]
		if parent.Tp == JS_TYPE_UNDEFINED || parent.Tp == JS_TYPE_NULL {
			if force {
				parent.Tp = JS_TYPE_OBJECT
			} else {
				return errors.New("Cannot assign key " + key + " to " + typeOfSpecific[parent.Tp] + " [keys:" + strings.Join(keys, ",") + "]")
			}
		}
		ok = parent.Refs != nil
		if ok {
			child, ok = parent.Refs[key]
		} else if force {
			parent.Refs = make(map[string]*DvVariable)
		}
		if !ok {
			if force {
				child = &DvVariable{}
				parent.Refs[key] = child
			} else {
				return errors.New("Cannot assign key " + key + " to  undefined [keys:" + strings.Join(keys, ",") + "]")
			}
		}
		parent = child
	}
	key := keys[l]
	if parent.Tp == JS_TYPE_UNDEFINED || parent.Tp == JS_TYPE_NULL {
		if force {
			parent.Tp = JS_TYPE_OBJECT
		} else {
			return errors.New("Cannot assign key " + key + " to " + typeOfSpecific[parent.Tp] + " [keys:" + strings.Join(keys, ",") + "]")
		}
	}
	if parent.Refs == nil {
		parent.Refs = make(map[string]*DvVariable)
	}
	parent.Refs[key] = value
	return nil
}

func GetVariableByKeys(parent *DvVariable, keys []string) (thisValue *DvVariable, child *DvVariable, prototyped bool, err error) {
	prototyped = false
	err = nil
	thisValue = parent
	child = parent
	l := len(keys)
	if l == 0 {
		return
	}
	if parent == nil {
		err = errors.New("Cannot get object from undefined")
		return
	}
	var ok bool
	for i := 0; i < l; i++ {
		parent = child
		key := keys[i]
		thisValue = parent
		if parent.Tp == JS_TYPE_UNDEFINED || parent.Tp == JS_TYPE_NULL {
			err = errors.New("Cannot get key " + key + " from " + typeOfSpecific[parent.Tp])
			return
		}
		ok = parent.Refs != nil
		if ok {
			child, ok = parent.Refs[key]
		}
		prototyped = false
		if !ok && parent.Prototype != nil {
			child = parent.Prototype.ObjectInPrototypedChain(key)
			if child != nil {
				ok = true
				prototyped = true
			}
		}
		if !ok {
			err = errors.New("Cannot get key " + key + " from  undefined")
			return
		}
	}
	return
}

func DvVariableFromString(parent *DvVariable, data string) *DvVariable {
	if parent == nil {
		parent = &DvVariable{}
	}
	parent.Value = data
	parent.Tp = JS_TYPE_STRING
	return parent
}

func DvVariableFromInt(parent *DvVariable, data int) *DvVariable {
	if parent == nil {
		parent = &DvVariable{}
	}
	parent.Value = strconv.Itoa(data)
	parent.Tp = JS_TYPE_NUMBER
	return parent
}

func DvVariableFromArray(parent *DvVariable, data []*DvVariable) *DvVariable {
	if parent == nil {
		parent = &DvVariable{Prototype: ArrayMaster}
	}
	parent.Refs = make(map[string]*DvVariable)
	parent.Tp = JS_TYPE_ARRAY
	n := len(data)
	parent.Refs[LENGTH_PROPERTY] = DvVariableFromInt(nil, n)
	for i := 0; i < n; i++ {
		parent.Refs[strconv.Itoa(i)] = data[i]
	}
	return parent
}

func DvVariableFromMap(parent *DvVariable, data map[string]*DvVariable, reuse bool) *DvVariable {
	if parent == nil {
		parent = &DvVariable{Prototype: ObjectMaster}
	}
	if data == nil {
		data = make(map[string]*DvVariable)
		reuse = true
	}
	parent.Tp = JS_TYPE_OBJECT
	if reuse {
		parent.Refs = data
	} else {
		parent.Refs = make(map[string]*DvVariable)
		for k, v := range data {
			parent.Refs[k] = v
		}
	}
	return parent
}

func AssignArrayStringToVariable(parent *DvVariable, variableDefinition string, data []string, force bool) error {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning array string to "+variableDefinition+" due to this name conversion")
	}
	value := DvVariableFromStringArray(nil, data)
	err = AssignVariable(parent, keys, value, force)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning array string to "+variableDefinition+" due to assignment")
	}
	return nil
}

func AssignMapStringToVariable(parent *DvVariable, variableDefinition string, data map[string]string, force bool) error {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning map string to "+variableDefinition+" due to this name conversion")
	}
	value := DvVariableFromStringMap(nil, data)
	err = AssignVariable(parent, keys, value, force)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning map string to "+variableDefinition+" due to assignment")
	}
	return nil
}

func AssignIntToVariable(parent *DvVariable, variableDefinition string, data int, force bool) error {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning int to "+variableDefinition+" due to this name conversion")
	}
	value := DvVariableFromInt(nil, data)
	err = AssignVariable(parent, keys, value, force)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning int to "+variableDefinition+" due to assignment")
	}
	return nil
}

func AssignVariableToVariable(parent *DvVariable, variableDefinition string, data *DvVariable, force bool) error {
	keys, err := ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning Variable to "+variableDefinition+" due to this name conversion")
	}
	err = AssignVariable(parent, keys, data, force)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At assigning Variable to "+variableDefinition+" due to assignment")
	}
	return nil
}

func GetVariableByDefinition(parent *DvVariable, variableDefinition string) (thisValue *DvVariable, child *DvVariable, prototyped bool, err error, keys []string) {
	keys, err = ConvertVariableNameToKeys(variableDefinition)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At getting Variable from "+variableDefinition+" due to this name conversion")
		return
	}
	thisValue, child, prototyped, err = GetVariableByKeys(parent, keys)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At getting Variable from "+variableDefinition+" due to this operation")
		return
	}
	return
}

func GetIntFromVariable(parent *DvVariable, variableDefinition string, force bool) (res int, err error) {
	res = 0
	var v *DvVariable
	_, v, _, err, _ = GetVariableByDefinition(parent, variableDefinition)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At getting int from "+variableDefinition)
		return
	}
	res, err = v.GetIntValue(force)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At getting int from "+variableDefinition)
	}
	return
}

func GetStringFromVariable(parent *DvVariable, variableDefinition string, force bool) (res string, err error) {
	var v *DvVariable
	_, v, _, err, _ = GetVariableByDefinition(parent, variableDefinition)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At getting string from "+variableDefinition)
		return
	}
	res = v.GetStringValue()
	return
}

func ModifySimpleValueAfterGet(thisVal *DvVariable, child *DvVariable, prototyped bool, value string, tp int, keys []string) error {
	n := len(keys)
	if thisVal == nil {
		return errors.New("Cannot modify variable by keys " + strings.Join(keys, ", "))
	}
	if child == nil || prototyped {
		if n == 0 {
			return errors.New("Cannot modify variable with no keys")
		}
		if thisVal.Refs == nil {
			thisVal.Refs = make(map[string]*DvVariable)
		}
		thisVal.Refs[keys[n-1]] = &DvVariable{Value: value, Tp: tp}
		return nil
	}
	child.Value = value
	child.Tp = tp
	return nil
}

func GetIntFromVariableAndModify(parent *DvVariable, variableDefinition string, modify int, after bool, force bool) (res int, err error) {
	res = 0
	var (
		v, thisVal *DvVariable
		prototyped bool
		keys       []string
	)
	thisVal, v, prototyped, err, keys = GetVariableByDefinition(parent, variableDefinition)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "While getting and modifying int from "+variableDefinition)
		return
	}
	res, err = v.GetIntValue(force)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "While getting and modifying int from "+variableDefinition+" due to int conversion")
		return
	}
	newRes := res + modify
	if !after {
		res = newRes
	}
	err = ModifySimpleValueAfterGet(thisVal, v, prototyped, strconv.Itoa(newRes), JS_TYPE_NUMBER, keys)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "While getting and modifying int from "+variableDefinition+" due to int modification")
	}
	return
}

func GetIntFromVariableAndIncrementAfter(parent *DvVariable, variableDefinition string, force bool) (res int, err error) {
	res, err = GetIntFromVariableAndModify(parent, variableDefinition, 1, true, force)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At ++ after operation from "+variableDefinition)
	}
	return
}
