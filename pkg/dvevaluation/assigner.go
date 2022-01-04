// package dvevaluation covers expression calculations
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strconv"
	"strings"
)

var ObjectMaster *DvVariable = RegisterMasterVariable("Object", &DvVariable{Kind: FIELD_OBJECT})
var ArrayMaster *DvVariable = RegisterMasterVariable("Array", &DvVariable{Kind: FIELD_OBJECT})

func AssignVariableDirect(parent *DvVariable, val interface{}) error {
	var value *DvVariable
	if val == nil {
		value = &DvVariable{Kind: FIELD_UNDEFINED}
	} else {
		value = AnyToDvVariable(val)
	}
	parent.Fields = value.Fields
	parent.Value = value.Value
	parent.Kind = value.Kind
	parent.Extra = value.Extra
	parent.Prototype = value.Prototype
	return nil
}

func (item *DvVariable) FindChildByKey(key string) (*DvVariable, bool) {
	if item == nil {
		return nil, false
	}
	n := len(item.Fields)
	k := []byte(key)
	for i := 0; i < n; i++ {
		f := item.Fields[i]
		if f != nil && bytes.Equal(f.Name, k) {
			return f, true
		}
	}
	return nil, false
}

func (item *DvVariable) FindChildIndexByKey(key string) int {
	if item == nil {
		return -1
	}
	n := len(item.Fields)
	k := []byte(key)
	for i := 0; i < n; i++ {
		f := item.Fields[i]
		if f != nil && bytes.Equal(f.Name, k) {
			return i
		}
	}
	return -1
}

func (item *DvVariable) DeleteChildByIndex(ind int) {
	if item == nil || ind < 0 {
		return
	}
	n := len(item.Fields)
	if ind < n {
		if ind == n-1 {
			item.Fields = item.Fields[:ind:ind]
		} else {
			item.Fields = append(item.Fields[:ind], item.Fields[ind+1:]...)
		}
	}
}

func (item *DvVariable) MakeCopyWithNewKey(key string) *DvVariable {
	k := []byte(key)
	if item == nil {
		return &DvVariable{Name: k, Kind: FIELD_NULL}
	}
	return &DvVariable{
		Name:      k,
		Kind:      item.Kind,
		Value:     item.Value,
		Fields:    item.Fields,
		Extra:     item.Extra,
		Prototype: item.Prototype,
	}
}

func AssignVariable(parent *DvVariable, keys []string, value interface{}, force bool) error {
	if parent == nil {
		return errors.New("Cannot assign to undefined [keys:" + strings.Join(keys, ",") + "]")
	}
	if value == nil {
		value = &DvVariable{Kind: FIELD_NULL}
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
		if parent.Kind == FIELD_UNDEFINED || parent.Kind == FIELD_NULL {
			if force {
				parent.Kind = FIELD_OBJECT
			} else {
				return errors.New("Cannot assign key " + key + " to " + typeOfSpecific[parent.Kind] + " [keys:" + strings.Join(keys, ",") + "]")
			}
		}
		ok = parent.Fields == nil
		if !ok {
			if force {
				parent.Fields = make([]*DvVariable, 0, 7)
			}
		} else {
			child, ok = parent.FindChildByKey(key)
		}
		if !ok {
			if force {
				child = &DvVariable{Name: []byte(key), Kind: FIELD_OBJECT, Fields: make([]*DvVariable, 0, 7)}
				parent.Fields = append(parent.Fields, child)
			} else {
				return errors.New("Cannot assign key " + key + " to  undefined [keys:" + strings.Join(keys, ",") + "]")
			}
		}
		parent = child
	}
	key := keys[l]
	if parent.Kind == FIELD_UNDEFINED || parent.Kind == FIELD_NULL {
		if force {
			parent.Kind = FIELD_OBJECT
		} else {
			return errors.New("Cannot assign key " + key + " to " + typeOfSpecific[parent.Kind] + " [keys:" + strings.Join(keys, ",") + "]")
		}
	}
	if parent.Fields == nil {
		parent.Fields = make([]*DvVariable, 0, 7)
	}
	parent.Fields = append(parent.Fields, AnyToDvVariable(value).MakeCopyWithNewKey(key))
	return nil
}

func GetVariableByKeys(parent *DvVariable, keys []string, silent bool) (thisValue *DvVariable, child *DvVariable, prototyped bool, err error) {
	prototyped = false
	err = nil
	thisValue = parent
	child = parent
	l := len(keys)
	if l == 0 {
		return
	}
	if parent == nil {
		if silent {
			return
		}
		err = errors.New("Cannot get object from undefined")
		return
	}
	var ok bool
	for i := 0; i < l; i++ {
		parent = child
		key := keys[i]
		thisValue = parent
		if parent.Kind == FIELD_UNDEFINED || parent.Kind == FIELD_NULL {
			err = errors.New("Cannot get key " + key + " from " + typeOfSpecific[parent.Kind])
			return
		}
		ok = parent.Fields != nil
		if ok {
			child, ok = parent.FindChildByKey(key)
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
	parent.Value = []byte(data)
	parent.Kind = FIELD_STRING
	return parent
}

func DvVariableFromInt(parent *DvVariable, data int) *DvVariable {
	if parent == nil {
		parent = &DvVariable{}
	}
	parent.Value = []byte(strconv.Itoa(data))
	parent.Kind = FIELD_NUMBER
	return parent
}

func DvVariableFromArray(parent *DvVariable, data []*DvVariable) *DvVariable {
	if parent == nil {
		parent = &DvVariable{Prototype: ArrayMaster}
	}
	n := len(data)
	parent.Fields = make([]*DvVariable, n)
	parent.Kind = FIELD_ARRAY
	for i := 0; i < n; i++ {
		parent.Fields[i] = data[i]
	}
	return parent
}

func DvVariableFromMap(parent *DvVariable, data map[string]*DvVariable) *DvVariable {
	if parent == nil {
		parent = &DvVariable{Prototype: ObjectMaster}
	}
	n := len(data)
	parent.Fields = make([]*DvVariable, n, n+1)
	if data != nil {
		i := 0
		for k, v := range data {
			item := v.MakeCopyWithNewKey(k)
			if i < n {
				parent.Fields[i] = item
				i++
			} else {
				parent.Fields = append(parent.Fields, item)
			}
		}
	}
	parent.Kind = FIELD_OBJECT
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

func AssignVariableToVariable(parent *DvVariable, variableDefinition string, data interface{}, force bool) error {
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
	thisValue, child, prototyped, err = GetVariableByKeys(parent, keys, true)
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
		if thisVal.Fields == nil {
			thisVal.Fields = make([]*DvVariable, 0, 7)
		}
		key := keys[n-1]
		f, ok := thisVal.FindChildByKey(key)
		if ok {
			f.Value = []byte(value)
			f.Kind = tp
		} else {
			thisVal.Fields = append(thisVal.Fields,
				&DvVariable{Value: []byte(value), Kind: tp, Name: []byte(key)})
		}
		return nil
	}
	child.Value = []byte(value)
	child.Kind = tp
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
	err = ModifySimpleValueAfterGet(thisVal, v, prototyped, strconv.Itoa(newRes), FIELD_NUMBER, keys)
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
