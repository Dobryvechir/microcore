/***********************************************************************
MicroCore
Copyright 2017 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strconv"
	"strings"
)

const (
	ERROR_POLICY_SKIP  = 0
	ERROR_POLICY_EMPTY = 1
	ERROR_POLICY_FATAL = 2
	ERROR_POLICY_NULL  = 3
)

func CreateLocalVariables(env *dvevaluation.DvObject, dat *dvevaluation.DvVariable) []string {
	locales := make([]string, 1, 16)
	locales[0] = "this"
	env.Set(locales[0], dat)
	if dat != nil {
		switch dat.Kind {
		case dvevaluation.FIELD_ARRAY:
			n := len(dat.Fields)
			for i := 0; i < n; i++ {
				k := strconv.Itoa(i)
				locales = append(locales, k)
				env.Set(k, dat.Fields[i])
			}
		case dvevaluation.FIELD_OBJECT:
			n := len(dat.Fields)
			for i := 0; i < n; i++ {
				if dat.Fields[i] == nil {
					continue
				}
				k := string(dat.Fields[i].Name)
				locales = append(locales, k)
				env.Set(k, dat.Fields[i])
			}
		default:
			v := "_value"
			n := "_key"
			env.Set(v, string(dat.Value))
			env.Set(n, string(dat.Name))
			locales = append(locales, v, n)
		}
	}
	return locales
}

func RemoveLocalVariables(env *dvevaluation.DvObject, locals []string) {
	n := len(locals)
	for i := 0; i < n; i++ {
		env.Delete(locals[i])
	}
}

func CreateObjectByArray(obj *dvevaluation.DvVariable, key string, value string, env *dvevaluation.DvObject, keyPolicy int, valuePolicy int) (res *dvevaluation.DvVariable, err error) {
	if obj == nil {
		return nil, nil
	}
	n := len(obj.Fields)
	res = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, n)}
	var v, k interface{}
createObjCycle:
	for i := 0; i < n; i++ {
		locals := CreateLocalVariables(env, obj.Fields[i])
		v, err = env.EvaluateAnyTypeExpression(value)
		if err != nil {
			switch valuePolicy {
			case ERROR_POLICY_SKIP:
				continue createObjCycle
			case ERROR_POLICY_FATAL:
				return
			case ERROR_POLICY_NULL:
				v = nil
			default:
				v = ""
			}
			err = nil
		}
		k, err = env.EvaluateAnyTypeExpression(key)
		if err != nil {
			switch valuePolicy {
			case ERROR_POLICY_SKIP:
				continue createObjCycle
			case ERROR_POLICY_FATAL:
				return
			case ERROR_POLICY_NULL:
				k = nil
			default:
				k = ""
			}
			err = nil
		}
		r := dvevaluation.AnyToDvVariable(v)
		if r == nil {
			r = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_NULL}
		}
		r.Name = dvevaluation.AnyToByteArray(k)
		res.Fields = append(res.Fields, r)
		RemoveLocalVariables(env, locals)
	}
	return
}

func RemoveByKeys(data *dvevaluation.DvVariable, keys []string, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, error) {
	n := len(keys)
	if data == nil || n == 0 || data.Kind != dvevaluation.FIELD_OBJECT {
		return data, nil
	}
	keyMap := make(map[string]string, n)
	for i := 0; i < n; i++ {
		keyMap[keys[i]] = ""
	}
	m := len(data.Fields)
	for i := 0; i < m; i++ {
		v := data.Fields[i]
		if v == nil {
			continue
		}
		_, ok := keyMap[string(v.Name)]
		if ok {
			if i != m-1 {
				data.Fields = append(data.Fields[:i], data.Fields[i+1:]...)
				m--
				i--
			} else {
				data.Fields = data.Fields[:i]
				m--
			}
		}
	}
	return data, nil
}

func ConcatObjects(dst *dvevaluation.DvVariable, src *dvevaluation.DvVariable) {
	if dst == nil || dst.Kind != dvevaluation.FIELD_OBJECT || src == nil || src.Kind != dvevaluation.FIELD_OBJECT {
		return
	}
	d := len(dst.Fields)
	s := len(src.Fields)
	if s == 0 {
		return
	}
	if dst.Fields == nil {
		dst.Fields = src.Fields
		return
	}
	mp := make(map[string]int, s+d)
	for i := 0; i < d; i++ {
		v := dst.Fields[i]
		if v == nil {
			continue
		}
		mp[string(v.Name)] = i
	}
	for i := 0; i < s; i++ {
		sr := src.Fields[i]
		if sr == nil {
			continue
		}
		key := string(sr.Name)
		pnt, ok := mp[key]
		if ok {
			dst.Fields[pnt] = sr
		} else {
			pnt = len(dst.Fields)
			dst.Fields = append(dst.Fields, sr)
			mp[key] = pnt
		}
	}
}

func ReplaceTextByObjectMap(s string, rules *dvevaluation.DvVariable) string {
	if rules == nil {
		return s
	}
	n := len(rules.Fields)
	for i := 0; i < n; i++ {
		r := rules.Fields[i]
		if r != nil {
			oldS := string(r.Name)
			newS := string(r.Value)
			s = strings.Replace(s, oldS, newS, -1)
		}
	}
	return s
}
