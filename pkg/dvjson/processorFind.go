/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
)

func ProcessorFind(parent *DvFieldInfo, params string) (*DvFieldInfo, error) {
	pattern, err := JsonFullParser([]byte(params))
	if err != nil {
		return nil, err
	}
	return FindInArray(parent, pattern, 0)
}

func FindInArray(parent *DvFieldInfo, pattern *DvFieldInfo, fromIndex int) (*DvFieldInfo, error) {
	if parent == nil || len(parent.Fields) <= fromIndex {
		return nil, nil
	}
	n := len(parent.Fields)
	for i := fromIndex; i < n; i++ {
		f := parent.Fields[i]
		if MatchDvFieldInfo(f, pattern) {
			return f, nil
		}
	}
	return nil, nil
}

func MatchDvFieldInfo(model *DvFieldInfo, pattern *DvFieldInfo) bool {
	if pattern == nil {
		return true
	}
	if model == nil {
		return false
	}
	n := len(pattern.Fields)
	switch pattern.Kind {
	case FIELD_OBJECT:
		for i := 0; i < n; i++ {
			item := model.ReadSimpleChild(string(pattern.Fields[i].Name))
			if item == nil {
				return false
			}
			if !MatchDvFieldInfo(item, pattern.Fields[i]) {
				return false
			}
		}
		return true
	case FIELD_ARRAY:
		if model.Kind == FIELD_ARRAY && len(model.Fields) >= n {
			for i := 0; i < n; i++ {
				if !MatchDvFieldInfo(model.Fields[i], pattern.Fields[i]) {
					return false
				}
			}
			return true
		}
	default:
		if model.Kind != FIELD_ARRAY && model.Kind != FIELD_OBJECT && bytes.Equal(model.Value, pattern.Value) {
			return true
		}
	}
	return false
}
