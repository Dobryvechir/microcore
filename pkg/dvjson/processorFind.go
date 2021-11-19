/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strings"
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

func CollectValuesByMap(data interface{}, params map[string]string, env *dvevaluation.DvObject) (res map[string]interface{}) {
	res = make(map[string]interface{})
	for k, _ := range params {
		val, err := CollectValueByKey(data, k, env)
		if err != nil {
			res[k] = val
		}
	}
	return res
}

func CollectValueByKey(data interface{}, key string, env *dvevaluation.DvObject) (interface{}, error) {
	switch data.(type) {
	case *DvFieldInfo:
		res, err := data.(*DvFieldInfo).ReadChildOfAnyLevel(key, env)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, errors.New("Unknown type to extrace " + key)
}

func CollectJsonVariables(data interface{}, params map[string]string, env *dvevaluation.DvObject, anyway bool) {
	if params == nil || (data == nil && !anyway) {
		return
	}
	src := CollectValuesByMap(data, params, env)
	CollectVariablesByAnyMap(src, params, env, anyway)
}

func CollectVariablesByStringMap(src map[string]string, params map[string]string, data *dvevaluation.DvObject, anyway bool) {
	if params == nil || (src == nil && !anyway) {
		return
	}
	for k, v := range params {
		if v != "" && v[0] >= 'A' && v[0] <= 'Z' {
			p := strings.Index(v, ":")
			if p > 0 {
				v = v[:p]
			}
			if src != nil {
				v1, ok := src[k]
				if ok {
					data.Properties[v] = v1
				} else if anyway {
					data.Properties[v] = ""
				}
			} else {
				data.Properties[v] = ""
			}
		}
	}
}

func CollectVariablesByAnyMap(src map[string]interface{}, params map[string]string, data *dvevaluation.DvObject, anyway bool) {
	if params == nil || (src == nil && !anyway) {
		return
	}
	for k, v := range params {
		if v != "" && v[0] >= 'A' && v[0] <= 'Z' {
			p := strings.Index(v, ":")
			if p > 0 {
				v = v[:p]
			}
			if src != nil {
				v1, ok := src[k]
				if ok {
					data.Properties[v] = v1
				} else if anyway {
					data.Properties[v] = ""
				}
			} else {
				data.Properties[v] = ""
			}
		}
	}
}
