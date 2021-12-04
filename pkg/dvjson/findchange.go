/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

const (
	FLAG_CREATED = "CREATED"
	FLAG_UPDATED = "UPDATED"
	FLAG_REMOVED = "REMOVED"
)

func ConvertAnyToDvVariable(s interface{}) (*dvevaluation.DvVariable, error) {
	switch s.(type) {
	case *dvevaluation.DvObject:
		return s.(*dvevaluation.DvVariable), nil
	}
	return nil, errors.New("Cannot convert to DvVariable")
}

func FindChangeAny(srcAny interface{}, refAny interface{}, field string, algorithm string, env *dvevaluation.DvObject) (interface{}, error) {
	src, err := ConvertAnyToDvVariable(srcAny)
	if err != nil {
		return nil, err
	}
	ref, err := ConvertAnyToDvVariable(refAny)
	if err != nil {
		return nil, err
	}
	res, err := FindChange(src, ref, field, algorithm, env)
	return res, err
}
func addFlagField(src *dvevaluation.DvVariable, fieldName string, fieldValue string) *dvevaluation.DvVariable {
	if src == nil || src.Kind != dvevaluation.FIELD_OBJECT {
		return src
	}
	src.Fields = append(src.Fields, &dvevaluation.DvVariable{
		Kind:  dvevaluation.FIELD_STRING,
		Name:  []byte(fieldName),
		Value: []byte(fieldValue),
	})
	return src
}

func FindChange(src *dvevaluation.DvVariable, ref *dvevaluation.DvVariable, field string, algorithm string, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, error) {
	if ref == nil {
		return addFlagField(src, field, FLAG_CREATED), nil
	}
	if src == nil {
		return addFlagField(ref, field, FLAG_REMOVED), nil
	}
	if src.Kind != ref.Kind {
		return addFlagField(src, field, FLAG_CREATED), nil
	}
	var err error = nil
	switch src.Kind {
	case dvevaluation.FIELD_ARRAY:
		src, err = FindChangeInArray(src, ref, field, algorithm, env)
	case dvevaluation.FIELD_OBJECT:
		src, err = FindChangeInObject(src, ref, field, algorithm, env)
	default:
		if bytes.Equal(src.Value, ref.Value) {
			src = nil
		}
	}
	return src, err
}

func FindChangeInArray(src *dvevaluation.DvVariable, ref *dvevaluation.DvVariable, field string, algorithm string, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, error) {
	if src.QuickSearch == nil {
		src.CreateQuickInfoForObjectType()
	}
	if ref.QuickSearch == nil {
		ref.CreateQuickInfoForObjectType()
	}
	added, removed, updated, _, counterParts := src.FindDifferenceByQuickMap(ref,
		true, true, true, false, true, false)
	na := len(added.Fields)
	nr := len(removed.Fields)
	nu := len(updated.Fields)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, na+nu, na+nu+nr)}
	for i := 0; i < na; i++ {
		f := added.Fields[i]
		res.Fields[i] = addFlagField(f, field, FLAG_CREATED)
	}
	for i := 0; i < nu; i++ {
		f := updated.Fields[i]
		if IsChangeUpdatedPossible(f) {
			fld, err := FindChangeInObject(f, counterParts, field, algorithm, env)
			if err != nil {
				return nil, err
			}
			res.Fields[na+i] = fld
		} else {
			res.Fields[na+i] = addFlagField(f, field, FLAG_CREATED)
		}
	}
	res.Fields = AppendRemovedFields(removed, res.Fields, field, env)
	return res, nil
}

func AppendRemovedFields(arr *dvevaluation.DvVariable, dst []*dvevaluation.DvVariable, field string, env *dvevaluation.DvObject) []*dvevaluation.DvVariable {
	n := len(arr.Fields)
	for i := 0; i < n; i++ {
		f := arr.Fields[i]
		if f == nil {
			continue
		}
		switch f.Kind {
		case dvevaluation.FIELD_OBJECT:
			f = addFlagField(GetMinifiedWithId(f, env), field, FLAG_REMOVED)
			dst = append(dst, f)
		case dvevaluation.FIELD_ARRAY:
			fa := &dvevaluation.DvVariable{
				Kind:   dvevaluation.FIELD_ARRAY,
				Fields: make([]*dvevaluation.DvVariable, 0, len(f.Fields)),
			}
			fa.Fields = AppendRemovedFields(f, fa.Fields, field, env)
		}
	}
	return dst
}

func GetMinifiedWithId(f *dvevaluation.DvVariable, env *dvevaluation.DvObject) *dvevaluation.DvVariable {
	if f.QuickSearch == nil {
		f.CreateQuickInfoForObjectType()
	}
	e, ok := f.QuickSearch.Looker["id"]
	if !ok {
		e, ok = f.QuickSearch.Looker["Id"]
		if !ok {
			e, ok = f.QuickSearch.Looker["ID"]
			if !ok {
				id, ok := env.Get("ID")
				if ok {
					e, ok = f.QuickSearch.Looker[dvevaluation.AnyToString(id)]
				}
			}
		}
	}
	if ok {
		return &dvevaluation.DvVariable{
			Kind:   dvevaluation.FIELD_OBJECT,
			Fields: []*dvevaluation.DvVariable{e},
		}
	}
	return GetMinifiedWithSimpleFields(f, env)
}

func GetMinifiedWithSimpleFields(f *dvevaluation.DvVariable, env *dvevaluation.DvObject) *dvevaluation.DvVariable {
	n := len(f.Fields)
	r := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, n)}
	for i := 0; i < n; i++ {
		fld := f.Fields[i]
		obj := &dvevaluation.DvVariable{Kind: fld.Kind, Name: fld.Name, Value: fld.Value}
		if obj.Kind == dvevaluation.FIELD_OBJECT || obj.Kind == dvevaluation.FIELD_ARRAY {
			obj.Kind = dvevaluation.FIELD_NULL
		}
		r.Fields[i] = obj
	}
	return r
}

func IsChangeUpdatedPossible(f *dvevaluation.DvVariable) bool {
	if f == nil {
		return false
	}
	switch f.Kind {
	case dvevaluation.FIELD_OBJECT:
		return true
	case dvevaluation.FIELD_ARRAY:
		n := len(f.Fields)
		for i := 0; i < n; i++ {
			if !IsChangeUpdatedPossible(f.Fields[i]) {
				return false
			}
			return true
		}
	}
	return false
}

func FindChangeInObject(src *dvevaluation.DvVariable, ref *dvevaluation.DvVariable, field string, algorithm string, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, error) {
	return nil, nil
}
