/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"sort"
)

type IterateProcessor func(string, interface{}, int, interface{}) (interface{}, bool)
type FilterProcessor func(string, interface{}, int, interface{}) (bool, bool)
type SortComparator func(d1 *dvevaluation.DvVariable, d2 *dvevaluation.DvVariable) int

func IterateOnAnyType(val interface{}, processor IterateProcessor, initial interface{}) interface{} {
	res := initial
	index := 0
	toBreak := false
	switch val.(type) {
	case map[string]string:
		for k, v := range val.(map[string]string) {
			res, toBreak = processor(k, v, index, res)
			if toBreak {
				break
			}
			index++
		}
	case map[string]interface{}:
		for k, v := range val.(map[string]interface{}) {
			res, toBreak = processor(k, v, index, res)
			if toBreak {
				break
			}
			index++
		}
	case *dvevaluation.DvVariable:
		fieldInfo := val.(*dvevaluation.DvVariable)
		fields := fieldInfo.Fields
		n := len(fields)
		for index = 0; index < n; index++ {
			f := fields[index]
			res, toBreak = processor(string(f.Name), f, index, res)
			if toBreak {
				break
			}
		}
	}
	return res
}

func IterateFilterOnAnyType(val interface{}, processor FilterProcessor) interface{} {
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	index := 0
	var toAdd, toBreak bool
	switch val.(type) {
	case *dvevaluation.DvVariable:
		fieldInfo := val.(*dvevaluation.DvVariable)
		if fieldInfo.Kind == dvevaluation.FIELD_OBJECT {
			res.Kind = dvevaluation.FIELD_OBJECT
		}
		fields := fieldInfo.Fields
		n := len(fields)
		res.Fields = make([]*dvevaluation.DvVariable, 0, n)
		for index = 0; index < n; index++ {
			f := fields[index]
			toAdd, toBreak = processor(string(f.Name), f, index, val)
			if toAdd {
				res.Fields = append(res.Fields, f)
			}
			if toBreak {
				break
			}
		}
	}
	return res
}

func IterateSortOnAnyType(val interface{}, processor SortComparator) interface{} {
	var res *dvevaluation.DvVariable
	switch val.(type) {
	case *dvevaluation.DvVariable:
		fieldInfo := val.(*dvevaluation.DvVariable)
		fields := fieldInfo.Fields
		n := len(fields)
		res = &dvevaluation.DvVariable{
			Kind:   fieldInfo.Kind,
			Name:   fieldInfo.Name,
			Fields: append(make([]*dvevaluation.DvVariable, 0, n), fields...),
			Value:  fieldInfo.Value,
		}
		fields = res.Fields
		sort.SliceStable(fields, func(i int, j int) bool {
			return processor(fields[i], fields[j]) < 0
		})
	}
	return res
}

func IterateFilterByExpression(val interface{}, expression string, env *dvevaluation.DvObject, errIsCritical bool) (res interface{}, err error) {
	res = IterateFilterOnAnyType(val, func(key string, item interface{}, index int, v interface{}) (bool, bool) {
		env.Set("KEY", key)
		env.Set("INDEX", index)
		switch item.(type) {
		case *dvevaluation.DvVariable:
			data, er := item.(*dvevaluation.DvVariable).EvaluateDvFieldItem(expression, env)
			if er != nil {
				if errIsCritical {
					err = er
					return false, true
				}
			}
			return data, false
		}
		return false, false
	})
	return
}

func IterateSortByFields(val interface{}, fields []string, env *dvevaluation.DvObject) (res interface{}, err error) {
	res = IterateSortOnAnyType(val, func(d1 *dvevaluation.DvVariable, d2 *dvevaluation.DvVariable) int {
		if d1 == nil {
			if d2 == nil {
				return 0
			}
			return -1
		}
		if d2 == nil {
			return 1
		}
		return d1.CompareDvFieldByFields(d2, fields)
	})
	return
}
